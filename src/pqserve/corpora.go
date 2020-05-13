package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Corpus struct {
	id          string
	description string
	status      string
	nline       int
	msg         string
	shared      string
	params      string
	datum       time.Time
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

	rows, err := sqlDB.Query(
		fmt.Sprintf(
			"SELECT `id`, `description`, `status`, `nline`, `nword`, `msg`, `shared`, `params`, `created` FROM `%s_info` WHERE `owner` = \"%s\" ORDER BY `description`",
			Cfg.Prefix,
			q.user))
	if hErr(q, err) {
		return
	}

	corpora := make([]Corpus, 0)
	gebruikt := 0
	var id, desc, status, msg, shared, params string
	var zinnen, woorden int
	var date time.Time
	for rows.Next() {
		err = rows.Scan(&id, &desc, &status, &zinnen, &woorden, &msg, &shared, &params, &date)
		if hErr(q, err) {
			rows.Close()
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
			datum:       date,
		})
	}

	writeHead(q, "Corpora", 4)
	fmt.Fprintln(q.w, `
<script type="text/javascript" src="jquery.js"></script>
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

<h1>Mijn corpora</h1>
`)

	foliafile := filepath.Join(foliadir(q), "error.txt")
	if first(q.r, "folia") == "clear" {
		os.Remove(foliafile)
	} else {
		data, err := ioutil.ReadFile(foliafile)
		if err == nil {
			fmt.Fprintf(q.w, `
<div class="error">
Er ging iets fout met de invoer van een corpus in FoLiA-formaat
<div class="output">
%s
</div>
<a href="corpora?folia=clear">Sluiten</a>
</div>
`, html.EscapeString(strings.Replace(string(data), foliadir(q), "...", -1)))
		}
	}

	if len(corpora) == 0 {
		fmt.Fprintln(q.w, "Je hebt nog geen corpora")
	} else {
		fmt.Fprint(q.w, `<script type="text/javascript"><!--
corpora = [`)
		p := ""
		for _, corpus := range corpora {
			fmt.Fprintf(q.w, "%s\n{\"id\": \"%s\", \"title\": %q, \"lines\": %q}", p, corpus.id, corpus.description, iformat(corpus.nline))
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
//--></script>
<table class="corpora">
<tr><th><th><th>Status
<th>Titel
<th>Regels
<th>Datum
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
				fmt.Fprintf(q.w, `<a href="javascript:void(0)" onclick="javascript:menu(%d)" class="opties a%d" title=" opties ">&#10020;</a>
<div class="submenu a%d" style="display: none;">
<ul class="optielijst">
`, i, i, i)
				if corpus.status == "gereed" {
					fmt.Fprintf(q.w, `
<li><a href=".?db=%s">doorzoeken</a>
<li><a href="rename?id=%s">hernoemen</a>
`, id, id)
					if !q.protected[id] {
						fmt.Fprintf(q.w, `
<li><a href="share?id=%s">delen</a>
`, id)
					} else {
						fmt.Fprintln(q.w, `<li><span class="disabled">delen</span>`)
					}
					fmt.Fprintf(q.w, `
<li><a href="browse?id=%s">overzicht</a>
<li><a href="download?dl=zinnen&amp;id=%s">bekijk zinnen</a>
`, id, id)
					if !q.protected[id] {
						fmt.Fprintf(q.w, `
<li><a href="download?dl=xml&amp;id=%s">download xml</a>
`, id)
					} else {
						fmt.Fprintln(q.w, `<li><span class="disabled">download xml</span>`)
					}
					if has_dbxml && Cfg.Dact {
						if !q.protected[id] {
							fmt.Fprintf(q.w, `<li><a href="download?dl=dact&amp;id=%s">download dact</a>`, id)
						} else {
							fmt.Fprintln(q.w, `<li><span class="disabled">download dact</span>`)
						}
					}
				}
				fmt.Fprintf(q.w, `
<li><a href="download?dl=stdout&amp;id=%s">bekijk stdout</a>
<li><a href="download?dl=stderr&amp;id=%s">bekijk stderr</a>
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
				if corpus.nline > 0 && (strings.HasPrefix(corpus.params, "run") ||
					strings.HasPrefix(corpus.params, "line") ||
					strings.HasPrefix(corpus.params, "folia") ||
					strings.HasPrefix(corpus.params, "tei")) {
					files := countXML(filepath.Join(paqudatadir, "data", corpus.id, "xml"))
					p := 1 + int(float64(files)/float64(corpus.nline)*98+.5)
					st = fmt.Sprintf("%s&nbsp;%d%%", st, p)
				}
			}
			fmt.Fprintf(q.w, "<td class=\"%s first\">%s\n", cl, st)
			fmt.Fprintf(q.w, "<td class=\"even\">%s\n", html.EscapeString(corpus.description))

			if corpus.status == "gereed" {
				fmt.Fprintf(q.w, "<td class=\"odd right\">%s\n", iformat(corpus.nline))
				fmt.Fprintf(q.w, "<td class=\"even right\">%s\n", strings.Replace(datum(corpus.datum), " ", "&nbsp;", -1))
				fmt.Fprintf(q.w, "<td class=\"odd\">%s\n", corpus.shared)
			} else {
				n := ""
				if corpus.nline > 0 {
					n = iformat(corpus.nline)
				}
				fmt.Fprintf(q.w, "<td class=\"odd right\">%s\n<td class=\"even\">\n<td class=\"odd\">\n", n)
			}
			fmt.Fprintf(q.w, "<td class=\"even\">%s\n", html.EscapeString(corpus.msg))
		}

		fmt.Fprintln(q.w, "</table>")
	}

	fmt.Fprint(q.w, `
<div class="submenu a9999">
<div class="corpushelp">
Keuzes voor het soort document. Zie ook <a href="info.html#corpora" target="_blank">nadere uitleg</a>.
<p>
<dl>
<dt>Automatisch bepaald
<dd>Een van onderstaande formaten wordt gedetecteerd. Dit gaat niet altijd goed. Kies in dat geval zelf het formaat.
<dt>Doorlopende tekst
<dd>Een bestand met platte tekst, met zinnen die doorlopen over regeleindes. Voorbeeld:
<pre>
  % commentaar
  label |
  Dit is de eerste zin. Dit is
  de tweede zin. Dit is zin nummer drie.
</pre>
Commentaar en label worden genegeerd.
<dt>Een zin per regel
<dd>Een bestand met platte tekst, met één zin per regel. Voorbeeld:
<pre>
  % commentaar
  Dit is de eerste zin.
  Dit is de tweede zin.
</pre>
Commentaar wordt genegeerd.
<dt>Een zin per regel, met labels
<dd>Als boven, met labels toegevoegd. Label wordt van zin gescheiden door verticale streep, zonder spaties aan weerszijde. Voorbeeld:
<pre>
  % commentaar
  zin 1|Dit is de eerste zin.
  zin 2|Dit is de tweede zin.
</pre>
Commentaar wordt genegeerd.
<dt>Een zin per regel, getokeniseerd
<dd>Tekst die al getokeniseerd is, één zin per regel. Voorbeeld:
<pre>
  % commentaar
  Dit is de eerste zin .
  Dit is de tweede zin .
</pre>
Commentaar wordt genegeerd.
<dt>Een zin per regel, met labels, getokeniseerd
<dd>Tekst die al getokeniseerd is, en waar labels aan de zinnen is toegevoegd. Voorbeeld:
<pre>
  % commentaar
  zin 1|Dit is de eerste zin .
  zin 2|Dit is de tweede zin .
</pre>
Commentaar wordt genegeerd.
</dl>
<b>Andere formaten</b> die automatisch worden herkend:
<ul>
<li>Een zip- of tarbestand met daarin bestanden met platte tekst. In dit geval geldt ook de keus die je hebt gemaakt over het soort bestand.
<li>XML-bestanden met één door Alpino geparste zin per bestand, samengevoegd in een zip- of tarbestand.
`)
	if has_dbxml {
		fmt.Fprint(q.w, `
<li>Een bestand in het DbXML-formaat, waarin Alpino XML-bestanden zijn opgeslagen. Dit formaat wordt onder andere gebruikt door het programma
<a href="http://rug-compling.github.io/dact/" target="_blank">dact</a>.
`)
	}
	fmt.Fprint(q.w, `
<li>Een bestand in FoLiA-formaat: <a href="http://proycon.github.io/folia/" target="_blank">Format for Linguistic Annotation</a>.
Het bestand moet gecodeerd zijn in UTF-8.
De tekst moet getokeniseerd zijn. Als je ook metadata hebt, klik dan <a href="folia">hier</a>.
<li>Een zip- of tarbestand met daarin bestanden in FoLiA-formaat. Als je ook metadata hebt, klik dan <a href="folia">hier</a>.
<li>Een bestand in TEI-formaat: <a href="http://www.tei-c.org/" target="_blank">Text Encoding Initiative</a>.
Het bestand moet gecodeerd zijn in UTF-8.
De tekst moet getokeniseerd zijn.
<li>Een zip- of tarbestand met daarin bestanden in TEI-formaat.
</ul>
Als je een zip- of tarbestand gebruikt, dan moeten de bestanden die daarin zitten allemaal hetzelfde formaat hebben.
<p>
Alle bestanden mogen gecomprimeerd zijn met gzip, behalve de bestanden in een zip- of tarbestand.
Een tarbestand zelf mag wel gecomprimeerd zijn met gzip.
</div>
</div>
<h2>Nieuw corpus maken</h2>
`)
	if q.quotum > 0 {
		fmt.Fprintf(q.w, "Je hebt nog ruimte voor %s woorden (tokens)\n", iformat(q.quotum-gebruikt))
	}
	fmt.Fprintf(q.w, `
    <div class="info">
    Heb je een corpus in <b>FoLiA</b>-formaat, met interne of externe <b>metadata</b>, klik dan <b><a href="folia">hier</a></b>.
    </div>
    <form name="newcorpus" action="submitcorpus" method="post" enctype="multipart/form-data"
      accept-charset="utf-8" onsubmit="javascript:return formtest()">
        De tekst die je uploadt moet platte tekst zijn, zonder opmaak (geen Word of zo), gecodeerd in utf-8.<br>
        Daarnaast worden een aantal andere formaten herkend, zie <a href="javascript:void(0)" onclick="javascript:menu(9999)">uitleg</a>.
        <p>
    Titel:<br>
	<input type="text" name="title" size="%d" maxlength="%d">
    <p>
	Upload document:<br>
	<input type="file" name="data">
        <p>
        Soort document (<a href="javascript:void(0)" onclick="javascript:menu(9999)">uitleg</a>):<br>
	<select name="how">
	  <option value="auto">Automatisch bepaald of ander formaat</option>
	  <option value="run">Doorlopende tekst</option>
	  <option value="line">Een zin per regel</option>
	  <option value="line-lbl">Een zin per regel, met labels</option>
	  <option value="line-tok">Een zin per regel, getokeniseerd</option>
	  <option value="line-lbl-tok">Een zin per regel, met labels, getokeniseerd</option>
	</select>
      <p>
    Optioneel, toelichting:<br>
    <textarea rows="6" cols="80" id="infotext" name="infotext" maxlength="4000" placeholder="tekst in markdown-formaat"></textarea>
      <p>
	<input type="submit">
    </form>
</body>
</html>
`, MAXTITLELEN+MAXTITLELEN/4, MAXTITLELEN)

}

func submitCorpus(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	how := firstf(q.form, "how")
	title := maxtitlelen(firstf(q.form, "title"))
	info := firstf(q.form, "infotext")

	if title == "" {
		http.Error(q.w, "Titel ontbreekt", http.StatusPreconditionFailed)
		return
	}

	dirname, fulldirname, ok := beginNewCorpus(q, title, hErr)
	if !ok {
		return
	}

	if len(q.form.File["data"]) > 0 {
		fpout, err := os.Create(filepath.Join(fulldirname, "data"))
		if hErr(q, err) {
			return
		}
		defer fpout.Close()
		fpin, err := q.form.File["data"][0].Open()
		if hErr(q, err) {
			return
		}
		defer fpin.Close()
		_, err = io.Copy(fpout, fpin)
		if hErr(q, err) {
			return
		}
	}

	newCorpus(q, dirname, title, info, how, 0, hErr, true)
}

func newCorpus(q *Context, dirname, title, info, how string, protected int, errCheck func(*Context, error) bool, htmlOutput bool) {

	_, err := sqlDB.Exec(fmt.Sprintf(
		"UPDATE %s_info SET `description` = %q, `owner` = %q, `status` = \"QUEUED\", `params` = %q, `msg` = %q, `protected` = %d WHERE `id` = %q;",
		Cfg.Prefix,
		title, q.user, how, "Bron: "+invoertabel[how], protected,
		dirname))
	if errCheck(q, err) {
		return
	}

	logf("QUEUED: " + dirname)
	p := &Process{
		id:     dirname,
		info:   strings.Replace(info, "\r\n", "\n", -1),
		chKill: make(chan bool),
		queued: true,
	}
	processLock.Lock()
	taskWaitNr++
	p.nr = taskWaitNr
	processes[dirname] = p
	processLock.Unlock()
	chWork <- p

	if htmlOutput {
		writeHtml(
			q,
			"Document wordt verwerkt",
			`
Je document wordt verwerkt. Als het klaar is zie je op de hoofdpagina een nieuw corpus bij de databases.
<p>
Let op: Dit kan even duren. Minuten, uren, of dagen, afhankelijk van de grootte van je document.
<p>
<b>Je krijgt een e-mail als het corpus klaar is.</b>
`)
	}
}

func beginNewCorpus(q *Context, title string, errCheck func(*Context, error) bool) (dirname, fulldirname string, ok bool) {

	// db is niet altijd gelijk aan q.db

	dirname = reNoAz.ReplaceAllString(strings.ToLower(title), "")
	if len(dirname) > 20 {
		dirname = dirname[:20]
	} else if dirname == "" {
		dirname = "a"
	}

	dirnameLock.Lock()
	defer dirnameLock.Unlock()
	for i := 0; true; i++ {
		d := dirname + abc(i)
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT 1 FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, d))
		if errCheck(q, err) {
			return
		}
		if rows.Next() {
			rows.Close()
			continue
		}
		dirname = d
		break
	}
	fulldirname = filepath.Join(paqudatadir, "data", dirname)
	err := os.Mkdir(fulldirname, 0700)
	if errCheck(q, err) {
		return
	}

	_, err = sqlDB.Exec(fmt.Sprintf(
		"INSERT %s_info (`id`,`description`,`msg`,`params`) VALUES (%q,\"\",\"\",\"\");",
		Cfg.Prefix,
		dirname))

	if errCheck(q, err) {
		return
	}

	return dirname, fulldirname, true
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

func countXML(dir string) int {
	fis, err := ioutil.ReadDir(dir)
	if err != nil || len(fis) == 0 {
		return 0
	}

	count := 0

	name := fis[len(fis)-1].Name()
	if files, err := ioutil.ReadDir(dir + "/" + name); err == nil {
		for _, f := range files {
			if n := f.Name(); strings.HasSuffix(n, ".xml") || strings.HasSuffix(n, ".xml.gz") {
				count++
			}
		}
	}

	if count == 0 {
		time := fis[0].ModTime()
		name = fis[0].Name()
		for _, fi := range fis {
			if t := fi.ModTime(); time.Before(t) {
				time = t
				name = fi.Name()
			}
		}

		if files, err := ioutil.ReadDir(dir + "/" + name); err == nil {
			for _, f := range files {
				if n := f.Name(); strings.HasSuffix(n, ".xml") || strings.HasSuffix(n, ".xml.gz") {
					count++
				}
			}
		}
	}

	for _, fi := range fis {
		if n := fi.Name(); n < name {
			count += 10000
		}
	}

	return count
}
