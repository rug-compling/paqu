package main

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func savez(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	writeHead(q, "Nieuw corpus maken", 0)

	fmt.Fprintln(q.w, "<form action=\"savez2\">\nKies een of meer corpora:\n<p>\n")
	choice := make(map[string]string)
	for _, c := range q.r.Form["db"] {
		choice[c] = " checked"
	}
	for _, c := range q.opt_db {
		if c[0] != 'B' {
			continue
		}
		p := strings.Fields(c)
		txt := strings.Join(p[1:], " ")
		opt := p[0][1:]
		fmt.Fprintf(q.w, "<input type=\"checkbox\" name=\"db\" value=\"%s\"%s>%s<br>\n",
			opt,
			choice[opt],
			html.EscapeString(txt))
	}

	fmt.Fprintf(q.w, `
<p>
Zoekopdracht:
<p>
<table>
<tr>
<td style="background-color: yellow">woord
<td>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
<td style="background-color: lightgreen">hoofdwoord
<tr><td><em>%s</em><td>%s<td><em>%s</em>
<tr><td>%s<td><td>%s
</table>
<p>
Titel:<br>
<input type="text" name="title" size="80" maxlength="64">
`,
		html.EscapeString(first(q.r, "word")),
		html.EscapeString(first(q.r, "rel")),
		html.EscapeString(first(q.r, "hword")),
		html.EscapeString(first(q.r, "postag")),
		html.EscapeString(first(q.r, "hpostag")))

	fmt.Fprintf(q.w, `<p>
<input type="hidden" name="word" value="%s">
<input type="hidden" name="postag" value="%s">
<input type="hidden" name="rel" value="%s">
<input type="hidden" name="hpostag" value="%s">
<input type="hidden" name="hword" value="%s">
<input type="submit" value="nieuw corpus maken">
</form>
`,
		html.EscapeString(first(q.r, "word")),
		html.EscapeString(first(q.r, "postag")),
		html.EscapeString(first(q.r, "rel")),
		html.EscapeString(first(q.r, "hpostag")),
		html.EscapeString(first(q.r, "hword")))

	fmt.Fprint(q.w, `
</body>
</html>
`)

}

func savez2(q *Context) {

	var fpz, fpgz *os.File
	var z *zip.Writer
	var gz *gzip.Reader

	close := func() {
		if z != nil {
			z.Close()
		}
		if fpz != nil {
			fpz.Close()
		}
		if gz != nil {
			gz.Close()
		}
		if fpgz != nil {
			fpgz.Close()
		}
	}

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	corpora := q.r.Form["db"]
	for _, corpus := range corpora {
		if !q.myprefixes[corpus] {
			http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
			return
		}
	}

	if len(corpora) == 0 {
		writeHtml(q, "Fout", "Geen corpora gekozen")
		return
	}

	word := first(q.r, "word")
	rel := first(q.r, "rel")
	hword := first(q.r, "hword")
	postag := first(q.r, "postag")
	hpostag := first(q.r, "hpostag")
	if word == "" && hword == "" && rel == "" && postag == "" && hpostag == "" {
		writeHtml(q, "Fout", "Zoektermen ontbreken")
		return
	}

	title := maxtitlelen(first(q.r, "title"))
	if title == "" {
		writeHtml(q, "Fout", "Titel ontbreekt")
		return
	}

	writeHead(q, "Nieuw corpus opslaan", 0)

	if title == "" {
		http.Error(q.w, "Titel ontbreekt", http.StatusPreconditionFailed)
		return
	}

	dirname := reNoAz.ReplaceAllString(strings.ToLower(title), "")
	if len(dirname) > 20 {
		dirname = dirname[:20]
	} else if dirname == "" {
		dirname = "a"
	}

	dirnameLock.Lock()
	defer dirnameLock.Unlock()
	for i := 0; true; i++ {
		d := dirname + abc(i)
		rows, err := q.db.Query(fmt.Sprintf("SELECT 1 FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, d))
		if doErr(q, err) {
			return
		}
		if rows.Next() {
			rows.Close()
			continue
		}
		dirname = d
		break
	}
	fulldirname := path.Join(paqudir, "data", dirname)
	err := os.Mkdir(fulldirname, 0700)
	if doErr(q, err) {
		return
	}

	fpz, err = os.Create(fulldirname + "/data")
	if doErr(q, err) {
		return
	}
	z = zip.NewWriter(fpz)

	chClose := make(<-chan bool)
	for _, prefix := range corpora {
		pathlen := len(path.Join(paqudir, "data", prefix, "xml")) + 1

		query, err := makeQuery(q, prefix, chClose)
		if doErr(q, err) {
			close()
			return
		}
		query = strings.Replace(query, " `", " `c`.`", -1)
		query = fmt.Sprintf("SELECT DISTINCT `f`.`file` FROM `%s_c_%s_deprel` `c`, `%s_c_%s_file` `f` WHERE `c`.%s AND `f`.`id` = `c`.`file`",
			Cfg.Prefix, prefix,
			Cfg.Prefix, prefix,
			query)
		rows, err := q.db.Query(query)
		if doErr(q, err) {
			close()
			return
		}
		var filename string
		for rows.Next() {
			err := rows.Scan(&filename)
			if doErr(q, err) {
				close()
				return
			}

			fpgz, err = os.Open(filename + ".gz")
			if doErr(q, err) {
				close()
				return
			}
			gz, err = gzip.NewReader(fpgz)
			if doErr(q, err) {
				close()
				return
			}
			data, err := ioutil.ReadAll(gz)
			if doErr(q, err) {
				close()
				return
			}
			gz.Close()
			gz = nil
			fpgz.Close()
			fpgz = nil

			newfile := filename[pathlen:]
			if strings.Contains(q.params[prefix], "-lbl") || q.params[prefix] == "folia" || q.params[prefix] == "tei" {
				newfile = decode_filename(newfile[10:])
			} else if q.params[prefix] == "xmlzip" || q.params[prefix] == "dact" {
				newfile = decode_filename(newfile[5:])
			}
			if len(corpora) > 1 {
				newfile = prefix + "/" + newfile
			}

			f, err := z.Create(newfile)
			if doErr(q, err) {
				close()
				return
			}
			_, err = f.Write(data)
			if doErr(q, err) {
				close()
				return
			}

		}
		err = rows.Err()
		if doErr(q, err) {
			close()
			return
		}

	}

	if doErr(q, z.Close()) {
		z = nil
		close()
		return
	}
	fpz.Close()

	newCorpus(q, dirname, title, "xmlzip")
}
