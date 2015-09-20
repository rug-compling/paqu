package main

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
)

type StructIS struct {
	i int
	s string
}

// TAB: metadata
func metadata(q *Context) {

	prefix := getprefix(q)
	if !q.prefixes[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	// HTML-uitvoer van begin van de pagina
	writeHead(q, "", 3)

	fmt.Fprint(q.w, `
<form action="metadata" method="get" accept-charset="utf-8">
corpus: <select name="db">
`)
	html_opts(q, q.opt_dbmeta, getprefix(q), "corpus")
	fmt.Fprintln(q.w, "</select>")
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpuslijst?meta=1\">meer/minder</a>")
	}
	fmt.Fprint(q.w, `
       <p>
       <input type="submit" value="Kies corpus">
	   </form>
       <p>
	   `)

	if !q.hasmeta[prefix] {
		html_footer(q)
		return
	}
	fmt.Fprintln(q.w, "<hr><p><div id=\"stats\">")

	metas := getMeta(q, prefix)
	for _, meta := range metas {
		fmt.Fprintf(q.w, "<b>%s</b><p>\n", html.EscapeString(meta.name))
		o := `idx`
		align := "right"
		limit := ""
		if meta.mtype == "TEXT" {
			o = "`n` DESC, `idx`"
			align = "left"
			limit = fmt.Sprintf(" LIMIT 0, %d", METAMAX)
		}
		rows, err := q.db.Query(fmt.Sprintf(
			"SELECT `text`, `n` FROM `%s_c_%s_mval` WHERE `id`=%d ORDER BY %s%s",
			Cfg.Prefix, prefix,
			meta.id,
			o,
			limit))
		if doErr(q, err) {
			return
		}
		fmt.Fprintln(q.w, "<table>")
		count := 0
		for rows.Next() {
			count++
			var s string
			var n int
			if doErr(q, rows.Scan(&s, &n)) {
				return
			}
			cnil := ""
			if s == "" {
				s = "(leeg)"
				cnil = " nil"
			} else {
				s = html.EscapeString(s)
			}
			fmt.Fprintf(q.w, "<tr><td class=\"right\">%s<td class=\"%s%s\">%s\n", iformat(n), align, cnil, s)
		}
		if meta.mtype == "TEXT" && count == METAMAX {
			fmt.Fprintln(q.w, "<tr><td><td>...")
		}
		fmt.Fprintln(q.w, "</table><p>")
		if doErr(q, rows.Err()) {
			return
		}
	}
	fmt.Fprintf(q.w, "<hr><a href=\"metadl?db=%s\">download</a></div>\n", prefix)

	meta2form(q, prefix, metas)

	fmt.Fprintln(q.w, `<p><b>In ontwikkeling</b>`)

	html_footer(q)
}

func metadl(q *Context) {

	prefix := first(q.r, "db")
	if !q.hasmeta[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	q.w.Header().Set("Content-Disposition", "attachment; filename=metadata.txt")
	cache(q)

	metas := getMeta(q, prefix)
	for _, meta := range metas {
		fmt.Fprintf(q.w, "# %s\n", meta.name)
		o := `idx`
		if meta.mtype == "TEXT" {
			o = "`n` DESC, `idx`"
		}
		rows, err := q.db.Query(fmt.Sprintf(
			"SELECT `text`, `n` FROM `%s_c_%s_mval` WHERE `id`=%d ORDER BY %s",
			Cfg.Prefix, prefix,
			meta.id,
			o))
		if doErr(q, err) {
			return
		}
		for rows.Next() {
			var s string
			var n int
			if doErr(q, rows.Scan(&s, &n)) {
				return
			}
			fmt.Fprintf(q.w, "%d\t%s\n", n, s)
		}
		if doErr(q, rows.Err()) {
			return
		}
	}
}

func meta2form(q *Context, prefix string, metas []MetaType) {
	if len(metas) < 2 {
		return
	}
	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
  var lastcall = null;
  $.fn.meta2 = function() {
    if (lastcall) {
      try {
        lastcall.abort();
      }
      catch(err) {}
    }
    $("#meta2res").html('<img src="busy.gif">');
    lastcall = $.ajax("meta2?" + $(document.meta2form).serialize())
      .done(function(data) {
        $("#meta2res").html(data);
      }).fail(function(e) {
        $("#meta2res").html(e.responseText);
      })
      .always(function() {
        lastcall = null;
      });
  }
//--></script>
<p>
<div id="statsrel">
<form action="javascript:$.fn.meta2()" name="meta2form" id="meta2form">
Selecteer twee elementen om ze te koppelen:
<p>
`)
	for i, meta := range metas {
		c := ""
		if i < 2 {
			c = `checked="checked"`
		}
		fmt.Fprintf(q.w, "<input type=\"checkbox\" %s name=\"cmeta\" value=\"%d\">%s<br>\n", c, meta.id, html.EscapeString(meta.name))
	}
	fmt.Fprintf(q.w, `
<p>
<input type="hidden" name="db" value="%s">
<input type="submit" value="telling van combinatie" id="meta2submit">
</form>
<p>
<div id="meta2res"></div>
</div>
<script type="text/javascript"><!--
  $('#meta2form input').on('change', function (e) {
    var f = document.forms["meta2form"];
    var n = 0;
    for (i = 0; i < f.cmeta.length; i++) {
       if (f.cmeta[i].checked) { n++; }
    }
    if (n ==  2) {
      $('#meta2submit').prop('disabled', false);
    } else {
      $('#meta2submit').prop('disabled', true);
    }
  });
//--></script>
`, prefix)

}

func meta2(q *Context) {

	download := first(q.r, "d") != ""

	if download {
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		q.w.Header().Set("Content-Disposition", "attachment; filename=metadata.txt")
	} else {
		q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	cache(q)

	prefix := first(q.r, "db")
	if !q.hasmeta[prefix] {
		http.Error(q.w, "Invalid corpus: "+prefix, http.StatusPreconditionFailed)
		return
	}

	if len(q.r.Form["cmeta"]) != 2 {
		http.Error(q.w, "Missing options", http.StatusPreconditionFailed)
		return
	}

	var mi [2]int
	mi[0], _ = strconv.Atoi(q.r.Form["cmeta"][0])
	mi[1], _ = strconv.Atoi(q.r.Form["cmeta"][1])
	var ok [2]bool
	var met [2]MetaType
	metas := getMeta(q, prefix)
	for _, meta := range metas {
		if meta.id == mi[0] {
			met[0] = meta
			ok[0] = true
		} else if meta.id == mi[1] {
			met[1] = meta
			ok[1] = true
		}
	}
	if !(ok[0] && ok[1]) {
		http.Error(q.w, "Invalid options", http.StatusPreconditionFailed)
		return
	}

	var keys [2][]StructIS
	for i := 0; i < 2; i++ {
		keys[i] = make([]StructIS, 0)
		rows, err := q.db.Query(fmt.Sprintf(
			"SELECT `idx`,`text` FROM `%s_c_%s_mval` WHERE `id`=%d ORDER BY `idx`",
			Cfg.Prefix, prefix, met[i].id))
		if logerrfrag(q, err) {
			return
		}
		for rows.Next() {
			var ii int
			var s string
			if logerrfrag(q, rows.Scan(&ii, &s)) {
				return
			}
			keys[i] = append(keys[i], StructIS{ii, s})
		}
		if logerrfrag(q, rows.Err()) {
			return
		}
	}
	if len(keys[0]) > len(keys[1]) {
		keys[0], keys[1] = keys[1], keys[0]
		met[0], met[1] = met[1], met[0]
	}
	count := make(map[[2]int]int)

	rows, err := q.db.Query(fmt.Sprintf(
		"SELECT COUNT(*), `a`.`idx`, `b`.`idx` FROM `%s_c_%s_meta` `a` JOIN `%s_c_%s_meta` `b` "+
			"USING (`arch`,`file`) WHERE `a`.`id` = %d AND `b`.`id` = %d GROUP BY `a`.`idx`, `b`.`idx`",
		Cfg.Prefix, prefix,
		Cfg.Prefix, prefix,
		met[0].id,
		met[1].id))
	if logerrfrag(q, err) {
		return
	}
	for rows.Next() {
		var c, i0, i1 int
		if logerrfrag(q, rows.Scan(&c, &i0, &i1)) {
			return
		}
		count[[2]int{i0, i1}] = c
	}
	if logerrfrag(q, rows.Err()) {
		return
	}

	yalign := "right"
	if met[1].mtype == "TEXT" {
		yalign = "left"
	}

	if download {
		fmt.Fprintf(q.w, "# %s %s\n", met[0].name, met[1].name)
	} else {
		fmt.Fprintf(q.w, `
<table class="breed">
<tr><td><td><td colspan="%d"><b>%s</b>
<tr><td><td>
`, len(keys[0]), html.EscapeString(met[0].name))
	}
	for i, k := range keys[0] {
		if download {
			if i > 0 {
				fmt.Fprint(q.w, "\t")
			}
			fmt.Fprint(q.w, k.s)
		} else {
			if k.s == "" {
				fmt.Fprintln(q.w, "<td class=\"nil right\">(leeg)")
			} else {
				fmt.Fprintf(q.w, "<td class=\"right\">%s\n", html.EscapeString(k.s))
			}
		}
	}
	if download {
		fmt.Fprintln(q.w)
	}
	for i, k := range keys[1] {
		if download {
			fmt.Fprint(q.w, k.s)
		} else {
			if i%2 == 1 {
				fmt.Fprint(q.w, "<tr class=\"odd\">")
			} else {
				fmt.Fprint(q.w, "<tr>")
			}
			if i == 0 {
				fmt.Fprintf(q.w, "<td rowspan=\"%d\"><b>%s</b>\n", len(keys[1]), html.EscapeString(met[1].name))
			}
			if k.s == "" {
				fmt.Fprintf(q.w, "<td class=\"nil %s\">(leeg)\n", yalign)
			} else {
				fmt.Fprintf(q.w, "<td class=\"%s\">%s\n", yalign, html.EscapeString(k.s))
			}
		}
		for j := 0; j < len(keys[0]); j++ {
			if download {
				fmt.Fprint(q.w, "\t", count[[2]int{keys[0][j].i, k.i}])
			} else {
				fmt.Fprintf(q.w, "<td class=\"right\">%s\n", iformat(count[[2]int{keys[0][j].i, k.i}]))
			}
		}
		if download {
			fmt.Fprintln(q.w)
		}
	}
	if !download {
		fmt.Fprintf(q.w, `
</table>
<p>
<hr>
<a href="meta2?db=%s&cmeta=%d&cmeta=%d&d=1">download</a>
`, html.EscapeString(prefix), mi[0], mi[1])
	}
}
