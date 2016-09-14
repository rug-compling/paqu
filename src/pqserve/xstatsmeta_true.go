// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"math"
	"net/http"
	"path/filepath"
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

	var errval error
	var download bool
	var db *dbxml.Db
	var docs, docs2 *dbxml.Docs
	defer func() {
		if docs2 != nil {
			docs2.Close()
		}
		if docs != nil {
			docs.Close()
		}
		if db != nil {
			db.Close()
		}
		if errval != nil {
			updateError(q, errval, !download)
		}
		completedmeta(q, download)
		if !download {
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
	}()

	var rows *sql.Rows

	now := time.Now()
	now2 := time.Now()

	if first(q.r, "d") != "" {
		download = true
	}

	itemselect := first(q.r, "item")

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
function setvalue(n) {
    window.parent._fn.setmetaval(n);
}
function setmetavars(idx, lbl, fl, max, ac, bc) {
    window.parent._fn.setmetavars(idx, lbl, fl, max, ac, bc);
}
function setmetalines(idx, a, b) {
    window.parent._fn.setmetalines(idx, a, b);
}
function makemetatable(idx) {
    window.parent._fn.makemetatable(idx);
}
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
<body>
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
		}
		return
	}

	query := first(q.r, "xpath")

	if query == "" {
		if download {
			fmt.Fprintln(q.w, "Query ontbreekt")
		} else {
			updateText(q, "Query ontbreekt")
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
	rows, errval = q.db.Query(fmt.Sprintf("SELECT `owner`,`nline` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, prefix))
	if logerr(errval) {
		return
	}
	for rows.Next() {
		errval = rows.Scan(&owner, &nlines)
		if logerr(errval) {
			rows.Close()
			return
		}
	}
	errval = rows.Err()
	if logerr(errval) {
		return
	}

	dactfiles := make([]string, 0)
	if strings.Contains(owner, "@") {
		dactfiles = append(dactfiles, filepath.Join(paqudir, "data", prefix, "data.dact"))
	} else {
		rows, errval = q.db.Query(fmt.Sprintf("SELECT `arch` FROM `%s_c_%s_arch` ORDER BY `id`", Cfg.Prefix, prefix))
		if logerr(errval) {
			return
		}
		for rows.Next() {
			var s string
			errval = rows.Scan(&s)
			if logerr(errval) {
				rows.Close()
				return
			}
			if strings.HasSuffix(s, ".dact") {
				dactfiles = append(dactfiles, s)
			}
		}
		errval = rows.Err()
		if logerr(errval) {
			return
		}
	}

	if len(dactfiles) == 0 {
		if download {
			fmt.Fprintln(q.w, "Er zijn geen dact-bestanden voor dit corpus")
		} else {
			updateText(q, "Er zijn geen dact-bestanden voor dit corpus")
		}
		return
	}

	if !q.hasmeta[prefix] {
		if download {
			fmt.Fprintln(q.w, "Geen metadata voor dit corpus")
		} else {
			updateText(q, "Geen metadata voor dit corpus")
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
			rows, errval = q.db.Query(fmt.Sprintf(
				"SELECT `idx`,`text` FROM `%s_c_%s_mval` WHERE `id` = %d",
				Cfg.Prefix, prefix, m.id))
			if logerr(errval) {
				return
			}
			for rows.Next() {
				var t string
				var i int
				errval = rows.Scan(&i, &t)
				if logerr(errval) {
					rows.Close()
					return
				}
				tranges[m.name][t] = i
			}
			errval = rows.Err()
			if logerr(errval) {
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
		errval = row.Scan(&indexed, &size, &dmin, &dmax, &dtype, &fmin, &fstep, &imin, &istep)
		if logerr(errval) {
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
	} // for _, m := range metas

	queryparts := strings.Split(query, "+|+")

	telling := make(map[string]map[string][3]int)
	for _, m := range metas {
		telling[m.name] = make(map[string][3]int)
	}

	seen := make(map[string]bool)

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

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
		db, errval = dbxml.Open(dactfile)
		if logerr(errval) {
			return
		}

		var qu *dbxml.Query
		qu, errval = db.Prepare(queryparts[0])
		if logerr(errval) {
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

		docs, errval = qu.Run()
		if logerr(errval) {
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
				filename = name
				doctxt := fmt.Sprintf("[dbxml:metadata('dbxml:name')=%q]", filename)
				for i := 1; i < len(queryparts)-1; i++ {
					docs2, errval = db.Query(doctxt + queryparts[i])
					if logerr(errval) {
						return
					}
					if !docs2.Next() {
						docs2 = nil
						continue NEXTDOC
					}
					docs2.Close()
					docs2 = nil
				}

				docs2, errval = db.Query(doctxt + queryparts[len(queryparts)-1])
				if logerr(errval) {
					return
				}
				for docs2.Next() {
					matches++
				}
				docs2 = nil
			}

			if matches == 0 {
				continue
			}

			counter += matches

			alpino := Alpino_ds_meta{}
			errval = xml.Unmarshal([]byte(docs.Content()), &alpino)
			if logerr(errval) {
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

		} // for docs.Next()
		errval = docs.Error()
		docs = nil
		if logerr(errval) {
			return
		}
		db.Close()
		db = nil
		done <- true
		select {
		case <-interrupted:
			return
		default:
		}
	} // for _, dactfile := range dactfiles
	if !download {
		updateCount(q, counter, len(seen))
	}

	var buf bytes.Buffer

	pow10 := math.Pow10(int(math.Log10(float64(q.lines[prefix])) + .5))
	if pow10 < 10 {
		pow10 = 10
	}

	if !download {
		fmt.Fprintf(q.w, `<script type="text/javascript">
setvalue(%d);
</script>
`, int(pow10))
	} else {
		fmt.Fprintf(q.w, "# items: %d\n# zinnen: %d\n# n = %d\n", counter, len(seen), int(pow10))
	}

	for number, meta := range metas {
		items := make([]*MetaItem, 0, len(telling[meta.name]))
		for name := range telling[meta.name] {
			items = append(items, &MetaItem{
				text:  name,
				idx:   telling[meta.name][name][0],
				count: [2]int{telling[meta.name][name][1], telling[meta.name][name][2]},
			})
		}
		rows, errval = q.db.Query(fmt.Sprintf(
			"SELECT `idx`, `text`, `n` FROM `%s_c_%s_mval` WHERE `id`=%d ORDER BY `idx`",
			Cfg.Prefix, prefix, metai[meta.name]))
		if logerr(errval) {
			return
		}
		nn := make(map[int]int)
		values := make([]StructIS, 0)
		for rows.Next() {
			var idx, n int
			var txt string
			errval = rows.Scan(&idx, &txt, &n)
			if logerr(errval) {
				rows.Close()
				return
			}
			nn[idx] = n
			values = append(values, StructIS{idx, txt})
		}
		errval = rows.Err()
		if logerr(errval) {
			return
		}

		if !download {
			var hide string
			if itemselect != meta.name {
				hide = " hide"
			}
			var hex string
			for _, c := range meta.name {
				hex += fmt.Sprintf("%04x", uint16(c))
			}
			fmt.Fprintf(&buf, `
<div class="metasub%s" id="meta%s">
<p>
<b>%s</b> &mdash; <a href="javascript:void(0)" onclick="javascript:metahelp()">toelichting bij tabel</a>
<p>
<table>
  <tr>
   <td>per item:
     <table class="right" id="meta%da">
     </table>
   <td class="next">per zin:
     <table class="right" id="meta%db">
     </table>
</table>
</div>
`, hide, hex, html.EscapeString(meta.name), number, number)
			updateText(q, buf.String())
			buf.Reset()

			fl := "right"
			max := 99999
			ac := 1
			bc := 2
			if meta.mtype == "TEXT" {
				fl = "left"
				max = METAMAX
				ac = 0
				bc = 0
			}
			fmt.Fprintf(q.w, `<script type="text/javascript">
setmetavars(%d,"%s","%s",%d,%d,%d);
setmetalines(%d`, number, meta.value, fl, max, ac, bc, number)
		}

		if metat[meta.name] != "TEXT" {
			sort.Sort(MetaItems(items))
		}
		for run := 0; run < 2; run++ {
			if !download {
				fmt.Fprint(q.w, ",[")
			}
			if metat[meta.name] == "TEXT" {
				if run == 0 {
					sort.Sort(MetaItems0(items))
				} else {
					sort.Sort(MetaItems1(items))
				}
			}

			lines := make([]Statline, 0)

			if download {
				if run == 0 {
					fmt.Fprintln(q.w, "# "+meta.name+" per item\t")
				} else {
					fmt.Fprintln(q.w, "# "+meta.name+" per zin\t")
				}
			}
			select {
			case <-chClose:
				logerr(errConnectionClosed)
				return
			default:
			}

			seen := make(map[int]*Statline)
			for _, item := range items {
				lines = append(lines, Statline{item.text, item.count[run], nn[item.idx], item.idx})
				seen[item.idx] = &lines[len(lines)-1]
			}
			if download || (meta.mtype != "TEXT" && len(seen)*NEEDALL > len(values)) {
				// ontbrekende waardes (count==0) toevoegen
				if meta.mtype == "TEXT" {
					for _, v := range values {
						if _, ok := seen[v.i]; !ok {
							lines = append(lines, Statline{v.s, 0, 1, v.i})
						}
					}
				} else {
					lines2 := make([]Statline, len(values))
					for i, v := range values {
						if s, ok := seen[v.i]; ok {
							lines2[i] = *s
						} else {
							lines2[i] = Statline{v.s, 0, 1, v.i}
						}
					}
					lines = lines2
				}
			}
			p := "\n"
			for _, line := range lines {
				if download {
					if run == 1 {
						v := int(.5 + pow10*float64(line.i)/float64(line.n))
						fmt.Fprintf(q.w, "%d\t%d\t%s\n", line.i, v, line.s)
					} else {
						fmt.Fprintf(q.w, "%d\t%s\n", line.i, line.s)
					}
				} else {
					fmt.Fprintf(q.w, "%s[%d,", p, line.i)
					p = ",\n"
					if run == 1 {
						v := int(.5 + pow10*float64(line.i)/float64(line.n))
						fmt.Fprintf(q.w, "%d,", v)
					}
					fmt.Fprintf(q.w, "%d,\"%s\"]", line.idx, line.s)
				}
			} // for _, line := range lines
			if !download {
				fmt.Fprintln(q.w, "]")
			}
		} // for run := 0; run < 2; run++

		if !download {
			fmt.Fprintf(q.w, `);
			   makemetatable(%d);
			   //--></script>
			   `, number)
		}
	} // for number, meta := range metas

	if !download {
		fmt.Fprintf(&buf,
			"<hr>tijd: %s\n<p>\n<a href=\"xstatsmeta?%s&amp;d=1\">download</a>\n",
			tijd(time.Now().Sub(now)),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
		updateText(q, buf.String())
		buf.Reset()
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

func completedmeta(q *Context, download bool) {
	if download {
		return
	}
	fmt.Fprintf(q.w, `<script type="text/javascript">
window.parent._fn.completedmeta();
</script>
`)
}
