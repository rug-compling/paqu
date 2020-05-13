package main

import (
	"github.com/pebbe/util"

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
		"tei":          "TEI-bestand)",
		"folia-arch":   "FoLiA-bestanden in zip/tar-bestand",
		"tei-arch":     "TEI-bestanden in zip/tar-bestand",
	}

	reRunLabel = regexp.MustCompile(`^[^|]*\|\s*$`)
	reEndPoint = regexp.MustCompile(`[.!?]\s*$`)
	reMidPoint = regexp.MustCompile(`\pL\pL\pP*[.!?]\s+\S`)
)

func setinvoer(soort string, id string, isarch bool) error {
	s := soort
	if isarch {
		s = soort + "-arch"
	}
	_, err := sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `params` = %q, `msg` = %q WHERE `id` = %q",
		Cfg.Prefix, s, "Bron: "+invoertabel[soort], id))
	return err
}

func invoersoort(data, id string) (string, error) {

	set := func(soort string) (string, error) {
		return soort, setinvoer(soort, id, false)
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
	rd := util.NewReader(fp)
	for i := 0; i < 20; i++ {
		line, e := rd.ReadLineString()
		if e != nil {
			break
		}
		line = strings.ToUpper(line)
		line = strings.Replace(line, "\000", "", -1) // utf-16, utf-32, grove methode
		if strings.TrimSpace(line) == "" ||
			strings.HasPrefix(line, "##PAQU") ||
			strings.HasPrefix(line, "##META") ||
			line[0] == '%' ||
			reRunLabel.MatchString(line) {
			i--
		} else {
			lines = append(lines, line)
		}
	}
	ln := len(lines)
	if ln < 2 {
		return set("run")
	}

	endletter := 0
	midpoint := 0
	nlabel := 0
	for _, line := range lines {
		if !reEndPoint.MatchString(line) {
			endletter++
		}
		midpoint += len(reMidPoint.FindAllString(line, -1))
		if strings.Contains(line, "|") {
			nlabel++
		}
	}
	if nlabel < ln && (endletter > ln/3 || midpoint > endletter/2) {
		return set("run")
	}

	soort := "line"

	ntok := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
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
