package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	invoertabel = map[string]string{
		"auto":         "moet nog bepaald worden",
		"run":          "doorlopende tekst",
		"line":         "een zin per regel",
		"line-lbl":     "een zin per regel, met labels",
		"line-tok":     "een zin per regel, getokeniseerd",
		"line-lbl-tok": "een zin per regel, met labels, getokeniseerd",
		"xmlzip":       "Alpino XML-bestanden in zipbestand of tarbestand",
		"xmlzip-d":     "afgeleid van een of meer andere corpora",
		"xmlzip-p":     "afgeleid van een of meer andere corpora (beschermd)",
		"dact":         "Dact-bestand",
		"folia":        "FoLiA-bestand",
		"tei":          "TEI-bestand",
	}

	reEndPoint = regexp.MustCompile(`[.!?]\s*$`)
	reMidPoint = regexp.MustCompile(`\pL\pL\pP*[.!?]\s+\S`)
)

func setinvoer(db *sql.DB, soort string, id string) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `params` = %q, `msg` = %q WHERE `id` = %q",
		Cfg.Prefix, soort, "Bron: "+invoertabel[soort], id))
	return err
}

func invoersoort(db *sql.DB, data, id string) (string, error) {

	set := func(soort string) (string, error) {
		return soort, setinvoer(db, soort, id)
	}

	fp, err := os.Open(data)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	b := make([]byte, 200)
	n, _ := io.ReadFull(fp, b)
	fp.Seek(0, 0)

	if n >= 3 {
		s := string(b[:4])
		if s == "\x00\x06\x15\x61" || s == "\x61\x15\x06\x00" || s == "\x00\x05\x31\x62" || s == "\x62\x31\x05\x00" {
			return set("dact")
		}
	}

	if n > 15 {
		s := string(b[12:16])
		if s == "\x00\x06\x15\x61" || s == "\x61\x15\x06\x00" ||
			s == "\x00\x05\x31\x62" || s == "\x62\x31\x05\x00" ||
			s == "\x00\x04\x22\x53" || s == "\x53\x22\x04\x00" {
			return set("dact")
		}
	}

	if strings.Contains(string(b), "<FoLiA") {
		return set("folia")
	}

	if strings.Contains(string(b), "<TEI") {
		return set("tei")
	}

	lines := make([]string, 0, 20)
	scanner := bufio.NewScanner(fp)
	for i := 0; i < 20 && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}
	ln := len(lines)
	if ln < 2 || scanner.Err() != nil {
		return set("run")
	}

	endletter := 0
	midpoint := 0
	for _, line := range lines {
		if !reEndPoint.MatchString(line) {
			endletter++
		}
		midpoint += len(reMidPoint.FindAllString(line, -1))
	}
	if endletter > ln/3 || midpoint > endletter/2 {
		return set("run")
	}

	soort := "line"

	nlabel := 0
	ntok := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			nlabel++
			ntok++
			continue
		}
		if strings.Contains(line, "|") {
			nlabel++
		}
		if strings.HasSuffix(line, " .") || strings.HasSuffix(line, " !") || strings.HasSuffix(line, " ?") {
			ntok++
		}
	}
	if nlabel == ln {
		soort += "-lbl"
	}
	if ntok > (3*ln)/4 {
		soort += "-tok"
	}

	return set(soort)
}
