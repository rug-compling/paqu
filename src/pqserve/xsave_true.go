// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"archive/zip"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func xsavez(q *Context) {

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
<form action="xsavez2" method="post" enctype="multipart/form-data" accept-charset="utf-8" onsubmit="javascript:return submitter()">
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
				t = "algemene corpora"
			case 'B':
				t = "mijn corpora"
			case 'C':
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

	fmt.Fprintf(q.w, `
<p>
XPATH query:
<p>
<pre>
%s
</pre>
<p>
Titel:<br>
<input type="text" name="title" size="80" maxlength="64">
`,
		html.EscapeString(firstf(q.form, "xpath")))

	s := "default: geen limiet"
	if Cfg.Maxdup > 0 {
		s = fmt.Sprintf("default is maximum van %d", Cfg.Maxdup)
	}
	fmt.Fprintf(q.w, `
<p>
Maximum aantal zinnen (%s):<br>
<input type="text" name="maxdup" size="8" maxlength="10">
`, s)

	fmt.Fprintf(q.w, `<p>
<input type="hidden" name="xpath" value="%s">
<input type="hidden" name="mt" value="%s">
<input type="submit" value="nieuw corpus maken" id="subbut">
<span id="subsp" class="hide">Even geduld...</span>
`,
		html.EscapeString(firstf(q.form, "xpath")), firstf(q.form, "mt"))

	fmt.Fprint(q.w, `
</body>
</html>
`)

}

func xsavez2(q *Context) {

	var fpz, fpgz *os.File
	var z *zip.Writer
	var gz *gzip.Reader
	var dact *dbxml.Db
	var docs *dbxml.Docs
	var dirname, fulldirname string
	var okall bool

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
		if docs != nil {
			docs.Close()
		}
		if dact != nil {
			dact.Close()
		}
		if !okall {
			os.RemoveAll(fulldirname)
			q.db.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, dirname))
		}
	}()

	protected := 0

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	methode := firstf(q.form, "mt")
	if methode != "dx" {
		methode = "std"
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
	}

	if len(corpora) == 0 {
		writeHtml(q, "Fout", "Geen corpora gekozen")
		return
	}

	xpath := firstf(q.form, "xpath")
	if xpath == "" {
		writeHtml(q, "Fout", "Zoekterm ontbreekt")
		return
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

	dirname, fulldirname, ok := beginNewCorpus(q, q.db, title, hErr)
	if !ok {
		return
	}

	fpz, err := os.Create(fulldirname + "/data")
	if hErr(q, err) {
		fpz = nil
		return
	}
	z = zip.NewWriter(fpz)

	linecount := 0
	for _, prefix := range corpora {
		if linecount == maxdup && maxdup > 0 {
			break
		}

		global, ok := isGlobal(q, prefix)
		if !ok {
			return
		}
		pathlen, ok := getPathLen(q, prefix, global, true)
		if !ok {
			return
		}

		dactfiles := make([]string, 0)
		if !global {
			dactfiles = append(dactfiles, fmt.Sprintf("%s/data/%s/data.dact", paqudir, prefix))
		} else {
			rows, err := q.db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
			if hErr(q, err) {
				return
			}
			for rows.Next() {
				var s string
				if hErr(q, rows.Scan(&s)) {
					rows.Close()
					return
				}
				if strings.HasSuffix(s, ".dact") {
					dactfiles = append(dactfiles, s)
				}
			}
			if hErr(q, rows.Err()) {
				return
			}
		}

		fullquery := xpath
		if strings.Contains(xpath, "%") {
			rules := getMacrosRules(q)
			fullquery = macroKY.ReplaceAllStringFunc(xpath, func(s string) string {
				return rules[s[1:len(s)-1]]
			})
		}
		queryparts := strings.Split(fullquery, "+|+")

		for _, dactfile := range dactfiles {
			if linecount == maxdup && maxdup > 0 {
				break
			}
			if Cfg.Dactx && methode == "dx" {
				dactfile += "x"
			}
			var data []byte
			dact, err = dbxml.Open(dactfile)
			if hErr(q, err) {
				dact = nil
				return
			}

			qu, err := dact.Prepare(queryparts[0])
			if hErr(q, err) {
				return
			}
			docs, err = qu.Run()
			if hErr(q, err) {
				docs = nil
				return
			}
			seen := make(map[string]bool)
		NEXTDOC:
			for docs.Next() {
				if linecount == maxdup && maxdup > 0 {
					break
				}
				filename := docs.Name()
				if seen[filename] {
					continue
				}
				seen[filename] = true
				found := false
				if len(queryparts) == 1 {
					found = true
					data = []byte(docs.Content())
				} else {
					doctxt := fmt.Sprintf("[dbxml:metadata('dbxml:name')=%q]", filename)
					for i := 1; i < len(queryparts)-1; i++ {
						docs2, err := dact.Query(doctxt + queryparts[i])
						if hErr(q, err) {
							return
						}
						if !docs2.Next() {
							continue NEXTDOC
						}
						docs2.Close()
					}
					docs2, err := dact.Query(doctxt + queryparts[len(queryparts)-1])
					if hErr(q, err) {
						return
					}
					found = false
					if docs2.Next() {
						found = true
						data = []byte(docs2.Content())
						docs2.Close()
					}

				}
				if !found {
					continue
				}

				newfile := filename
				if global {
					newfile = dactfile[pathlen:len(dactfile)-5] + "::" + filename
				}
				if len(corpora) > 1 {
					newfile = prefix + "/" + newfile
					data = xmlSetSource(data, prefix)
				}
				f, err := z.Create(newfile)
				if hErr(q, err) {
					return
				}
				if methode == "dx" {
					data, err = unexpandDact(data)
					if hErr(q, err) {
						return
					}
				}
				_, err = f.Write(data)
				if hErr(q, err) {
					return
				}
				linecount++
			} // for docs.Next()
			err = docs.Error()
			docs = nil
			if hErr(q, err) {
				return
			}
			dact.Close()
			dact = nil
		} // for range dactfiles
	} // for range corpora

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
	newCorpus(q, q.db, dirname, title, s, protected, hErr, true)
	okall = true
}

func unexpandDact(data []byte) ([]byte, error) {
	alpino := Alpino_ds_complete{}
	err := xml.Unmarshal(data, &alpino)
	if err != nil {
		return nil, err
	}
	for strings.HasPrefix(alpino.Version, "X-") {
		alpino.Version = alpino.Version[2:]
	}
	unexpandDactNode(alpino.Node0)
	b, err := xml.Marshal(alpino)
	if err != nil {
		return nil, err
	}
	return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		strings.Replace(
			strings.Replace(string(b), "  <metadata></metadata>\n", "", 1),
			"  <comments></comments>\n", "", 1) + "\n"), nil
}

func unexpandDactNode(node *Node) {
	if node.NodeList == nil {
		return
	}
	for i, n := range node.NodeList {
		if n.OtherId == "" {
			unexpandDactNode(n)
		} else {
			node.NodeList[i] = &Node{
				FullNode: FullNode{
					Begin: n.Begin,
					End:   n.End,
					Id:    n.Id,
					Index: n.Index,
					Rel:   n.Rel,
				},
			}
		}
	}
}
