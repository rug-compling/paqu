// +build !nodbxml

package main

import (
	"github.com/pebbe/dbxml"

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

func (s StructIS) String() string {
	return fmt.Sprintf("%12d%s", int64(s.i)-math.MinInt32, s.s)
}

func getDeepAttr(attr string, n *Node, values *[]*ValueItem) {
	if n.Index != "" && len(n.NodeList) == 0 && n.Word == "" {
		*values = append(*values, &ValueItem{-1, -1, n.Index})
		return
	}
	if s := strings.TrimSpace(getAttr(attr, &n.FullNode)); s != "" {
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

	is_word := attr == "is_word_"
	if is_word {
		attr = "word"
	}

	if s := strings.TrimSpace(getAttr(attr, &n.FullNode)); s != "" {
		if is_word {
			return "+"
		}
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
		if is_word {
			v.value = "+"
		}
		if v.begin > p {
			p = v.begin
			if v.begin > q && len(s) > 0 {
				s = append(s, "...")
			}
			q = v.end
			s = append(s, v.value)
		}
	}
	if len(s) == 0 {
		// alleen indexen
		return ""
	}
	return "  " + strings.Join(s, " ")

}

func xpathstats(q *Context) {

	download := false
	if first(q.r, "d") != "" {
		download = true
	}

	methode := first(q.r, "mt")
	if methode != "dx" {
		methode = "std"
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
    window.parent._fn.update(s);
}
function e(s) {
    window.parent._fn.error(s);
}
function init(s) {
    window.parent._fn.init(s);
}
//--></script>
</head>
<body>
`)
	}

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		if download {
			fmt.Fprintf(q.w, "Invalid corpus: "+prefix)
		} else {
			updateJsonErr(q, "Invalid corpus: "+html.EscapeString(prefix))
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

	attr := make([]string, 5)
	attr[0], attr[1], attr[2], attr[3], attr[4] =
		first(q.r, "attr1"), first(q.r, "attr2"), first(q.r, "attr3"), first(q.r, "attr4"), first(q.r, "attr5")
	wantRel := false
	j := 0
	for i := 0; i < 5; i++ {
		if attr[i] != "" {
			attr[j] = attr[i]
			j++
			if attr[i] == "rel" {
				wantRel = true
			}
		}
	}
	for ; j < 5; j++ {
		attr[j] = ""
	}
	/*
		if attr[0] == "" {
			if download {
				fmt.Fprintln(q.w, "Geen attribuut gekozen")
			} else {
				updateJsonErr(q, "Geen attribuut gekozen")
				fmt.Fprintln(q.w, "</body>\n</html>")
			}
			return
		}
	*/

	nAttr := 0
	for i := 0; i < 5; i++ {
		if attr[i] != "" {
			nAttr = i + 1
		}
	}

	var isMeta [5]bool
	var isInt [5]bool
	var isFloat [5]bool
	var isDate [5]bool
	var isDateTime [5]bool
	var iranges [5]*irange
	var franges [5]*frange
	var dranges [5]*drange
	var isidx [5]bool
	var aligns [5]string
	for i := 0; i < 5; i++ {
		aligns[i] = "left"
		if attr[i] != "" && attr[i][0] == ':' {
			isMeta[i] = true
			name := attr[i][1:]
			rows, err := q.db.Query(fmt.Sprintf("SELECT `type` FROM `%s_c_%s_midx` WHERE `name` = %q",
				Cfg.Prefix, prefix, name))
			if err != nil {
				updateError(q, err, !download)
				logerr(err)
				return
			}
			var t string
			for rows.Next() {
				rows.Scan(&t)
			}
			if t == "INT" {
				rows, err := q.db.Query(fmt.Sprintf(
					"SELECT MIN(`ival`), MAX(`ival`), COUNT(DISTINCT `ival`) FROM `%s_c_%s_meta` JOIN `%s_c_%s_midx` USING (`id`) WHERE `name` = %q AND `idx` != 2147483647",
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					name))
				if err != nil {
					updateError(q, err, !download)
					logerr(err)
					return
				}
				var min, max, count int
				for rows.Next() {
					if rows.Scan(&min, &max, &count) == nil {
						iranges[i] = newIrange(min, max, count)
						isInt[i] = true
						aligns[i] = "right"
						isidx[i] = true
					}
				}
			} else if t == "FLOAT" {
				rows, err := q.db.Query(fmt.Sprintf(
					"SELECT MIN(`fval`), MAX(`fval`) FROM `%s_c_%s_meta` JOIN `%s_c_%s_midx` USING (`id`) WHERE `name` = %q AND `idx` != 2147483647",
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					name))
				if err != nil {
					updateError(q, err, !download)
					logerr(err)
					return
				}
				var min, max float64
				for rows.Next() {
					if rows.Scan(&min, &max) == nil {
						franges[i] = newFrange(min, max)
						isFloat[i] = true
						aligns[i] = "right"
						isidx[i] = true
					}
				}
			} else if t == "DATE" || t == "DATETIME" {
				rows, err := q.db.Query(fmt.Sprintf(
					"SELECT MIN(`dval`), MAX(`dval`) FROM `%s_c_%s_meta` JOIN `%s_c_%s_midx` USING (`id`) WHERE `name` = %q AND `idx` != 2147483647",
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					name))
				if err != nil {
					updateError(q, err, !download)
					logerr(err)
					return
				}
				var min, max time.Time
				for rows.Next() {
					if rows.Scan(&min, &max) == nil {
						aligns[i] = "right"
						isidx[i] = true
						if t == "DATE" {
							dranges[i] = newDrange(min, max, 0, false)
							isDate[i] = true
						} else {
							dranges[i] = newDrange(min, max, 0, true)
							isDateTime[i] = true
						}
					}
				}
			}
		}
	}

	query := first(q.r, "xpath")

	if query == "" {
		if download {
			fmt.Fprintln(q.w, "Query ontbreekt")
		} else {
			updateJsonErr(q, "Query ontbreekt")
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
		dactfiles = append(dactfiles, filepath.Join(paqudir, "data", prefix, "data.dact"))
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

	if !download {
		fmt.Fprintf(q.w, `<script type="text/javascript">
init({
"download": %q,
"aligns": ["right"`,
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1)+"&amp;d=1")
		for i := 0; i < nAttr; i++ {
			fmt.Fprint(q.w, `,"`, aligns[i], `"`)
		}
		fmt.Fprint(q.w, "],\n\"labels\": [\"items\"")
		for i := 0; i < nAttr; i++ {
			a := attr[i]
			if a[0] == ':' {
				a = a[1:]
			}
			fmt.Fprint(q.w, `,"`, html.EscapeString(a), `"`)
		}
		fmt.Fprint(q.w, "],\n\"isidx\": [true")
		for i := 0; i < nAttr; i++ {
			fmt.Fprintf(q.w, ",%v", isidx[i])
		}
		fmt.Fprintln(q.w, "]});\n</script>")
		if ff, ok := q.w.(http.Flusher); ok {
			ff.Flush()
		}
	}

	now := time.Now()
	now2 := time.Now()

	queryparts := strings.Split(query, "+|+")

	sums := make(map[[5]StructIS]int)
	count := 0
	linecount := 0
	tooMany := false
	var seen uint64
	for _, dactfile := range dactfiles {
		if !download && time.Now().Sub(now2) > 2*time.Second {
			xpathout(q, sums, attr, isidx, count, linecount, tooMany, now, download, seen, nlines, false)
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
		if Cfg.Dactx && methode == "dx" {
			dactfile += "x"
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
		var seenId map[string]bool
	NEXTDOC:
		for docs.Next() {
			if !download && time.Now().Sub(now2) > 2*time.Second {
				xpathout(q, sums, attr, isidx, count, linecount, tooMany, now, download, seen, nlines, false)
				now2 = time.Now()
			}

			matches := make([]string, 0)
			if len(queryparts) == 1 {
				matches = append(matches, docs.Match())
				name := docs.Name()
				if name != filename {
					filename = name
					linecount++
					seenId = make(map[string]bool)
				}
			} else {
				name := docs.Name()
				if name == filename {
					continue
				}
				filename = name
				seenId = make(map[string]bool)
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
				first := true
				for docs2.Next() {
					matches = append(matches, docs2.Match())
					if first {
						first = false
						linecount++
					}
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

			var mm [5][]StructIS
			for i := 0; i < 5; i++ {
				mm[i] = make([]StructIS, 0, 4)
			}
			for i := 0; i < 5; i++ {
				if isMeta[i] {
					name := attr[i][1:]
					for _, m := range alpino.Meta {
						if m.Name == name {
							if isInt[i] {
								v, err := strconv.Atoi(m.Value)
								if err == nil {
									vv, idx := iranges[i].value(v)
									mm[i] = append(mm[i], StructIS{idx, vv})
								} else {
									mm[i] = append(mm[i], StructIS{math.MinInt32, err.Error()})
								}
							} else if isFloat[i] {
								v, err := strconv.ParseFloat(m.Value, 64)
								if err == nil {
									vv, idx := franges[i].value(v)
									mm[i] = append(mm[i], StructIS{idx, vv})
								} else {
									mm[i] = append(mm[i], StructIS{math.MinInt32, err.Error()})
								}
							} else if isDate[i] {
								v, err := time.Parse("2006-01-02", m.Value)
								if err == nil {
									vv, idx := dranges[i].value(v)
									mm[i] = append(mm[i], StructIS{idx, vv})
								} else {
									mm[i] = append(mm[i], StructIS{math.MinInt32, err.Error()})
								}
							} else if isDateTime[i] {
								v, err := time.Parse("2006-01-02 15:04", m.Value)
								if err == nil {
									vv, idx := dranges[i].value(v)
									mm[i] = append(mm[i], StructIS{idx, vv})
								} else {
									mm[i] = append(mm[i], StructIS{math.MinInt32, err.Error()})
								}
							} else {
								mm[i] = append(mm[i], StructIS{0, m.Value})
							}
						}
					}
				}
			}
			for i := 0; i < 5; i++ {
				if len(mm[i]) == 0 {
					mm[i] = append(mm[i], StructIS{2147483647, ""})
				}
			}

			var at [5]StructIS
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
				sid := ""
				if alp.Node0 != nil {
					sid = alp.Node0.Id
					if alp.Node0.OtherId != "" {
						sid = alp.Node0.OtherId
					}
					if wantRel {
						sid = sid + " " + alp.Node0.Rel
					}
				}
				if seenId[sid] {
					continue
				}
				seenId[sid] = true
				for _, at[0] = range mm[0] {
					for _, at[1] = range mm[1] {
						for _, at[2] = range mm[2] {
							for _, at[3] = range mm[3] {
								for _, at[4] = range mm[4] {
									if nAttr > 0 && attr[0][0] != ':' {
										at[0] = StructIS{0, getFullAttr(attr[0], alp.Node0, alpino.Node0)}
									}
									if nAttr > 1 && attr[1][0] != ':' {
										at[1] = StructIS{0, getFullAttr(attr[1], alp.Node0, alpino.Node0)}
									}
									if nAttr > 2 && attr[2][0] != ':' {
										at[2] = StructIS{0, getFullAttr(attr[2], alp.Node0, alpino.Node0)}
									}
									if nAttr > 3 && attr[3][0] != ':' {
										at[3] = StructIS{0, getFullAttr(attr[3], alp.Node0, alpino.Node0)}
									}
									if nAttr > 4 && attr[4][0] != ':' {
										at[4] = StructIS{0, getFullAttr(attr[4], alp.Node0, alpino.Node0)}
									}
									sums[at]++
									count++
									if len(sums) >= 100000 {
										tooMany = true
										docs.Close()
									}
								}
							}
						}
					}
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

	xpathout(q, sums, attr, isidx, count, linecount, tooMany, now, download, 0, 0, true)

	if !download {
		fmt.Fprintln(q.w, "</body>\n</html>")
	}

}

func xpathout(q *Context, sums map[[5]StructIS]int, attr []string, isidx [5]bool, count int, linecount int, tooMany bool, now time.Time, download bool, seen, total uint64, final bool) {
	nAttr := 0
	for i := 0; i < 5; i++ {
		if attr[i] != "" {
			nAttr = i + 1
		}
	}

	data := &StatSorter{
		lines: make([]StatLine, 0),
		isInt: make([]bool, 1+nAttr),
	}
	data.isInt[0] = true
	for i := 0; i < nAttr; i++ {
		data.isInt[i+1] = isidx[i]
	}
	for key, value := range sums {
		is := make([]StructIS, nAttr+1)
		is[0].i = value
		is[0].s = fmt.Sprintf("%.1f%%", float64(value)/float64(count)*100)
		for i := 0; i < nAttr; i++ {
			is[i+1] = key[i]
		}
		data.lines = append(data.lines, StatLine{cols: is})
	}
	if len(data.lines) > WRDMAX {
		for i := 1; i < len(data.isInt); i++ {
			data.n = i
			sort.Sort(data)
			for j := range data.lines {
				if j == WRDMAX {
					break
				}
				data.lines[j].used = true
			}
		}
	}
	data.n = 0
	sort.Sort(data)

	if download {
		if tooMany {
			fmt.Fprintln(q.w, "# ONDERBROKEN VANWEGE TE VEEL COMBINATIES")
		}
		fmt.Fprintf(q.w, "# %d zinnen met %d items in %d combinaties\n", linecount, count, len(data.lines))
	} else {
		fmt.Fprintf(q.w, `<script type="text/javascript">
f({
"toomany": %v,
"final": %v,
"matches": "%s",
"linecount": "%s",
"combis": "%s",
"tijd": "%s",
`,
			tooMany,
			final,
			iformat(count),
			iformat(linecount),
			iformat(len(data.lines)),
			tijd(time.Now().Sub(now)))
		if seen > 0 {
			fmt.Fprintf(q.w, "\"perc\": \"%.1f%%\",\n", float64(seen)*100/float64(total))
		}
	}

	if download {
		fmt.Fprint(q.w, "aantal\tperc.")
	}
	for i := 0; i < nAttr; i++ {
		a := attr[i]
		if a[0] == ':' {
			a = a[1:]
		}
		if download {
			fmt.Fprintf(q.w, "\t%s", a)
		}
	}
	if download {
		fmt.Fprintln(q.w)
	} else {
		fmt.Fprint(q.w, "\"lines\": [")
	}

	p1 := ""
	for n, a := range data.lines {
		if !download && n >= WRDMAX && !a.used {
			continue
		}
		if !download {
			fmt.Fprintf(q.w, "%s\n[", p1)
		}
		p1 = ","
		p2 := ""
		for m, c := range a.cols {
			if download {
				if m == 0 {
					fmt.Fprintf(q.w, "%d\t%s", c.i, c.s)
				} else {
					fmt.Fprintf(q.w, "\t%s", c.s)
				}
			} else {
				fmt.Fprintf(q.w, "%s[%q,%d]", p2, c.s, c.i)
				p2 = ","
			}
		}
		if !download {
			fmt.Fprint(q.w, "]")
		}
		if download {
			fmt.Fprintln(q.w)
		}
	}

	if !download {
		fmt.Fprintln(q.w, "]});\n</script>")
		if ff, ok := q.w.(http.Flusher); ok {
			ff.Flush()
		}
	}
}

func updateJsonErr(q *Context, s string) {
	fmt.Fprintf(q.w, `<script type="text/javascript">
e(%q);
</script>
`, html.EscapeString(s))
	if ff, ok := q.w.(http.Flusher); ok {
		ff.Flush()
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
