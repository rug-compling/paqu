package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"math"
	"net/http"
	"strings"
	"time"
)

type Statline struct {
	s   string
	i   int
	n   int
	idx int
}

func statsmeta(q *Context) {

	var errval error
	var download bool
	defer func() {
		if errval != nil {
			updateError(q, errval, !download)
		}
		completedmeta(q, download)
		if !download {
			fmt.Fprintln(q.w, "</body>\n</html>")
		}
	}()

	now := time.Now()

	if first(q.r, "d") != "" {
		download = true
	}

	if download {
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		q.w.Header().Set("Content-Disposition", "attachment; filename=telling.txt")
		cache(q)
	} else {
		q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
		cache(q)
		fmt.Fprint(q.w, `<!DOCTYPE html>
<html>
<head>
<title></title>
<script type="text/javascript"><!--
function f(s) {
    window.parent._fn.updatemeta(s);
}
//--></script>
</head>
<body">
<script type="text/javascript">
window.parent._fn.startedmeta();
</script>
`)
	}

	option := make(map[string]string)
	for _, t := range []string{"word", "postag", "rel", "hpostag", "hword"} {
		option[t] = first(q.r, t)
	}
	if option["word"] == "" && option["postag"] == "" && option["rel"] == "" && option["hpostag"] == "" && option["hword"] == "" {
		updateError(q, errors.New("Query ontbreekt"), !download)
		return
	}

	prefix := first(q.r, "db")
	if prefix == "" {
		updateError(q, errors.New("Geen corpus opgegeven"), !download)
		return
	}
	if !q.prefixes[prefix] {
		updateError(q, errors.New("Ongeldig corpus"), !download)
		return
	}

	metas := getMeta(q, prefix)

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	var query string
	query, errval = makeQuery(q, prefix, "", chClose)
	if logerr(errval) {
		return
	}

	pow10 := math.Pow10(int(math.Log10(float64(q.lines[prefix])) + .5))
	if pow10 < 10 {
		pow10 = 10
	}

	var buf bytes.Buffer

	// DEBUG: HTML-uitvoer van de query
	if !download {
		fmt.Fprintf(&buf, `
<div style="font-family:monospace">%s</div>
<p>
<a href="javascript:void(0)" onclick="javascript:metahelp()">toelichting bij tabellen</a>
<p>
`, html.EscapeString(query))
		updateText(q, buf.String())
		buf.Reset()
	} else {
		fmt.Fprintln(q.w, "# n =", int(pow10))
	}

	// Tellingen van onderdelen
	for _, meta := range metas {

		// telling van metadata in matchende zinnen

		values := make([]StructIS, 0)
		rows, errval := q.db.Query(fmt.Sprintf(
			"SELECT `idx`,`text` FROM `%s_c_%s_mval` WHERE `id`=%d ORDER BY `idx`",
			Cfg.Prefix, prefix,
			meta.id))
		if logerr(errval) {
			return
		}
		for rows.Next() {
			var i int
			var s string
			errval = rows.Scan(&i, &s)
			if logerr(errval) {
				rows.Close()
				return
			}
			values = append(values, StructIS{i, s})
		}
		errval = rows.Err()
		if logerr(errval) {
			return
		}

		for run := 0; run < 2; run++ {
			seen := make(map[int]*Statline)
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
					fmt.Fprintln(&buf, "<tr><td><em>aantal</em><td><em>"+meta.value+"</em>")
				} else {
					fmt.Fprintln(&buf, "<td class=\"next\">per zin:<table class=\"right\">")
					fmt.Fprintf(&buf, "<tr><td><em>aantal</em><td><em>per&nbsp;%s</em><td><em>%s</em>", iformat(int(pow10)), meta.value)
				}
			}
			select {
			case <-chClose:
				logerr(errConnectionClosed)
				return
			default:
			}
			var qu string
			var order string
			if meta.mtype == "TEXT" {
				order = "1 DESC,"
			}
			if run == 0 {
				qu = fmt.Sprintf(
					"SELECT COUNT(`text`), `idx`, `text`, 0 "+
						"FROM `%s_c_%s_deprel` "+
						"JOIN `%s_c_%s_meta` USING(`arch`,`file`) "+
						"JOIN `%s_c_%s_mval` USING (`id`,`idx`) "+
						"WHERE `id` = %d AND %s "+
						"GROUP BY `text` ORDER BY %s `idx`",
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					meta.id,
					query,
					order)
			} else {
				qu = fmt.Sprintf(
					"SELECT DISTINCT `arch`,`file`,`idx`,`text`,`n` "+
						"FROM `%s_c_%s_deprel` "+
						"JOIN `%s_c_%s_meta` USING(`arch`,`file`) "+
						"JOIN `%s_c_%s_mval` USING (`id`,`idx`) "+
						"WHERE `id` = %d AND %s",
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					Cfg.Prefix, prefix,
					meta.id,
					query)
				qu = fmt.Sprintf(
					"SELECT COUNT(`a`.`text`), `a`.`idx`, `a`.`text`,`a`.`n` FROM ( %s ) `a` GROUP BY `a`.`idx` ORDER BY %s `a`.`idx`",
					qu,
					order)
			}
			var rows *sql.Rows
			rows, errval = timeoutQuery(q, chClose, qu)
			if logerr(errval) {
				return
			}
			for rows.Next() {
				select {
				case <-chClose:
					rows.Close()
					logerr(errConnectionClosed)
					return
				default:
				}
				var idx, cnt, n int
				var text string
				errval = rows.Scan(&cnt, &idx, &text, &n)
				if logerr(errval) {
					rows.Close()
					return
				}
				if len(lines) < METAMAX || download || meta.mtype != "TEXT" {
					text = unHigh(text)
					lines = append(lines, Statline{text, cnt, n, idx})
					seen[idx] = &lines[len(lines)-1]
				}
			}
			errval = rows.Err()
			if logerr(errval) {
				return
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
			}
			if !download {
				if count == METAMAX {
					fmt.Fprint(&buf, "<tr><td>")
					if run == 1 {
						fmt.Fprint(&buf, "<td>...")
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
			"<hr>tijd: %s\n<p>\n<a href=\"statsmeta?%s&amp;d=1\">download</a>\n",
			tijd(time.Now().Sub(now)),
			strings.Replace(q.r.URL.RawQuery, "&", "&amp;", -1))
		updateText(q, buf.String())
		buf.Reset()
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

func metahelp(q *Context) {
	fmt.Fprint(q.w, `
<div class="submenu a9999" id="helpmeta">
<div class="corpushelp">

De tabellen bestaan uit twee delen. Het linkerdeel,
<em>per item</em>, geeft het aantal matches per metadata-waarde. Dit
is het totaal aantal matches in het corpus, en dat
kan soms hoger zijn dan het aantal matchende zinnen
omdat er soms binnen één zin twee of meer matches zijn.
<p>
De rechterdeel, <em>per zin</em>, geeft het aantal zinnen waarin een
match gevonden is. Dit aantal staat in de eerste kolom.
De tweede kolom binnen de tabel <em>per zin</em> geeft aan wat
de relatieve frequentie is voor deze metadata-waarde.
Dit is het aantal hits per <em>n</em> zinnen waarbij <em>n</em>
bijvoorbeeld 10&#8239;000 of 100&#8239;000 is (afhankelijk van de
grootte van het corpus).
<p>
Voorbeeld, sekse per zin:
<p>
<table class="right">
<tr><td><em>aantal</em><td><em>per&nbsp;10&#8239;000</em><td><em>waarde</em>
<tr><td>36<td>40<td class="left"> female
<tr><td>30<td>50<td class="left"> male
</table>
<p>
In dit voorbeeld zijn er meer zinnen die matchen voor
vrouwelijke sprekers, dan voor mannelijke sprekers.
<p>
Stel dat in het complete corpus meer zinnen van vrouwelijke sprekers zijn opgenomen, dan van mannen:
<p>
<table class="right">
<tr><td><em>totaal</em><td><em>waarde</em>
<tr><td>9&#8239;000<td class="left"> female
<tr><td>6&#8239;000<td class="left"> male
</table>
<p>
Alleen al op basis hiervan zou je meer matches van vrouwen dan van mannen verwachten. De tweede kolom compenseert hiervoor.
De telling van matches genormaliseerd over tienduizend zinnen geeft voor
vrouwen 36 / 9&#8239;000 &times; 10&#8239;000 = 40, en voor mannen 30 / 6&#8239;000 &times; 10&#8239;000 = 50.
<p>
Dus, <em>absoluut</em> zijn er meer treffers voor vrouwen, maar <em>relatief</em> zijn er meer treffers voor mannen.

</div>
</div>
`)
}
