package main

import (
	"fmt"
	"html"
	"strings"
)

func corpusinfo(q *Context) {

	writeHead(q, "Corpusinfo", 0)
	fmt.Fprint(q.w, `
<h1>Corpusinformatie</h1>
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
var infop = {
`)

	p := ""
	ids := make([]string, 0)
	for _, corpus := range q.opt_db {
		a := strings.SplitN(corpus, " ", 2)
		opt := a[0]
		if opt[0] == '-' {
			opt = opt[1:]
		}
		db := opt[1:]
		if t := q.infops[db]; t != "" {
			fmt.Fprintf(q.w, "%s%s: %q", p, db, t)
			p = ",\n"
			ids = append(ids, fmt.Sprintf("%s: %q", db, a[1]))
		}
	}

	fmt.Fprintf(q.w, `
};
var titles = {
%s
};`, strings.Join(ids, ",\n"))

	fmt.Fprint(q.w, `
var infovisible = false;
function info(id) {
  $("#title").text(titles[id]);
  $("#content").html(infop[id]);
  var e = $("#infocorpus");
  e.show();
  e.css("zIndex", 9999);
  infovisible = true;
  return false;
}
$(document).mouseup(
  function(e) {
    if (infovisible) {
      var e = $("#infocorpus");
      e.hide();
      e.css("zIndex", 1);
      infovisible = false;
    }
  });
//--></script>
<div class="submenu a9999" id="infocorpus">
<div id="title"></div>
<div id="content"></div>
</div>
<table class="corpusinfo">
<tr><th><th class="left">Titel<th class="right">Regels<th class="right">Datum<th class="left">Metadata<th class="left">UD<th class="left">Eigenaar</tr>
`)

	t := ""
	var c byte
	for _, corpus := range q.opt_db {
		a := strings.SplitN(corpus, " ", 2)
		opt := a[0]
		if opt[0] == '-' {
			opt = opt[1:]
		}
		if opt[0] != c {
			c = opt[0]
			switch c {
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
			if c != 'Z' {
				fmt.Fprintf(q.w, "<tr class=\"sub\"><th><th colspan=\"6\" class=\"left\">%s</tr>\n", t)
			}
		}
		db := opt[1:]
		if q.infops[db] == "" {
			fmt.Fprint(q.w, "<tr><td>")
		} else {
			fmt.Fprintf(q.w, "<tr><td><a href=\"javascript:void(0)\" onclick=\"javascript:info('%s')\">[?]</a>", db)
		}
		metadata := "nee"
		if q.hasmeta[db] {
			metadata = "ja"
		}
		ud := "nee"
		if q.hasud[db] {
			ud = "ja"
		}
		fmt.Fprintf(q.w, "<td>%s<td class=\"right\">%s<td class=\"right\">%s<td>%s<td>%s<td>%s</tr>\n",
			html.EscapeString(a[1]), iformat(q.lines[db]), datum(q.dates[db]), metadata, ud, html.EscapeString(displayEmail(q.owners[db])))
	}

	fmt.Fprint(q.w, `
</table>
</body>
</html>
`)
}
