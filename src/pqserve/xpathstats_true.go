// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"
)

type AttrItem struct {
	n int
	a string
}

type AttrItems []*AttrItem

type ValueItem struct {
	begin int
	end   int
	value string
}

type ValueItems []*ValueItem

type Alpino_ds_full_node struct {
	XMLName xml.Name  `xml:"alpino_ds"`
	Node0   *FullNode `xml:"node"`
}

func getDeepAttr(attr string, n *Node, values *[]*ValueItem) {
	if n.Index != "" && len(n.NodeList) == 0 && n.Word == "" {
		*values = append(*values, &ValueItem{-1, -1, n.Index})
		return
	}
	if s := getAttr(attr, &n.FullNode); s != "" {
		*values = append(*values, &ValueItem{n.Begin, n.End, s})
		return
	}
	for _, n2 := range n.NodeList {
		getDeepAttr(attr, n2, values)
	}
}

func getIndexValue(idx, attr string, n *Node, values *[]*ValueItem) bool {
	if n.Index == idx && (len(n.NodeList) > 0 || n.Word != "") {
		getDeepAttr(attr, n, values)
		return true
	}
	for _, n2 := range n.NodeList {
		if getIndexValue(idx, attr, n2, values) {
			return true
		}
	}
	return false
}

func getFullAttr(attr string, n, top *Node) string {

	if n == nil || top == nil {
		return ""
	}

	if s := getAttr(attr, &n.FullNode); s != "" {
		return s
	}
	values := make([]*ValueItem, 0)
	getDeepAttr(attr, n, &values)

	if len(values) == 0 {
		return ""
	}

	// er kunnen nog index-waardes aan 'values' worden toegevoegd, dus geen range gebruiken
	for i := 0; i < len(values); i++ {
		v := values[i]
		if v.begin < 0 {
			getIndexValue(v.value, attr, top, &values)
		}
	}

	sort.Sort(ValueItems(values))

	s := make([]string, 0, len(values))
	p := -1
	q := 0
	for _, v := range values {
		if v.begin > p {
			p = v.begin
			if v.begin > q && len(s) > 0 {
				s = append(s, "...")
			}
			q = v.end
			s = append(s, v.value)
		}
	}
	return "  " + strings.Join(s, " ")
}

func xpathstats(q *Context) {

	download := false
	if first(q.r, "d") != "" {
		download = true
	}

	if download {
		contentType(q, "text/plain; charset=utf-8")
		cache(q)
	} else {
		contentType(q, "text/html; charset=utf-8")
		cache(q)
		fmt.Fprint(q.w, `<!DOCTYPE html>
<html>
<head>
<title></title>
<script type="text/javascript"><!--
function f(s) {
    window.parent._fn.update(s);
}
//--></script>
</head>
<body">
`)
		updateText(q, `<img src="busy.gif" alt="aan het werk...">`)
	}

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		if download {
			fmt.Fprintf(q.w, "Invalid corpus: "+prefix)
		} else {
			updateText(q, "Invalid corpus: "+html.EscapeString(prefix))
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	attr := make([]string, 3)
	attr[0], attr[1], attr[2] = first(q.r, "attr1"), first(q.r, "attr2"), first(q.r, "attr3")
	if attr[1] == "" {
		attr[1], attr[2] = attr[2], ""
	}
	if attr[0] == "" {
		attr[0], attr[1], attr[2] = attr[1], attr[2], ""
	}
	if attr[0] == "" {
		if download {
			fmt.Fprintln(q.w, "Geen attribuut gekozen")
		} else {
			updateText(q, "Geen attribuut gekozen")
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}

	nAttr := 0
	for i := 0; i < 3; i++ {
		if attr[i] != "" {
			nAttr = i + 1
		}
	}

	query := first(q.r, "xpath")

	if query == "" {
		if download {
			fmt.Fprintln(q.w, "Query ontbreekt")
		} else {
			updateText(q, "Query ontbreekt")
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}

	if strings.Contains(query, "%") {
		rules := getMacrosRules(q)
		query = macroKY.ReplaceAllStringFunc(query, func(s string) string {
			return rules[s[1:len(s)-1]]
		})
	}

	var owner string
	var nlines uint64
	rows, err := q.db.Query(fmt.Sprintf("SELECT `owner`,`nline` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if err != nil {
		updateError(q, err, !download)
		logerr(err)
		return
	}
	for rows.Next() {
		if err := rows.Scan(&owner, &nlines); err != nil {
			updateError(q, err, !download)
			logerr(err)
			rows.Close()
			return
		}
	}
	if err := rows.Err(); err != nil {
		updateError(q, err, !download)
		logerr(err)
		return
	}

	dactfiles := make([]string, 0)
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, path.Join(paqudir, "data", prefix, "data.dact"))
	} else {
		rows, err := q.db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
		if err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}
		for rows.Next() {
			var s string
			err := rows.Scan(&s)
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				rows.Close()
				return
			}
			if strings.HasSuffix(s, ".dact") {
				dactfiles = append(dactfiles, s)
			}
		}
		if err := rows.Err(); err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}
	}

	if len(dactfiles) == 0 {
		fmt.Fprintln(q.w, "Er zijn geen dact-bestanden voor dit corpus")
		return
	}

	now := time.Now()
	now2 := time.Now()

	queryparts := strings.Split(query, "+|+")

	sums := make(map[string]int)
	count := 0
	tooMany := false
	var seen uint64
	for _, dactfile := range dactfiles {
		if !download && time.Now().Sub(now2) > 2*time.Second {
			xpathout(q, sums, attr, count, tooMany, now, download, seen, nlines, false)
			now2 = time.Now()
		}
		if tooMany {
			break
		}
		select {
		case <-chClose:
			logerr(errConnectionClosed)
			return
		default:
		}
		db, err := dbxml.Open(dactfile)
		if err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}

		qu, err := db.Prepare(queryparts[0])
		if err != nil {
			updateError(q, err, !download)
			db.Close()
			return
		}
		done := make(chan bool, 1)
		interrupted := make(chan bool, 1)
		go func() {
			select {
			case <-chClose:
				interrupted <- true
				logerr(errConnectionClosed)
				qu.Cancel()
			case <-done:
			}
		}()

		docs, err := qu.Run()
		if err != nil {
			updateError(q, err, !download)
			db.Close()
			return
		}
		filename := ""
	NEXTDOC:
		for docs.Next() {
			if !download && time.Now().Sub(now2) > 2*time.Second {
				xpathout(q, sums, attr, count, tooMany, now, download, seen, nlines, false)
				now2 = time.Now()
			}

			matches := make([]string, 0)
			if len(queryparts) == 1 {
				matches = append(matches, docs.Match())
			} else {
				name := docs.Name()
				if name == filename {
					continue
				}
				filename = name
				doctxt := fmt.Sprintf("[dbxml:metadata('dbxml:name')=%q]", name)
				for i := 1; i < len(queryparts)-1; i++ {
					docs2, err := db.Query(doctxt + queryparts[i])
					if err != nil {
						updateError(q, err, !download)
						logerr(err)
						docs.Close()
						db.Close()
						return
					}
					if !docs2.Next() {
						continue NEXTDOC
					}
					docs2.Close()
				}

				docs2, err := db.Query(doctxt + queryparts[len(queryparts)-1])
				if err != nil {
					updateError(q, err, !download)
					logerr(err)
					docs.Close()
					db.Close()
					return
				}
				for docs2.Next() {
					matches = append(matches, docs2.Match())
				}
			}

			if len(matches) == 0 {
				continue
			}

			alpino := Alpino_ds{}
			err := xml.Unmarshal([]byte(docs.Content()), &alpino)
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				docs.Close()
				db.Close()
				return
			}
			for _, match := range matches {
				alp := Alpino_ds{}
				err = xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<alpino_ds version="1.3">
`+match+`
</alpino_ds>`), &alp)
				if err != nil {
					updateError(q, err, !download)
					logerr(err)
					docs.Close()
					db.Close()
					return
				}
				switch nAttr {
				case 1:
					sums[getFullAttr(attr[0], alp.Node0, alpino.Node0)]++
				case 2:
					sums[getFullAttr(attr[0], alp.Node0, alpino.Node0)+"\t"+
						getFullAttr(attr[1], alp.Node0, alpino.Node0)]++
				case 3:
					sums[getFullAttr(attr[0], alp.Node0, alpino.Node0)+"\t"+
						getFullAttr(attr[1], alp.Node0, alpino.Node0)+"\t"+
						getFullAttr(attr[2], alp.Node0, alpino.Node0)]++
				}
				count++
				if len(sums) >= 100000 {
					tooMany = true
					docs.Close()
				}
			}
		}
		if err := docs.Error(); err != nil {
			logerr(err)
		}
		if n, err := db.Size(); err == nil {
			seen += n
		}
		db.Close()
		done <- true
		select {
		case <-interrupted:
			return
		default:
		}
	}

	xpathout(q, sums, attr, count, tooMany, now, download, 0, 0, true)

	if !download {
		fmt.Fprintln(q.w, "</body>\n</html>")
	}

}

func xpathout(q *Context, sums map[string]int, attr []string, count int, tooMany bool, now time.Time, download bool, seen, total uint64, final bool) {

	attrList := make([]*AttrItem, 0, len(sums))
	for key, value := range sums {
		attrList = append(attrList, &AttrItem{value, key})
	}
	sort.Sort(AttrItems(attrList))

	var buf bytes.Buffer

	if download {
		if tooMany {
			fmt.Fprintln(q.w, "# ONDERBROKEN VANWEGE TE VEEL COMBINATIES")
		}
		fmt.Fprintf(q.w, "# %d matches in %d combinaties\n", count, len(attrList))
	} else {
		if tooMany {
			fmt.Fprintln(&buf, "<div class=\"warning\">Onderbroken vanwege te veel combinaties</div>")
		}
		f := ""
		if !final {
			f = `<img src="busy.gif" alt="aan het werk...">`
			if seen > 0 {
				f = fmt.Sprintf("%s %.1f%%", f, float64(seen)*100/float64(total))
			}
		}
		fmt.Fprintf(&buf, `<table>
<tr><td>Matches:<td class="right">%s<td rowspan="3">%s
<tr><td>Combinaties:<td class="right">%s
<tr><td>Tijd:<td class="right">%s
</table>
<p>
`, iformat(count), f, iformat(len(attrList)), tijd(time.Now().Sub(now)),
		)
	}

	nAttr := 0
	for i := 0; i < 3; i++ {
		if attr[i] != "" {
			nAttr = i + 1
		}
	}

	if download {
		fmt.Fprint(q.w, "aantal\tperc.")
	} else {
		fmt.Fprint(&buf, "<table class=\"breed\">\n<tr class=\"odd\"><th><th>")
	}
	for i := 0; i < nAttr; i++ {
		if download {
			fmt.Fprintf(q.w, "\t%s", attr[i])
		} else {
			fmt.Fprintf(&buf, "<th>%s", html.EscapeString(attr[i]))
		}
	}
	fmt.Fprintln(q.w)

	for n, a := range attrList {
		o := ""
		if n%2 == 1 {
			o = " class=\"odd\""
		}
		if !download && n == WRDMAX {
			fmt.Fprintf(&buf, "<tr%s><td class=\"right\">...<td class=\"right\">...\n", o)
			for i := 0; i < nAttr; i++ {
				fmt.Fprintln(&buf, "<td class=\"nil\">...")
			}
			break
		}
		if download {
			fmt.Fprintf(q.w, "%d\t%.1f%%", a.n, float64(a.n)/float64(count)*100)
		} else {
			fmt.Fprintf(&buf, "<tr%s><td class=\"right\">%d<td class=\"right\">%.1f%%\n", o, a.n, float64(a.n)/float64(count)*100)
		}
		v := strings.Split(a.a, "\t")
		for i := 0; i < nAttr; i++ {
			if strings.TrimSpace(v[i]) == "" {
				if download {
					fmt.Fprintf(q.w, "\tNIL")
				} else {
					fmt.Fprintln(&buf, "<td class=\"nil\">(leeg)")
				}
			} else {
				if download {
					fmt.Fprintf(q.w, "\t%s", strings.TrimSpace(v[i]))
				} else {
					if strings.HasPrefix(v[i], "  ") {
						fmt.Fprintf(&buf, "<td class=\"multi\">%s\n", html.EscapeString(strings.TrimSpace(v[i])))
					} else {
						fmt.Fprintf(&buf, "<td>%s\n", html.EscapeString(v[i]))
					}
				}
			}
		}
		if download {
			fmt.Fprintln(q.w)
		}
	}

	if !download {
		fmt.Fprintf(&buf,
			"</table>\n<hr><a href=\"xpathstats?%s&amp;d=1\" target=\"_blank\">download</a>\n",
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
		updateText(q, buf.String())
	}
}

func (x AttrItems) Less(i, j int) bool {
	if x[i].n != x[j].n {
		return x[i].n > x[j].n
	}
	return x[i].a < x[j].a
}

func (x AttrItems) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x AttrItems) Len() int {
	return len(x)
}

func (x ValueItems) Less(i, j int) bool {
	if x[i].begin < x[j].begin {
		return true
	}
	return x[i].end < x[j].end
}

func (x ValueItems) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ValueItems) Len() int {
	return len(x)
}
