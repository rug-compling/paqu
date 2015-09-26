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
	s string
	i int
	n int
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

		for run := 0; run < 2; run++ {
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
					"SELECT COUNT(`text`), `text`, 0 "+
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
					"SELECT COUNT(`a`.`text`), `a`.`text`,`a`.`n` FROM ( %s ) `a` GROUP BY `a`.`idx` ORDER BY %s `a`.`idx`",
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
				var cnt, n int
				var text string
				errval = rows.Scan(&cnt, &text, &n)
				if logerr(errval) {
					rows.Close()
					return
				}
				if len(lines) < METAMAX || download {
					text = unHigh(text)
					lines = append(lines, Statline{text, cnt, n})
				}
			}
			errval = rows.Err()
			if logerr(errval) {
				return
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

In de tabel <em>per zin</em> staan tussen de kolommen met aantallen en met de metadata-waardes nog een kolom.
<p>
<b>Kolom 2: <em>n</em> keer de fractie</b>
<p>
De fractie is het aantal gevonden zinnen met de metadata-waarde, gedeeld door het totaal aantal zinnen dat deze waarde heeft.
<p>
De Tweede kolom geeft deze fractie vermenigvuldigd met <em>n</em>. De waarde van <em>n</em> is afankelijk van de grootte van het corpus,
en staat boven de tabellen vermeld.
<p>
Een waarde gelijk aan <em>n</em> wil zeggen dat alle zinnen waarin de waarde voorkomt voldoen aan je zoekopdracht.
</div>
</div>
`)
}
