package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
)

type CorpusType struct {
	id    string
	title string
	lines int
	owner string
}

func corpuslijst(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	meta := false
	if first(q.r, "meta") != "" {
		meta = true
	}

	rows, err := q.db.Query(fmt.Sprintf(
		"SELECT `i`.`id`, `i`.`description`, `i`.`nline`, `i`.`owner` "+
			"FROM `%s_info` `i`, `%s_corpora` `c` "+
			"WHERE `i`.`id` = `c`.`prefix` "+
			"AND `c`.`enabled` = 1 "+
			"AND (`c`.`user` = %q OR `c`.`user` = \"all\") "+
			"AND `i`.`owner` != %q "+
			"AND `i`.`owner` != \"none\" "+
			"ORDER BY 2", Cfg.Prefix, Cfg.Prefix, q.user, q.user))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	corpora := make([]CorpusType, 0)
	for rows.Next() {
		var id, title, owner string
		var lines int
		err := rows.Scan(&id, &title, &lines, &owner)
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if !meta || q.hasmeta[id] {
			corpora = append(corpora, CorpusType{
				id:    id,
				title: title,
				lines: lines,
				owner: owner,
			})
		}
	}

	writeHead(q, "Corpuslijst", 0)

	if len(corpora) == 0 {
		if meta {
			fmt.Fprint(q.w, "Niemand heeft corpora met metadata met je gedeeld")
		} else {
			fmt.Fprint(q.w, "Niemand heeft corpora met je gedeeld")
		}
		return
	}

	p := ""
	fmt.Fprint(q.w, `
<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
var field = "title";
var reverse = false;
var data = [`)
	for _, c := range corpora {
		fmt.Fprintf(q.w, `%s
{id:%q,title:%q,lower:%q,lines:%d,liness:%q,owner:%q,show:%v}`,
			p,
			c.id,
			html.EscapeString(c.title),
			strings.ToLower(c.title),
			c.lines,
			iformat(c.lines),
			html.EscapeString(displayEmail(c.owner)),
			!q.ignore[c.id])
		p = ","
	}
	fmt.Fprint(q.w, `
];

function opslaan() {
    $("#submit").attr("disabled", "disabled");
    var lines = "";
    for (i in data) {
        if (!data[i].show) {
            lines += data[i].id + "\n";
        }
    }
    var jqxhr = $.post(
    "corsave",
    lines,
    function(result) {
        $("#msg").html(result);
    }).fail(function(jqXHR, textStatus, errorThrown) {
        $("#msg").html(errorThrown);
    });

}

function toggle(idx) {
    $("#msg").html("");
    $("#submit").removeAttr("disabled");
    data[idx].show = !data[idx].show;
}

function redraw () {
    var lines = "";
    $.each(data, function(index, value) {
        var c = "";
        var eo = "odd";
        if (index % 2 == 0) {
            eo = "even";
        }
        if (index == 0) {
            eo += " first";
        }
        if (index == data.length - 1) {
            eo += " last";
        }
        if (value.show) {
            c = " checked=\"checked\"";
        }
        lines += "<tr class=\"" + eo + "\"><td><input type=\"checkbox\"" + c + " onchange=\"toggle(" + index + ")\"><td class=\"odd first\">" + value.title + "<td class=\" even right\">" + value.liness + "<td class=\"odd\">" + value.owner + "\n";
    });
    $("#items").html(lines);
    $("#title").html("");
    $("#lines").html("");
    $("#owner").html("");
    var pijl = "&darr;";
    if (reverse) {
        pijl = "&uarr;";
    }
    if (field == "title") {
        $("#title").html(pijl);
    }
    if (field == "lines") {
        $("#lines").html(pijl);
    }
    if (field == "owner") {
        $("#owner").html(pijl);
    }
}

function doTitle() {
    if (field == "title") {
        reverse = !reverse;
    } else {
        field = "title";
    }
    var rv = reverse ? -1 : 1;
    data.sort(function(a, b) {
      if (a.lower > b.lower) { return rv; }
      if (a.lower < b.lower) { return -rv; }
      return 0;
    });
    redraw();
}

function doLines() {
    if (field == "lines") {
        reverse = !reverse;
    } else {
        field = "lines";
    }
    var rv = reverse ? -1 : 1;
    data.sort(function(a, b) {
      if (a.lines == b.lines) {
        if (a.lower > b.lower) { return rv; }
        if (a.lower < b.lower) { return -rv; }
        return 0;
      }
      return rv * (a.lines - b.lines);
    });
    redraw();
}

function doOwner() {
    if (field == "owner") {
        reverse = !reverse;
    } else {
        field = "owner";
    }
    var rv = reverse ? -1 : 1;
    data.sort(function(a, b) {
      if (a.owner > b.owner) { return rv; }
      if (a.owner < b.owner) { return -rv; }
      if (a.lower > b.lower) { return rv; }
      if (a.lower < b.lower) { return -rv; }
      return 0;
    });
    redraw();
}

$(document).ready(redraw);

//--></script>
Hieronder zie je een lijst van corpora die door anderen met je zijn gedeeld.
<p>`)
	if meta {
		fmt.Fprint(q.w, `
Je ziet nu alleen corpora die <b>metadata</b> bevatten
<p>
`)
	}
	fmt.Fprint(q.w, `	
Selecteer welke corpora je in het menu wilt zien.
<p>
<form action="javascript:void(0)" onsubmit="javascript:opslaan()">
<table id="corpustabel" class="corpora">
<thead>
<tr>
  <th>
  <th><a href="javascript:void(0)" onclick="javascript:doTitle()">Titel<span id="title"></span></a>
  <th><a href="javascript:void(0)" onclick="javascript:doLines()">Regels<span id="lines"></span></a>
  <th><a href="javascript:void(0)" onclick="javascript:doOwner()">Eigenaar<span id="owner"></span></a>
</tr>
</thead>
<tbody id="items">
</tbody>
</table>
<input id="submit" type="submit" value="Opslaan" disabled="disabled">
<span id="msg"></span>
</form>
</body>
</html>
`)

}

func corsave(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	data, err := ioutil.ReadAll(q.r.Body)
	if err != nil {
		http.Error(q.w, html.EscapeString(err.Error()), http.StatusInternalServerError)
		logerr(err)
		return
	}

	_, err = q.db.Exec(fmt.Sprintf("DELETE FROM `%s_ignore` WHERE `user` = %q", Cfg.Prefix, q.user))
	if err != nil {
		http.Error(q.w, html.EscapeString(err.Error()), http.StatusInternalServerError)
		logerr(err)
		return
	}

	for _, id := range strings.Fields(string(data)) {
		_, err = q.db.Exec(fmt.Sprintf("INSERT INTO `%s_ignore` (`user`, `prefix`) VALUES (%q, %q)", Cfg.Prefix, q.user, id))
		if err != nil {
			http.Error(q.w, html.EscapeString(err.Error()), http.StatusInternalServerError)
			logerr(err)
			return
		}
	}

	fmt.Fprint(q.w, "Gegevens opgeslagen")
}
