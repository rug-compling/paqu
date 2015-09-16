package main

import (
	"bytes"
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
	r float64
	n int
}

func statsmeta(q *Context) {

	var buf bytes.Buffer

	var chClose <-chan bool
	if f, ok := q.w.(http.CloseNotifier); ok {
		chClose = f.CloseNotify()
	} else {
		chClose = make(<-chan bool)
	}

	now := time.Now()

	download := false
	if first(q.r, "d") != "" {
		download = true
	}

	option := make(map[string]string)
	for _, t := range []string{"word", "postag", "rel", "hpostag", "hword"} {
		option[t] = first(q.r, t)
	}
	if option["word"] == "" && option["postag"] == "" && option["rel"] == "" && option["hpostag"] == "" && option["hword"] == "" {
		http.Error(q.w, "Query ontbreekt", http.StatusPreconditionFailed)
		return
	}

	prefix := first(q.r, "db")
	if prefix == "" {
		http.Error(q.w, "Geen corpus opgegeven", http.StatusPreconditionFailed)
		return
	}
	if !q.prefixes[prefix] {
		http.Error(q.w, "Ongeldig corpus", http.StatusPreconditionFailed)
		return
	}

	metas := getMeta(q, prefix)

	query, err := makeQuery(q, prefix, "", chClose)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	// BEGIN UITVOER

	pow10 := math.Pow10(int(math.Log10(float64(q.lines[prefix])) + .5))
	if pow10 < 10 {
		pow10 = 10
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

	// DEBUG: HTML-uitvoer van de query
	if !download {
		fmt.Fprintf(&buf, `
<div style="font-family:monospace">%s</div>
<p>
<a href="javascript:void(0)" onclick="javascript:metahelp()">toelichting bij tabellen</a>
<p>
<em>n</em> = %s
<p>
`, html.EscapeString(query), iformat(int(pow10)))
		updateText(q, buf.String())
		buf.Reset()
	}

	// Tellingen van onderdelen
	for _, meta := range metas {

		// telling van metadata in matchende zinnen

		for run := 0; run < 2; run++ {
			lines := make([]Statline, 0)
			var count int
			var sum = 0.0
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
			var qu string
			var order string
			if meta.mtype == "TEXT" {
				order = "1 DESC,"
			}
			if run == 0 {
				qu = fmt.Sprintf(
					"SELECT COUNT(`text`), `text`, 0, 0 "+
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
					"SELECT DISTINCT `arch`,`file`,`idx`,`text`,`n`,`r` "+
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
					"SELECT COUNT(`a`.`text`), `a`.`text`,`a`.`n`,`a`.`r` FROM ( %s ) `a` GROUP BY `a`.`idx` ORDER BY %s `a`.`idx`",
					qu,
					order)
			}
			rows, err := timeoutQuery(q, chClose, qu)
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
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
				var r float64
				err := rows.Scan(&cnt, &text, &n, &r)
				if err != nil {
					updateError(q, err, !download)
					completedmeta(q, download)
					logerr(err)
					return
				}
				if len(lines) < METAMAX || download {
					text = unHigh(text)
					lines = append(lines, Statline{text, cnt, float64(cnt) / r, n})
				}
				sum += float64(cnt) / r
			}
			err = rows.Err()
			if err != nil {
				updateError(q, err, !download)
				completedmeta(q, download)
				logerr(err)
				return
			}
			for _, line := range lines {
				var p float64
				if run == 1 {
					p = float64(line.r) / sum * 100.0
				}
				if download {
					if run == 1 {
						v := int(.5 + pow10*float64(line.i)/float64(line.n))
						fmt.Fprintf(q.w, "%d\t%.2f%%\t%d\t%s\n", line.i, float32(p), v, line.s)
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
						fmt.Fprintf(&buf, "<td>%.2f%%<td>%s", float32(p), iformat(v))
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
		completedmeta(q, download)
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

In de tabel <em>per zin</em> staan tussen de kolommen met aantallen en met de metadata-waardes nog twee kolommen.
<p>
Kolom 2: genormaliseerd percentage
<p>
Kolom 3: <em>n</em> keer de fractie
<p>
<b>Kolom 2: genormaliseerd percentage</b>
<p>
Deze kolom geeft voor elke metadata-waarde het percentage van het aantal van alle metadata-waardes,
genormaliseerd voor het verwachte percentage op basis van verhoudingen van metadata-waardes in alle zinnen.
<p>
Stel dat je twee waardes hebt in de data, <em>man</em> en <em>vrouw</em>. Wanneer de genormaliseerde percentage 50%/50% zijn wil dat niet zeggen
dat je van voor beide waardes hetzelfde aantal zinnen hebt gevonden, maar dat de verhouding tussen <em>man</em> en <em>vrouw</em>
in de gevonden zinnen gelijk is aan de verhouding in alle zinnen. Wanneer de ene waarde een groter percentage geeft dan een andere,
betekent dat dat de eerste waarde verhoudingsgewijs vaker wordt gevonden. Dat slaat los van de absolute aantallen van de gevonden zinnen.
<p>
TODO: Formule?
<p>
<b>Kolom 3: <em>n</em> keer de fractie</b>
<p>
De fractie is het aantal gevonden zinnen met de metadata-waarde, gedeeld door het totaal aantal zinnen dat deze waarde heeft.
<p>
De derde kolom geeft deze fractie vermenigvuldigd met <em>n</em>. De waarde van <em>n</em> is afankelijk van de grootte van het corpus,
en staat boven de tabellen vermeld.
<p>
Een waarde gelijk aan <em>n</em> wil zeggen dat alle zinnen waarin de waarde voorkomt voldoen aan je zoekopdracht.
</div>
</div>
`)
}
