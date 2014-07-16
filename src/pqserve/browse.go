package main

import (
	"github.com/pebbe/util"

	"compress/gzip"
	"fmt"
	"html"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type ZinArchFile struct {
	zin  string
	arch int
	file int
}

// TAB: browse (zinnen)
func browse(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := first(q.r, "id")
	if !q.myprefixes[id] {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	datadir := path.Join(paqudir, "data", id)
	fp, err := os.Open(path.Join(datadir, "summary.txt.gz"))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}
	defer fp.Close()
	gz, err := gzip.NewReader(fp)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}
	defer gz.Close()

	rd := util.NewReader(gz)
	line, _ := rd.ReadLineString()
	a := strings.SplitN(line, "\t", 3)
	nline, _ := strconv.Atoi(a[0])
	nerr, _ := strconv.Atoi(a[1])

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "Overzicht", 0)
	fmt.Fprintf(q.w, `
<script type="text/javascript"><!--
  function formclear(f) {
    f.lbl.value = "";
  }
//--></script>
Corpus: <b>%s</b>
<p>
Bron: %s
<p>
`, q.desc[id], a[2])

	if nerr > 0 {
		fmt.Fprintf(q.w, `
Er waren problemen met %d van de %d zinnen:
<p>
<table class="corpora">
<tr><th>Label<th>Fout<th>Zin</tr>
`, nerr, nline)
		lineno := 0
		for {
			lineno++
			line, e := rd.ReadLineString()
			if e != nil {
				break
			}
			a := strings.SplitN(line, "\t", 4)
			eo := "even"
			if lineno%2 == 1 {
				eo = "odd"
			}
			if lineno == 1 {
				eo += " first"
			}
			if lineno == nerr {
				eo += " last"
			}
			fmt.Fprintf(q.w, "<tr class=\"%s\"><td class=\"odd first\">%s<td class=\"even\">%s&nbsp;|&nbsp;%s<td class=\"odd\">%s",
				eo,
				html.EscapeString(a[0]), html.EscapeString(a[1]), a[2], html.EscapeString(a[3]))
		}
		fmt.Fprint(q.w, "</table>\n<p>\n")
	}

	// HTML-uitvoer van het formulier
	fmt.Fprintf(q.w, `
<form action="browse" method="get" accept-charset="utf-8">
<input type="hidden" name="id", value="%s">
Label: <input type="text" name="lbl" size="20" value="%s">
<input type="submit" value="Zoeken">
<input type="button" value="Wissen" onClick="javascript:formclear(form)">
</form>
`, id, html.EscapeString(first(q.r, "lbl")))

	// Maximaal 2*ZINMAX matchende xml-bestanden opvragen

	offset := 0
	o, err := strconv.Atoi(first(q.r, "offset"))
	if err == nil {
		offset = o
	}

	lbl := first(q.r, "lbl")
	query := ""
	if lbl != "" {
		query = fmt.Sprintf("WHERE `lbl` LIKE %q", lbl)
	}

	rows, err := q.db.Query(
		fmt.Sprintf(
			"SELECT `arch`,`file`,`sent`,`lbl` FROM `%s_c_%s_sent` %s LIMIT %d,%d",
			Cfg.Prefix,
			id,
			query,
			offset,
			2*ZINMAX))
	if doErr(q, err) {
		return
	}

	zinnen := make([]ZinArchFile, 0, 2*ZINMAX)

	nzin := 0
	for rows.Next() {
		nzin++
		var arch, file int
		var sent, lbl string
		err := rows.Scan(&arch, &file, &sent, &lbl)
		if doErr(q, err) {
			return
		}
		zinnen = append(zinnen, ZinArchFile{zin: sent, arch: arch, file: file})
	}

	fmt.Fprintln(q.w, "<p>\n<table class=\"corpora\">\n<tr><th>Label<th>Zin</tr>")
	for i, zin := range zinnen {
		rows, err := q.db.Query(
			fmt.Sprintf(
				"SELECT `lbl` FROM `%s_c_%s_sent` WHERE `file` = %d AND `arch` = %d", Cfg.Prefix, id, zin.file, zin.arch))
		if err == nil && rows.Next() {
			var lbl string
			rows.Scan(&lbl)
			rows.Close()

			eo := "even"
			if i%2 == 0 {
				eo = "odd"
			}
			if i == 0 {
				eo += " first"
			}
			if i == len(zinnen)-1 {
				eo += " last"
			}
			fmt.Fprintf(q.w, "<tr class=\"%s\"><td class=\"first odd\"><a href=\"tree?db=%s&amp;arch=%d&amp;file=%d\">%s</a><td class=\"even\">%s\n",
				eo, id, zin.arch, zin.file,
				html.EscapeString(lbl), html.EscapeString(zin.zin))
		} else {
			doErr(q, fmt.Errorf("Database error"))
		}
	}
	fmt.Fprint(q.w, "</table>\n<p>\n")

	// Links naar volgende en vorige pagina's met resultaten
	qs := fmt.Sprintf("id=%s&amp;lbl=%s", urlencode(id), urlencode(lbl))
	if offset > 0 || nzin == 2*ZINMAX {
		if offset > 0 {
			fmt.Fprintf(q.w, "<a href=\"/browse?%s&amp;offset=%d\">vorige</a>", qs, offset-2*ZINMAX)
		} else {
			fmt.Fprint(q.w, "vorige")
		}
		fmt.Fprint(q.w, " | ")
		if nzin == 2*ZINMAX {
			fmt.Fprintf(q.w, "<a href=\"/browse?%s&amp;offset=%d\">volgende</a>", qs, offset+2*ZINMAX)
		} else {
			fmt.Fprint(q.w, "volgende")
		}
	}

	fmt.Fprint(q.w, "</body>\n</html>\n")

}
