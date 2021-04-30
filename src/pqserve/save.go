package main

import (
	"github.com/rug-compling/paqu/internal/dir"

	"archive/zip"
	"bytes"
	"compress/gzip"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func savez(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	writeHead(q, "Nieuw corpus maken", 0)

	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
var submitted = false;
function submitter() {
    if (submitted) {
        return false;
    }
    submitted = true;
    $('#subbut').addClass('hide');
    $('#subsp').removeClass('hide');
    return true;
}
//--></script>
<form action="savez2" method="post" enctype="multipart/form-data" accept-charset="utf-8" onsubmit="javascript:return submitter()">
Kies een of meer corpora:
<p>
`)
	choice := make(map[string]string)
	choice[firstf(q.form, "db")] = " checked"
	var gr byte
	for _, c := range q.opt_db {
		if c[0] != gr {
			gr = c[0]
			var t string
			switch gr {
			case 'A':
				t = "algemene corpora &mdash; handmatig geannoteerd"
			case 'B':
				t = "algemene corpora &mdash; automatisch geannoteerd"
			case 'C':
				t = "algemene corpora"
			case 'D':
				t = "mijn corpora"
			case 'E':
				t = "corpora gedeeld door anderen"
			}
			fmt.Fprintln(q.w, "<b>&mdash;", t, "&mdash;</b><br>")
		}
		p := strings.Fields(c)
		txt := strings.Join(p[1:], " ")
		opt := p[0][1:]
		fmt.Fprintf(q.w, "<input type=\"checkbox\" name=\"db\" value=\"%s\"%s>%s<br>\n",
			opt,
			choice[opt],
			html.EscapeString(txt))
	}

	qword := firstf(q.form, "word")
	qrel := firstf(q.form, "rel")
	qhword := firstf(q.form, "hword")
	qpostag := firstf(q.form, "postag")
	qhpostag := firstf(q.form, "hpostag")
	qmeta := firstf(q.form, "meta")

	if qword != "" || qrel != "" || qhword != "" || qpostag != "" || qhpostag != "" {
		fmt.Fprintf(q.w, `
<p>
Zoekopdracht:
<p>
<table style="margin-left: 2em">
<tr>
<td style="background-color: yellow">woord
<td>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
<td style="background-color: lightgreen">hoofdwoord
<tr><td><em>%s</em><td>%s<td><em>%s</em>
<tr><td>%s<td><td>%s
</table>
`,
			html.EscapeString(qword),
			html.EscapeString(qrel),
			html.EscapeString(qhword),
			html.EscapeString(qpostag),
			html.EscapeString(qhpostag))
	}

	if qmeta != "" {

		fmt.Fprintf(q.w, `
<p>
Metadata:
<p>
<div style="margin-left:2.5em">
%s
</div>
`,
			html.EscapeString(qmeta))
	}

	s := "default: geen limiet"
	if Cfg.Maxdup > 0 {
		s = fmt.Sprintf("default is maximum van %d", Cfg.Maxdup)
	}
	fmt.Fprintf(q.w, `
<p>
Titel:<br>
<input type="text" name="title" size="80" maxlength="64">
<p>
Maximum aantal zinnen (%s):<br>
<input type="text" name="maxdup" size="8" maxlength="10">
`, s)

	fmt.Fprintf(q.w, `<p>
<input type="hidden" name="word" value="%s">
<input type="hidden" name="postag" value="%s">
<input type="hidden" name="rel" value="%s">
<input type="hidden" name="hpostag" value="%s">
<input type="hidden" name="hword" value="%s">
<input type="hidden" name="meta" value="%s">
<input type="submit" value="nieuw corpus maken" id="subbut">
<span id="subsp" class="hide">Even geduld...</span>
`,
		html.EscapeString(qword),
		html.EscapeString(qpostag),
		html.EscapeString(qrel),
		html.EscapeString(qhpostag),
		html.EscapeString(qhword),
		html.EscapeString(qmeta))

	fmt.Fprint(q.w, `
</body>
</html>
`)

}

func savez2(q *Context) {

	var fpz, fpgz *os.File
	var z *zip.Writer
	var gz *gzip.Reader
	var dact interface{}
	var err error
	var dirname, fulldirname string
	var okall bool

	var buf bytes.Buffer

	fmt.Fprint(&buf, "Corpus afgeleid van de volgende corpora:\n\n")

	defer func() {
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
		saveCloseDact(dact)
		if !okall {
			os.RemoveAll(fulldirname)
			sqlDB.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, dirname))
		}
	}()

	protected := 0

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	corpora := make([]string, 0, len(q.form.Value["db"]))
	for _, c := range q.form.Value["db"] {
		if s := strings.TrimSpace(c); s != "" {
			corpora = append(corpora, s)
		}
	}
	for _, corpus := range corpora {
		if !q.prefixes[corpus] {
			http.Error(q.w, "Geen toegang tot corpus", http.StatusUnauthorized)
			return
		}
		if q.protected[corpus] || !q.myprefixes[corpus] {
			protected = 1
		}

		fmt.Fprintf(&buf, " * %s\n", q.desc[corpus])
	}

	if len(corpora) == 0 {
		writeHtml(q, "Fout", "Geen corpora gekozen")
		return
	}

	word := firstf(q.form, "word")
	rel := firstf(q.form, "rel")
	hword := firstf(q.form, "hword")
	postag := firstf(q.form, "postag")
	hpostag := firstf(q.form, "hpostag")
	meta := firstf(q.form, "meta")
	if word == "" && hword == "" && rel == "" && postag == "" && hpostag == "" && meta == "" {
		writeHtml(q, "Fout", "Zoektermen ontbreken")
		return
	}
	fmt.Fprint(&buf, "\n\nCriteria:\n\n")
	if word != "" {
		fmt.Fprintf(&buf, " * Woord: %s\n", word)
	}
	if postag != "" {
		fmt.Fprintf(&buf, " * Woord, postag: %s\n", postag)
	}
	if rel != "" {
		fmt.Fprintf(&buf, " * Relatie: %s\n", rel)
	}
	if hword != "" {
		fmt.Fprintf(&buf, " * Hoofdwoord: %s\n", hword)
	}
	if hpostag != "" {
		fmt.Fprintf(&buf, " * Hoofdwoord, postag: %s\n", hpostag)
	}
	if meta != "" {
		fmt.Fprintf(&buf, " * Metadata: %s\n", strings.Join(strings.Fields(meta), " "))
	}

	title := maxtitlelen(firstf(q.form, "title"))
	if title == "" {
		writeHtml(q, "Fout", "Titel ontbreekt")
		return
	}

	maxdup, _ := strconv.Atoi(firstf(q.form, "maxdup"))
	if maxdup < 1 || maxdup > Cfg.Maxdup {
		maxdup = Cfg.Maxdup
	}

	dirname, fulldirname, ok := beginNewCorpus(q, title, hErr)
	if !ok {
		return
	}

	fpz, err = os.Create(fulldirname + "/data")
	if hErr(q, err) {
		fpz = nil
		return
	}
	z = zip.NewWriter(fpz)

	linecount := 0

	chClose := make(<-chan bool)
	for _, prefix := range corpora {
		if linecount == maxdup && maxdup > 0 {
			break
		}

		global, ok := isGlobal(q, prefix)
		if !ok {
			return
		}
		pathlen, ok := getPathLen(q, prefix, global, false)
		if !ok {
			return
		}

		query, joins, usererr, syserr := makeQueryF(q, prefix, "c", chClose)
		if hErr(q, syserr) {
			return
		}
		if uhErr(q, usererr) {
			return
		}

		query = fmt.Sprintf(
			"SELECT DISTINCT `f`.`file`, `c`.`arch` FROM `%s_c_%s_deprel` `c` "+
				"JOIN `%s_c_%s_file` `f` ON (`f`.`id`=`c`.`file`) %s WHERE %s",
			Cfg.Prefix, prefix,
			Cfg.Prefix, prefix,
			joins,
			query)
		rows, err := sqlDB.Query(query)
		if hErr(q, err) {
			return
		}
		currentarch := -1
		dact = nil
		var arch int
		var filename, dactname string
		for rows.Next() {
			if linecount == maxdup && maxdup > 0 {
				rows.Close()
				break
			}
			err = rows.Scan(&filename, &arch)
			if hErr(q, err) {
				rows.Close()
				return
			}
			if strings.HasPrefix(filename, "$$") {
				filename = dir.Data + "/data/" + prefix + "/xml" + filename[2:]
			}

			var data []byte
			if arch < 0 {
				fpgz, err = os.Open(filename + ".gz")
				if err == nil {
					gz, err = gzip.NewReader(fpgz)
					if hErr(q, err) {
						gz = nil
						rows.Close()
						return
					}
					data, err = ioutil.ReadAll(gz)
					if hErr(q, err) {
						rows.Close()
						return
					}
					gz.Close()
					gz = nil
					fpgz.Close()
					fpgz = nil
				} else {
					fpgz, err = os.Open(filename)
					if hErr(q, err) {
						fpgz = nil
						rows.Close()
						return
					}
					data, err = ioutil.ReadAll(fpgz)
					if hErr(q, err) {
						rows.Close()
						return
					}
					fpgz.Close()
					fpgz = nil
				}
			} else {
				if arch != currentarch {
					currentarch = arch
					saveCloseDact(dact)
					dact, dactname = saveOpenDact(q, prefix, arch)
				}
				data = saveGetDact(q, dact, filename)
			}

			var newfile string
			if arch < 0 {
				newfile = filename[pathlen:]
				if !global {
					if strings.Contains(q.params[prefix], "-lbl") || strings.HasPrefix(q.params[prefix], "folia") || strings.HasPrefix(q.params[prefix], "tei") {
						newfile = decode_filename(newfile[10:])
					} else if strings.HasPrefix(q.params[prefix], "xmlzip") || q.params[prefix] == "dact" {
						newfile = decode_filename(newfile[5:])
					}
				}
			} else {
				newfile = dactname[pathlen:] + "::" + filename
			}
			if len(corpora) > 1 {
				newfile = prefix + "/" + newfile
				data = xmlSetSource(data, prefix)
			}

			f, err := z.Create(newfile)
			if hErr(q, err) {
				rows.Close()
				return
			}
			_, err = f.Write(data)
			if hErr(q, err) {
				rows.Close()
				return
			}
			linecount++
		} // for rows.Next()
		err = rows.Err()
		if hErr(q, err) {
			return
		}
		saveCloseDact(dact)
		dact = nil
	}

	err = z.Close()
	z = nil
	if hErr(q, err) {
		return
	}
	fpz.Close()
	fpz = nil

	s := "xmlzip-d"
	if protected != 0 {
		s = "xmlzip-p"
	}
	newCorpus(q, dirname, title, buf.String(), s, protected, hErr, true)
	okall = true
}

func isGlobal(q *Context, prefix string) (global bool, ok bool) {

	rows, err := sqlDB.Query(fmt.Sprintf("SELECT `owner` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if hErr(q, err) {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if hErr(q, err) {
			return
		}
		global = !strings.Contains(s, "@")
	}
	return global, true
}

func getPathLen(q *Context, prefix string, global, archonly bool) (length int, ok bool) {

	if !global {
		return len(filepath.Join(dir.Data, "data", prefix, "xml")) + 1, true
	}

	var min, max sql.NullInt64

	if !archonly {

		rows, err := sqlDB.Query(fmt.Sprintf("SELECT min(`file`), max(`file`) FROM `%s_c_%s_sent` WHERE `arch` = -1", Cfg.Prefix, prefix))
		if hErr(q, err) {
			return
		}
		for rows.Next() {
			err = rows.Scan(&min, &max)
			if hErr(q, err) {
				rows.Close()
				return
			}
		}
		rows.Close()

	}

	files := make([]string, 0, 2)
	if !archonly && min.Valid && max.Valid {
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT `file` FROM `%s_c_%s_file` WHERE `id` = %d OR `id` = %d",
			Cfg.Prefix, prefix, min.Int64, max.Int64))
		if hErr(q, err) {
			return
		}
		for rows.Next() {
			var s string
			err = rows.Scan(&s)
			if hErr(q, err) {
				rows.Close()
				return
			}
			files = append(files, s)
		}
		rows.Close()
	} else {
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT min(`id`), max(`id`) FROM `%s_c_%s_arch`", Cfg.Prefix, prefix))
		if hErr(q, err) {
			return
		}
		for rows.Next() {
			err = rows.Scan(&min, &max)
			if hErr(q, err) {
				rows.Close()
				return
			}
		}
		rows.Close()
		rows, err = sqlDB.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` WHERE `id` = %d OR `id` = %d",
			Cfg.Prefix, prefix, min.Int64, max.Int64))
		if hErr(q, err) {
			return
		}
		for rows.Next() {
			var s string
			err = rows.Scan(&s)
			if hErr(q, err) {
				rows.Close()
				return
			}
			files = append(files, s)
		}
		rows.Close()
	}
	if len(files) == 1 {
		files = append(files, files[0])
	}
	if len(files) != 2 {
		if archonly {
			hErr(q, errors.New("Er zijn geen dactbestanden voor corpus "+prefix))
		} else {
			hErr(q, errors.New("Missing records in file or arch"))
		}
		return
	}

	a1 := strings.Split(files[0], "/")
	a2 := strings.Split(files[1], "/")
	var i int
	for i = 0; i < len(a1)-1 && i < len(a2)-1; i++ {
		if a1[i] != a2[i] {
			break
		}
	}
	return len(strings.Join(a1[:i], "/")) + 1, true
}

var (
	reMetaSrcChange  = regexp.MustCompile(`(<\s*meta\s[^>]*name\s*=\s*['"]paqu)((?:\.source)+['"][^>]*>)`)
	reMetaNulRemove  = regexp.MustCompile(`[ \t]*<\s*metadata\s*/\s*>\s*`)
	reMetaScrInsert1 = regexp.MustCompile(`[ \t]*<\s*/\s*metadata\s*>`)
	reMetaScrInsert2 = regexp.MustCompile(`[ \t]*<\s*node`)
)

func xmlSetSource(data []byte, prefix string) []byte {
	s := string(reMetaNulRemove.ReplaceAll(data, []byte{}))
	s = reMetaSrcChange.ReplaceAllString(s, "${1}.source${2}")
	loc := reMetaScrInsert1.FindStringIndex(s)
	if loc != nil {
		return []byte(
			fmt.Sprintf("%s<meta type=\"text\" name=\"paqu.source\" value=\"%s\"/>\n%s",
				s[:loc[0]], prefix, s[loc[0]:]))
	}
	loc = reMetaScrInsert2.FindStringIndex(s)
	if loc != nil {
		return []byte(
			fmt.Sprintf("%s<metadata>\n<meta type=\"text\" name=\"paqu.source\" value=\"%s\"/>\n</metadata>\n%s",
				s[:loc[0]], prefix, s[loc[0]:]))
	}
	return data
}
