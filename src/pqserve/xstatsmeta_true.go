// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"math"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MetaItem struct {
	text  string
	idx   int
	count [2]int
}

type MetaItems []*MetaItem

type MetaItems0 []*MetaItem

type MetaItems1 []*MetaItem

func xstatsmeta(q *Context) {
	now := time.Now()
	now2 := time.Now()

	download := false
	if first(q.r, "d") != "" {
		download = true
	}

	if download {
		contentType(q, "text/plain; charset=utf-8")
		q.w.Header().Set("Content-Disposition", "attachment; filename=telling.txt")
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
    window.parent._fn.updatemeta(s);
}
function f1(s) {
    window.parent._fn.updatemetatop(s);
}
function c(i, j) {
    window.parent._fn.countmeta(i, j);
}
//--></script>
</head>
<body">
<script type="text/javascript">
window.parent._fn.startedmeta();
c("0", "0");
</script>
`)
		if ff, ok := q.w.(http.Flusher); ok {
			ff.Flush()
		}
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
		if download {
			fmt.Fprintln(q.w, "Er zijn geen dact-bestanden voor dit corpus")
		} else {
			updateText(q, "Er zijn geen dact-bestanden voor dit corpus")
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}

	if !q.hasmeta[prefix] {
		if download {
			fmt.Fprintln(q.w, "Geen metadata voor dit corpus")
		} else {
			updateText(q, "Geen metadata voor dit corpus")
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
		return
	}
	metas := getMeta(q, prefix)
	metat := make(map[string]string)
	metai := make(map[string]int)
	tranges := make(map[string]map[string]int)
	dranges := make(map[string]*drange)
	franges := make(map[string]*frange)
	iranges := make(map[string]*irange)
	for _, m := range metas {
		metat[m.name] = m.mtype
		metai[m.name] = m.id
		if m.mtype == "TEXT" {
			tranges[m.name] = make(map[string]int)
			rows, err := q.db.Query(fmt.Sprintf(
				"SELECT `idx`,`text` FROM `%s_c_%s_mval` WHERE `id` = %d",
				Cfg.Prefix, prefix, m.id))
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				return
			}
			for rows.Next() {
				var t string
				var i int
				err := rows.Scan(&i, &t)
				if err != nil {
					updateError(q, err, !download)
					logerr(err)
					return
				}
				tranges[m.name][t] = i
			}
			err = rows.Err()
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				return
			}
			continue
		}
		var indexed bool
		var size, dtype, imin, istep int
		var dmin, dmax time.Time
		var fmin, fstep float64
		row := q.db.QueryRow(fmt.Sprintf(
			"SELECT `indexed`, `size`, `dmin`, `dmax`, `dtype`, `fmin`, `fstep`, `imin`, `istep` FROM `%s_c_%s_minf` WHERE `id` = %d",
			Cfg.Prefix, prefix, m.id))
		err := row.Scan(&indexed, &size, &dmin, &dmax, &dtype, &fmin, &fstep, &imin, &istep)
		if err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}
		switch m.mtype {
		case "INT":
			iranges[m.name] = oldIrange(imin, istep, size, indexed)
		case "FLOAT":
			franges[m.name] = oldFrange(fmin, fstep, size)
		case "DATE", "DATETIME":
			dranges[m.name] = oldDrange(dmin, dmax, dtype, indexed)
		}
	}

	queryparts := strings.Split(query, "+|+")

	telling := make(map[string]map[string][3]int)
	for _, m := range metas {
		telling[m.name] = make(map[string][3]int)
	}

	seen := make(map[string]bool)

	counter := 0
	for _, dactfile := range dactfiles {
		if !download && time.Now().Sub(now2) > 2*time.Second {
			updateCount(q, counter, len(seen))
			now2 = time.Now()
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
				updateCount(q, counter, len(seen))
				now2 = time.Now()
			}
			matches := 0
			if len(queryparts) == 1 {
				matches = 1
			} else {
				name := docs.Name()
				if name == filename {
					continue
				}
				filename = docs.Name()
				doctxt := fmt.Sprintf("[dbxml:metadata('dbxml:name')=%q]", filename)
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
					matches++
				}
			}

			if matches == 0 {
				continue
			}

			counter += matches

			alpino := Alpino_ds_meta{}
			err := xml.Unmarshal([]byte(docs.Content()), &alpino)
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				docs.Close()
				db.Close()
				return
			}

			values := make(map[string][]string)
			for _, m := range metas {
				values[m.name] = make([]string, 0)
			}
			for _, m := range alpino.Meta {
				values[m.Name] = append(values[m.Name], m.Value)
			}
			for _, m := range metas {
				if len(values[m.name]) == 0 {
					values[m.name] = append(values[m.name], "")
				}
			}
			f := dactfile + "\t" + docs.Name()
			c := 0
			if !seen[f] {
				seen[f] = true
				c = 1
			}
			for name := range values {
				for _, value := range values[name] {
					var idx int
					if value == "" {
						idx = 2147483647
					} else {
						switch metat[name] {
						case "TEXT":
							idx = tranges[name][value]
						case "INT":
							v, _ := strconv.Atoi(value)
							value, idx = iranges[name].value(v)
						case "FLOAT":
							v, _ := strconv.ParseFloat(value, 32) // 32 is dezelfde precisie als gebruikt door MySQL
							value, idx = franges[name].value(v)
						case "DATE":
							v, _ := time.Parse("2006-01-02", value)
							value, idx = dranges[name].value(v)
						case "DATETIME":
							v, _ := time.Parse("2006-01-02 15:04", value)
							value, idx = dranges[name].value(v)
						}
					}
					telling[name][value] = [3]int{idx, telling[name][value][1] + matches, telling[name][value][2] + c}
				}
			}

		}
		if err := docs.Error(); err != nil {
			logerr(err)
		}
		db.Close()
		done <- true
		select {
		case <-interrupted:
			return
		default:
		}
	}
	if !download {
		updateCount(q, counter, len(seen))
	}

	var buf bytes.Buffer

	pow10 := math.Pow10(int(math.Log10(float64(q.lines[prefix])) + .5))
	if pow10 < 10 {
		pow10 = 10
	}

	if !download {
		fmt.Fprintf(&buf, `
<div style="font-family:monospace">%s</div>
<p>
<a href="javascript:void(0)" onclick="javascript:metahelp()">toelichting bij tabellen</a>
<p>
<em>n</em> = %s
<p>
`, html.EscapeString(query), iformat(int(pow10)))
		updateTextTop(q, buf.String())
		buf.Reset()
	}

	for _, meta := range metas {
		items := make([]*MetaItem, 0, len(telling[meta.name]))
		for name := range telling[meta.name] {
			items = append(items, &MetaItem{
				text:  name,
				idx:   telling[meta.name][name][0],
				count: [2]int{telling[meta.name][name][1], telling[meta.name][name][2]},
			})
		}
		rows, err := q.db.Query(fmt.Sprintf(
			"SELECT `idx`, `n` FROM `%s_c_%s_mval` WHERE `id`=%d",
			Cfg.Prefix, prefix, metai[meta.name]))
		if err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}
		nn := make(map[int]int)
		for rows.Next() {
			var idx, n int
			err := rows.Scan(&idx, &n)
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				return
			}
			nn[idx] = n
		}
		err = rows.Err()
		if err != nil {
			updateError(q, err, !download)
			logerr(err)
			return
		}

		if metat[meta.name] != "TEXT" {
			sort.Sort(MetaItems(items))
		}
		for run := 0; run < 2; run++ {
			if metat[meta.name] == "TEXT" {
				if run == 0 {
					sort.Sort(MetaItems0(items))
				} else {
					sort.Sort(MetaItems1(items))
				}
			}

			lines := make([]Statline, 0)
			var count int

			if download {
				if run == 0 {
					fmt.Fprintln(q.w, "# "+meta.name+" per item\t")
				} else {
					fmt.Fprintln(q.w, "# "+meta.name+" per zin\t")
				}
			} else {
				if run == 0 {
					fmt.Fprintln(&buf, "<p><b>"+html.EscapeString(meta.name)+"</b><table><tr><td>per item:<table class=\"right\">")
				} else {
					fmt.Fprintln(&buf, "<td class=\"next\">per zin:<table class=\"right\">")
				}
			}
			select {
			case <-chClose:
				logerr(errConnectionClosed)
				return
			default:
			}

			for _, item := range items {
				lines = append(lines, Statline{item.text, item.count[run], nn[item.idx]})
			}
			for _, line := range lines {
				if download {
					if run == 1 {
						v := int(.5 + pow10*float64(line.i)/float64(line.n))
						fmt.Fprintf(q.w, "%d\t%d\t%s\n", line.i, v, line.s)
					} else {
						fmt.Fprintf(q.w, "%d\t%s\n", line.i, line.s)
					}
				} else {
					td := "<td>"
					if line.s == "" {
						td = `<td class="nil">`
					}
					if meta.mtype == "TEXT" {
						if line.s == "" {
							td = `<td class="left nil">`
						} else {
							td = `<td class="left">`
						}
						count++
					}
					fmt.Fprintln(&buf, "<tr><td>", iformat(line.i))
					if run == 1 {
						v := int(.5 + pow10*float64(line.i)/float64(line.n))
						fmt.Fprintf(&buf, "<td>%s", iformat(v))
					}
					if line.s == "" {
						line.s = "(leeg)"
					}
					fmt.Fprintln(&buf, td, html.EscapeString(line.s))
				}
				if count == METAMAX && meta.mtype == "TEXT" {
					break
				}
			}
			if !download {
				if count == METAMAX && meta.mtype == "TEXT" {
					fmt.Fprint(&buf, "<tr><td>")
					if run == 1 {
						fmt.Fprint(&buf, "<td>...<td>...")
					}
					if meta.mtype == "TEXT" {
						fmt.Fprint(&buf, "<td class=\"left\">...")
					} else {
						fmt.Fprint(&buf, "<td>...")
					}
				}
				fmt.Fprintln(&buf, "</table>")
				if run == 1 {
					fmt.Fprintln(&buf, "</table>")
					updateText(q, buf.String())
					buf.Reset()
				}
			}
		}
	}

	if !download {
		fmt.Fprintf(&buf,
			"<hr>tijd: %s\n<p>\n<a href=\"xstatsmeta?%s&amp;d=1\">download</a>\n",
			tijd(time.Now().Sub(now)),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
		updateText(q, buf.String())
		completedmeta(q, download)
	}

}

func (s MetaItems) Len() int {
	return len(s)
}

func (s MetaItems0) Len() int {
	return len(s)
}

func (s MetaItems1) Len() int {
	return len(s)
}

func (s MetaItems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s MetaItems0) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s MetaItems1) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s MetaItems) Less(i, j int) bool {
	return s[i].idx < s[j].idx
}

func (s MetaItems0) Less(i, j int) bool {
	if s[i].count[0] != s[j].count[0] {
		return s[i].count[0] > s[j].count[0]
	}
	return s[i].idx < s[j].idx
}

func (s MetaItems1) Less(i, j int) bool {
	if s[i].count[1] != s[j].count[1] {
		return s[i].count[1] > s[j].count[1]
	}
	return s[i].idx < s[j].idx
}

func updateCount(q *Context, i, j int) {
	fmt.Fprintf(q.w, `<script type="text/javascript">
c(%q,%q);
</script>
`, iformat(i), iformat(j))
	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}
}

func updateTextTop(q *Context, s string) {
	fmt.Fprintf(q.w, `<script type="text/javascript">
f1(%q);
</script>
`, s)
	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
	}
}
