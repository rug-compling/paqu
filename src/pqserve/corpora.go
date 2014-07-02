package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type Corpus struct {
	id          string
	description string
	status      string
	nline       int
	msg         string
	shared      string
	params      string
}

var (
	tr = map[string]string{
		"FINISHED": "gereed",
		"QUEUED":   "wachtrij",
		"WORKING":  "bezig",
		"FAILED":   "fout",
		"PRIVATE":  "priv&eacute;",
		"PUBLIC":   "openbaar",
		"SHARED":   "gedeeld",
	}
)

// TAB: Corpora
func corpora(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	rows, err := q.db.Query(
		fmt.Sprintf(
			"SELECT `id`, `description`, `status`, `nline`, `nword`, `msg`, `shared`, `params`FROM `%s_info` WHERE `owner` = \"%s\" ORDER BY `description`",
			Cfg.Prefix,
			q.user))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	corpora := make([]Corpus, 0)
	gebruikt := 0
	var id, desc, status, msg, shared, params string
	var zinnen, woorden int
	for rows.Next() {
		err := rows.Scan(&id, &desc, &status, &zinnen, &woorden, &msg, &shared, &params)
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		gebruikt += woorden
		corpora = append(corpora, Corpus{
			id:          id,
			description: desc,
			status:      tr[status],
			nline:       zinnen,
			msg:         msg,
			shared:      tr[shared],
			params:      params,
		})
	}

	writeHead(q, "Corpora", 2)
	fmt.Fprintln(q.w, "<h1>Mijn corpora</h1>")

	if len(corpora) == 0 {
		fmt.Fprintln(q.w, "Je hebt nog geen corpora")
	} else {
		fmt.Fprint(q.w, `<script type="text/javascript" src="jquery.js"></script>
<script type="text/javascript"><!--
var visible = -1;
function menu(idx) {
    var e = $(".submenu").filter(".a" + idx);
    if (visible == idx) {
        e.hide();
        visible = -1;
    } else {
        if (visible > -1) {
            var e2 = $(".submenu").filter(".a" + visible);
            e2.hide();
            e2.css("zIndex", 1);
        }
        e.show();
        e.css("zIndex", 9999);
        visible = idx;
    }
    return false;
}

$(document).mouseup(
    function(e) {
        if (! $(e.target).hasClass("opties")) {
            $(".submenu").hide();
            visible = -1;
        }
    });

corpora = [`)
		p := ""
		for _, corpus := range corpora {
			fmt.Fprintf(q.w, "%s\n{\"id\": \"%s\", \"title\": %q, \"lines\": %d}", p, corpus.id, corpus.description, corpus.nline)
			p = ","
		}
		fmt.Fprint(q.w, `];
function rm(idx) {
    if (window.confirm("Verwijderen: " + corpora[idx].title)) {
        var del = true;
        if (corpora[idx].lines >= 500) {
            del = window.confirm('Het corpus "' + corpora[idx].title + '" heeft ' + corpora[idx].lines +
                ' regels\n\nWeet je zeker dat je dit wilt verwijderen?')
        }
        if (del) {
            window.location.assign("delete?id=" + corpora[idx].id);
        }
        return false;
    }
    return false;
}

function trim(stringToTrim) {
  return stringToTrim.replace(/^\s+|\s+$/g,"");
}

function formtest() {
  var f = document.newcorpus;
  if (trim(f.title.value) == '') {
    alert('Titel ontbreekt');
    return false;
  }
  if (trim(f.data.value) == '') {
    alert('Geen document gekozen');
    return false;
  }
  return true;
}
//--></script>
<table class="corpora">
<tr><th><th><th>Status
<th>Titel
<th>Regels
<th>Toegang
<th>Opmerkingen
`)

		for i, corpus := range corpora {
			o := "odd"
			if i%2 == 0 {
				o = "even"
			}
			e := ""
			if i == 0 {
				e = " first"
			}
			if i == len(corpora)-1 {
				e = e + " last"
			}
			fmt.Fprintf(q.w, "<tr class=\"%s%s\">\n", o, e)

			fmt.Fprint(q.w, "<td class=\"more\">")
			if corpus.status == "gereed" || corpus.status == "fout" {
				id := urlencode(corpus.id)
				fmt.Fprintf(q.w, `<a href="javascript:void(0)" onclick="javascript:menu(%d)" class="opties a%d" title="opties">&#10020;</a>
<div class="submenu a%d" style="display: none;">
<ul class="optielijst">
`, i, i, i)
				if corpus.status == "gereed" {
					fmt.Fprintf(q.w, `
<li><a href=".?db=%s">doorzoeken</a>
<li><a href="rename?id=%s">hernoemen</a>
<li><a href="share?id=%s">delen</a>
<li><a href="download?dl=summary&id=%s">overzicht</a>
<li><a href="download?dl=zinnen&id=%s">download zinnen</a>
<li><a href="download?dl=xml&id=%s">download xml</a>
`, id, id, id, id, id, id)
					if has_dbxml && Cfg.Dact {
						fmt.Fprintf(q.w, `<li><a href="download?dl=dact&id=%s">download dact</a>`, id)
					}
				}
				fmt.Fprintf(q.w, `
<li><a href="download?dl=stdout&id=%s">download stdout</a>
<li><a href="download?dl=stderr&id=%s">download stderr</a>
</ul>
</div>
`, id, id)
			} else {
				fmt.Fprintln(q.w)
			}

			fmt.Fprintf(q.w,
				"<td class=\"delete\"><a href=\"javascript:void(0)\" onclick=\"rm(%d)\" class=\"delete\" title=\" verwijderen \">&#10008;</a>\n",
				i)

			var cl string
			switch corpus.status {
			case "gereed":
				cl = "ok"
			case "fout":
				cl = "error"
			default:
				cl = "odd"
			}
			st := corpus.status
			if st == "wachtrij" {
				processLock.RLock()
				n := taskWorkNr
				m := n + 1
				if c, ok := processes[corpus.id]; ok {
					m = c.nr
				}
				processLock.RUnlock()
				st = fmt.Sprintf("%s&nbsp;#%d", st, m-n-1)
			} else if st == "bezig" {
				if corpus.params != "dact" && corpus.params != "xmlzip" {
					p := 0
					files, err := filenames2(path.Join(paqudir, "data", corpus.id, "xml"))
					if err == nil {
						p = 1 + int(float64(len(files))/float64(corpus.nline)*98+.5)
					}
					st = fmt.Sprintf("%s&nbsp;%d%%", st, p)
				}
			}
			fmt.Fprintf(q.w, "<td class=\"%s first\">%s\n", cl, st)
			fmt.Fprintf(q.w, "<td class=\"even\">%s\n", html.EscapeString(corpus.description))

			if corpus.status == "gereed" {
				fmt.Fprintf(q.w, "<td class=\"odd right\">%s\n", iformat(corpus.nline))
				fmt.Fprintf(q.w, "<td class=\"even\">%s\n", corpus.shared)
			} else {
				n := ""
				if corpus.nline > 0 {
					n = iformat(corpus.nline)
				}
				fmt.Fprintf(q.w, "<td class=\"odd right\">%s\n<td class=\"even\">\n", n)
			}
			fmt.Fprintf(q.w, "<td class=\"odd\">%s\n", html.EscapeString(corpus.msg))
		}

		fmt.Fprintln(q.w, "</table>")
	}

	fmt.Fprintln(q.w, "<h2>Nieuw corpus maken</h2>")
	if q.quotum > 0 {
		fmt.Fprintf(q.w, "Je hebt nog ruimte voor %d woorden (tokens)\n<p>\n", q.quotum-gebruikt)
	}
	fmt.Fprintf(q.w, `
    <form name="newcorpus" action="submitcorpus" method="post" enctype="multipart/form-data"
      accept-charset="utf-8" onsubmit="javascript:return formtest()">
        De tekst die je uploadt moet platte tekst zijn, zonder opmaak (geen Word of zo), gecodeerd in utf-8.
        <p>
    Titel:<br>
	<input type="text" name="title" size="%d" maxlength="%d">
    <p>
	Upload document:<br>
	<input type="file" name="data">
        <p>
        Structuur van document:<br>
	<select name="how">
	  <option value="run">Doorlopende tekst</option>
	  <option value="line">Een zin per regel</option>
	  <option value="xmlzip">Alpino XML-bestanden in zipfile</option>
`, MAXTITLELEN+MAXTITLELEN/4, MAXTITLELEN)
	if has_dbxml {
		fmt.Fprintln(q.w, "<option value=\"dact\">Dact-bestand</option>")
	}
	fmt.Fprint(q.w, `
	</select>
      <p>
	<input type="submit">
    </form>
</body>
</html>
`)

}

func submitCorpus(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	how := firstf(q.form, "how")
	title := maxtitlelen(firstf(q.form, "title"))

	if title == "" {
		http.Error(q.w, "Titel ontbreekt", http.StatusPreconditionFailed)
		return
	}

	dirname := reNoAz.ReplaceAllString(strings.ToLower(title), "")
	if len(dirname) > 20 {
		dirname = dirname[:20]
	} else if dirname == "" {
		dirname = "a"
	}

	dirnameLock.Lock()
	defer dirnameLock.Unlock()
	for i := 0; true; i++ {
		d := dirname + abc(i)
		rows, err := q.db.Query(fmt.Sprintf("SELECT 1 FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, d))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			rows.Close()
			continue
		}
		dirname = d
		break
	}
	fulldirname := path.Join(paqudir, "data", dirname)
	err := os.Mkdir(fulldirname, 0700)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	if len(q.form.File["data"]) > 0 {
		fname := "data"
		if how == "dact" {
			fname = "data.dact"
		}
		fpout, err := os.Create(path.Join(fulldirname, fname))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		defer fpout.Close()
		fpin, err := q.form.File["data"][0].Open()
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		defer fpin.Close()
		_, err = io.Copy(fpout, fpin)
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
	}

	_, err = q.db.Exec(fmt.Sprintf("INSERT %s_info (id, description, owner, status, params) VALUES (%q, %q, %q, \"QUEUED\", %q);",
		Cfg.Prefix,
		dirname, title, q.user, how))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	logf("QUEUED: " + dirname)
	p := &Process{
		id:     dirname,
		chKill: make(chan bool),
		queued: true,
	}
	processLock.Lock()
	taskWaitNr++
	p.nr = taskWaitNr
	processes[dirname] = p
	processLock.Unlock()
	go func() {
		chWork <- p
	}()

	writeHtml(
		q,
		"Document word verwerkt",
		`
Je document wordt verwerkt. Als het klaar is zie je op de hoofdpagina een nieuw corpus bij de databases.
<p>
Let op: Dit kan even duren. Minuten, uren, of dagen, afhankelijk van de grootte van je document.
<p>
<b>Je krijgt een e-mail als het corpus klaar is.</b>
`)
}

func abc(i int) string {
	b := make([]byte, 0)
	for {
		j := i % 26
		i = i / 26
		b = append(b, byte(j+97))
		if i == 0 {
			break
		}
	}
	return string(b)
}

func vertaal(s string) string {
	if s1, ok := tr[s]; ok {
		return s1
	}
	return s
}
