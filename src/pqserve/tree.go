package main

//. C-code voor omzetten van dot naar svg

/*
#cgo LDFLAGS: -lgvc -lcgraph
#include <graphviz/gvc.h>
#include <graphviz/cgraph.h>
#include <stdlib.h>

char *makeGraph(char *data) {
	Agraph_t *G;
	char *s;
	unsigned int n;
	GVC_t *gvc;

	s = NULL;
	gvc = gvContext();
	G = agmemread(data);
	free(data);
	if (G == NULL) {
		gvFreeContext(gvc);
		return s;
	}
	gvLayout(gvc, G, "dot");
	gvRenderData(gvc, G, "svg", &s, &n);
	gvFreeLayout(gvc, G);
	agclose(G);
	gvFreeContext(gvc);

	return s;
}
*/
import "C"

//. Imports

import (
	"github.com/pebbe/compactcorpus"

	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

//. Variables

var (
	treeMu sync.Mutex
)

//. Functies

func tree(q *Context) {

	// Zeldzame crash toen er zo te zien twee bomen tegelijk getekend werden.
	// Is graphviz soms niet thread-safe?
	treeMu.Lock()
	defer treeMu.Unlock()

	ctx := &TreeContext{
		yellow: make(map[int]bool),
		green:  make(map[int]bool),
		marks:  make(map[string]bool),
		refs:   make(map[string]bool),
		mnodes: make(map[string]bool),
		words:  make([]string, 0),
	}

	var data []byte
	var dot bool

	prefix := first(q.r, "db")
	if prefix == "" {
		http.Error(q.w, "Geen corpus in query", http.StatusPreconditionFailed)
		return
	}
	if !q.prefixes[prefix] {
		http.Error(q.w, "Ongeldig corpus", http.StatusPreconditionFailed)
		return
	}

	has_names := false
	if first(q.r, "names") == "true" {
		has_names = true
	}

	if n := first(q.r, "marknodes"); n != "" {
		for _, m := range strings.Split(n, ",") {
			ctx.mnodes[m] = true
		}
	}

	file := 0
	arch := 0
	if !has_names {
		var err error
		file, err = strconv.Atoi(first(q.r, "file"))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusPreconditionFailed)
			return
		}
		arch, err = strconv.Atoi(first(q.r, "arch"))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusPreconditionFailed)
			return
		}
	}

	label := ""
	if has_names {
		label = first(q.r, "file")
		if first(q.r, "global") == "true" {
			label = filepath.Base(first(q.r, "arch")) + "::" + label
		}
	} else {
		rows, err := q.db.Query(fmt.Sprintf("SELECT `lbl` FROM `%s_c_%s_sent` WHERE `file` = %d AND `arch` = %d", Cfg.Prefix, prefix, file, arch))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			err := rows.Scan(&label)
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			rows.Close()
		}
	}

	archive := ""
	if has_names {
		archive = first(q.r, "arch")
	} else {
		if arch >= 0 {
			rows, err := q.db.Query(fmt.Sprintf("SELECT arch FROM %s_c_%s_arch WHERE id = %d", Cfg.Prefix, prefix, arch))
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			if rows.Next() {
				err := rows.Scan(&archive)
				if err != nil {
					http.Error(q.w, err.Error(), http.StatusInternalServerError)
					logerr(err)
					return
				}
				rows.Close()
			}
		}
	}

	filename := ""
	if has_names {
		filename = first(q.r, "file")
	} else {
		rows, err := q.db.Query(fmt.Sprintf("SELECT file FROM %s_c_%s_file WHERE id = %d", Cfg.Prefix, prefix, file))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			err := rows.Scan(&filename)
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			rows.Close()
		}
	}

	// xml-bestand inlezen
	if archive != "" {
		if strings.HasSuffix(archive, ".dact") {
			var err error
			data, err = get_dact(archive, filename)
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
		} else {
			reader, err := compactcorpus.RaOpen(archive)
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			data, err = reader.Get(filename)
			reader.Close()
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
		}
	} else {
		var err error
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			fp, err := os.Open(filename + ".gz")
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			r, err := gzip.NewReader(fp)
			if err != nil {
				fp.Close()
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
			data, err = ioutil.ReadAll(r)
			r.Close()
			fp.Close()
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
		}
	}

	if first(q.r, "xml") != "" {
		q.w.Header().Set("Content-type", "application/xml")
		fmt.Fprint(q.w, string(data))
		return
	}

	// markeringen voor gewone woorden
	ctx.yellow = indexes(first(q.r, "yl"))

	// markeringen voor hoofdwoorden
	ctx.green = indexes(first(q.r, "gr"))

	// markeringen van nodes
	for _, m := range strings.Split(first(q.r, "ms"), ",") {
		ctx.marks[m] = true
	}

	// uitvoer als dot i.p.v. html ?
	if first(q.r, "dot") != "" {
		dot = true
	}

	alpino := Alpino_ds{}
	err := xml.Unmarshal(data, &alpino)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	title := html.EscapeString(alpino.Sentence)
	ctx.words = strings.Fields(title)

	// Multi-word units "ineenvouwen".
	// Past ook 'words' aan, bijv: "Koninkrijk" "der" "Nederlanden" -> "Koninkrijk der Nederlanden" "" ""
	if first(q.r, "mwu") != "false" {
		mwu(ctx, alpino.Node0)
	}

	// Uitvoer content-type.
	// Als niet 'dot', dan ook de kop van het HTML-bestand, inclusief zin met gekleurde woorden.
	if dot {
		q.w.Header().Set("Content-type", "application/x-dot; charset=utf-8")
		q.w.Header().Set("Content-disposition", "attachment; filename=tree.dot")
	} else {
		q.w.Header().Set("Content-type", "text/html; charset=utf-8")

		for i, w := range ctx.words {
			if ctx.yellow[i] {
				if ctx.green[i] {
					ctx.words[i] = "<span style=\"background-color: #00ffff;\">" + w + "</span>"
				} else {
					ctx.words[i] = "<span style=\"background-color: #ffff00;\">" + w + "</span>"
				}
			} else if ctx.green[i] {
				ctx.words[i] = "<span style=\"background-color: #90ee90;\">" + w + "</span>"
			}
		}
		fmt.Fprintf(q.w, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="noindex,nofollow">
<title>PaQu: %s</title>
<link rel="stylesheet" type="text/css" href="tooltip.css" />
<script type="text/javascript" src="tooltip.js"></script>
</head>
<body>
<em>%s</em>
<p>
`, title, strings.Join(ctx.words, " "))

		if len(alpino.Meta) > 0 {
			for _, m := range alpino.Meta {
				var v string
				if m.Type == "date" {
					t, err := time.Parse("2006-01-02", m.Value)
					if err != nil {
						v = err.Error()
					} else {
						v = printDate(t, false)
					}
				} else if m.Type == "datetime" {
					t, err := time.Parse("2006-01-02 15:04", m.Value)
					if err != nil {
						v = err.Error()
					} else {
						v = printDate(t, true)
					}
				} else {
					v = m.Value
				}
				fmt.Fprintf(q.w, "%s: %s<br>\n", html.EscapeString(m.Name), html.EscapeString(v))
			}
			fmt.Fprintln(q.w, "<p>")
		}
	}

	// BEGIN: definitie van dot-bestand aanmaken.

	ctx.graph.WriteString(`strict graph gr {

    ranksep=".25 equally"
    nodesep=.05
    ordering=out

    node [shape=plaintext, height=0, width=0, fontsize=12, fontname="Helvetica"];

`)

	// Registreer alle markeringen van nodes met een verwijzing.
	set_refs(ctx, alpino.Node0)

	// Nodes
	print_nodes(q, ctx, alpino.Node0)

	// Terminals
	ctx.graph.WriteString("\n    node [fontname=\"Helvetica-Oblique\", shape=box, color=\"#d3d3d3\", style=filled];\n\n")
	ctx.start = 0
	terms := print_terms(ctx, alpino.Node0)
	sames := strings.Split(strings.Join(terms, " "), "|")
	for _, same := range sames {
		same = strings.TrimSpace(same)
		if same != "" {
			ctx.graph.WriteString("\n    {rank=same; " + same + " }\n")
		}
	}

	// Edges
	ctx.graph.WriteString("\n    edge [sametail=true, color=\"#d3d3d3\"];\n\n")
	print_edges(ctx, alpino.Node0)

	ctx.graph.WriteString("}\n")

	// EINDE: definitie van dot-bestand aanmaken.

	// Als gebruiker dot wil, dan print dot, en exit.
	if dot {
		fmt.Fprint(q.w, ctx.graph.String())
		return
	}

	// Dot omzetten naar svg.
	// De C-string wordt in de C-code ge'free'd.
	s := C.makeGraph(C.CString(ctx.graph.String()))
	svg := C.GoString(s)
	C.free(unsafe.Pointer(s))

	// BEGIN: svg nabewerken en printen

	// XML-declaratie en DOCtype overslaan
	if i := strings.Index(svg, "<svg"); i < 0 {
		logerr(errors.New(fmt.Sprintf("BUG: %v %v", q.r.Method, q.r.URL)))
	} else {
		svg = svg[i:]
	}

	a := ""
	for _, line := range strings.SplitAfter(svg, "\n") {
		// alles wat begint met <title> weghalen
		i := strings.Index(line, "<title")
		if i >= 0 {
			line = line[:i] + "\n"
		}

		// <a xlink> omzetten in tooltip
		i = strings.Index(line, "<a xlink")
		if i >= 0 {
			s := line[i:]
			line = line[:i] + "\n"
			i = strings.Index(s, "\"")
			s = s[i+1:]
			i = strings.LastIndex(s, "\"")
			a = strings.TrimSpace(s[:i])

		}
		if strings.HasPrefix(line, "<text ") && a != "" {
			line = "<text onmouseover=\"tooltip.show('" + html.EscapeString(a) + "')\" onmouseout=\"tooltip.hide()\"" + line[5:]
		}
		if strings.HasPrefix(line, "</a>") {
			line = ""
			a = ""
		}

		fmt.Fprint(q.w, line)
	}
	// EIND: svg nabewerken en printen

	fmt.Fprintf(q.w, "<p>\nopslaan als: <a href=\"tree?%s&amp;dot=1\">dot</a><p>\n", q.r.URL.RawQuery)

	fmt.Fprintf(q.w, "bestand: <a href=\"tree?%s&amp;xml=1\">%s</a>\n", q.r.URL.RawQuery, html.EscapeString(label))

	fmt.Fprint(q.w, "\n</body>\n</html>\n")

}

//. Multi-word units "ineenvouwen".

// Deze code verschilt aan die in wordrel.go doordat hier niet alleen 'node' wordt
// aangepast, maar ook de globale variabele 'words'
// Hier wordt node.Root niet aangepast, omdat Root niet gebruikt wordt.
// Verder worden hier postags verzameld voor de tooltip, waarbij Lemma leeg gemaakt wordt.
func mwu(ctx *TreeContext, node *Node) {
	if node.Cat == "mwu" {
		/*
			Voorwaardes:
			- De dochters moeten op elkaar aansluiten, en heel het bereik van de parentnode beslaan.
			- Er mogen geen indexen in gebruikt zijn (met of zonder inhoud)
		*/
		p1 := node.Begin
		p2 := node.End
		ok := true
		for _, n := range node.NodeList {
			if n.Index != "" {
				ok = false
				break
			}
			if n.Begin != p1 {
				ok = false
				break
			}
			p1 = n.End
		}
		if ok && p1 == p2 {
			node.Cat = ""
			node.Pt = "mwu"
			wrds := make([]string, 0, node.End-node.Begin) // not 'words' !
			postags := make([]string, 0, node.End-node.Begin)
			for _, n := range node.NodeList {
				wrds = append(wrds, n.Word)
				postags = append(postags, fmt.Sprintf("%s:%s", n.Lemma, n.Postag))
			}
			node.Word = strings.Join(wrds, " ")
			node.Lemma = ""
			node.Postag = strings.Join(postags, " ")
			node.NodeList = node.NodeList[0:0]

			// aanpassing 'words', code die ontbreekt in wordrel.go
			ctx.words[node.Begin] = strings.Join(ctx.words[node.Begin:node.End], " ")
			for i := node.Begin + 1; i < node.End; i++ {
				ctx.words[i] = ""
			}
		}
	}
	for _, n := range node.NodeList {
		mwu(ctx, n)
	}

}

//. Genereren van dot

// Registreer alle markeringen van nodes met een verwijzing:
// Als node X gemarkeerd is, en deze node wijst naar node X, dan deze node ook
// markeren voor weergave in vet (maar niet voor een vette egde naar de parent node).
// (Een node met inhoud EN index verwijst naar zichzelf.)
func set_refs(ctx *TreeContext, node *Node) {
	if node.Index != "" && ctx.marks[node.Id] {
		ctx.refs[node.Index] = true
	}
	for _, d := range node.NodeList {
		set_refs(ctx, d)
	}
}

func print_nodes(q *Context, ctx *TreeContext, node *Node) {
	idx := ""
	style := ""

	// Als dit een node met inhoud EN index is, dan in vierkant zetten.
	// Als de node gemarkeerd is, dan in zwart, anders in lichtgrijs.
	// Index als nummer in label zetten.
	if (len(node.NodeList) > 0 || node.Word != "") && node.Index != "" {
		idx = fmt.Sprintf("%s\\n", node.Index)
		if !ctx.refs[node.Index] {
			style += ", color=\"#d3d3d3\""
		} else {
			style += ", color=\"#000000\""
		}
		if ctx.mnodes[node.Id] {
			style += ", style=filled, fillcolor=\"#7FCDBB\""
		} else {
			style += ", shape=box"
		}
	} else if ctx.mnodes[node.Id] {
		style += ", color=\"#7FCDBB\", style=filled"
	}

	// attributen
	var tooltip bytes.Buffer
	tooltip.WriteString("<table class=\"attr\">")
	for _, attr := range NodeTags {
		if value := getAttr(attr, &node.FullNode); value != "" {
			tooltip.WriteString(fmt.Sprintf("<tr><td class=\"lbl\">%s:<td>%s", html.EscapeString(attr), html.EscapeString(value)))
		}
	}
	tooltip.WriteString("</table>")

	lbl := dotquote(node.Rel)
	if node.Cat != "" && node.Cat != node.Rel {
		lbl += "\\n" + dotquote(node.Cat)
	} else if node.Pt != "" && node.Pt != node.Rel {
		lbl += "\\n" + dotquote(node.Pt)
	}

	ctx.graph.WriteString(fmt.Sprintf("    n%v [label=\"%v%v\"%s, tooltip=\"%s\"];\n", node.Id, idx, lbl, style, dotquote2(tooltip.String())))
	for _, d := range node.NodeList {
		print_nodes(q, ctx, d)
	}
}

// Geeft een lijst terminals terug die op hetzelfde niveau moeten komen te staan,
// met "|" ingevoegd voor onderbrekingen in niveaus.
func print_terms(ctx *TreeContext, node *Node) []string {
	terms := make([]string, 0)

	if len(node.NodeList) == 0 {
		if node.Word != "" {
			// Een terminal
			idx := ""
			col := ""
			if ctx.yellow[node.Begin] {
				if ctx.green[node.Begin] {
					col = ", color=\"#00ffff\""
				} else {
					col = ", color=\"#ffff00\""
				}
			} else if ctx.green[node.Begin] {
				col = ", color=\"#90ee90\""
			}
			if node.Begin != ctx.start {
				// Onderbeking
				terms = append(terms, "|")
				// Onzichtbare node invoegen om te scheiden van node die links staat
				ctx.graph.WriteString(fmt.Sprintf("    e%v [label=\" \", tooltip=\" \", style=invis];\n", node.Id))
				terms = append(terms, fmt.Sprintf("e%v", node.Id))
				node.skip = true
			}
			ctx.start = node.End
			terms = append(terms, fmt.Sprintf("t%v", node.Id))
			if node.Lemma == "" {
				ctx.graph.WriteString(fmt.Sprintf("    t%v [label=\"%s%s\", tooltip=\"%s\"%s];\n",
					node.Id, idx, dotquote(node.Word), dotquote2(node.Postag), col))
			} else {
				ctx.graph.WriteString(fmt.Sprintf("    t%v [label=\"%s%s\", tooltip=\"%s:%s\"%s];\n",
					node.Id, idx, dotquote(node.Word), dotquote2(node.Lemma), dotquote(node.Postag), col))
			}
		} else {
			// Een lege node met index
			col := ""
			if ctx.marks[node.Id] {
				col = ", color=black"
			}
			ctx.graph.WriteString(fmt.Sprintf("    t%v [fontname=\"Helvetica\", label=\"%s\", style=solid%s];\n", node.Id, node.Index, col))
		}
	} else {
		for _, d := range node.NodeList {
			t := print_terms(ctx, d)
			terms = append(terms, t...)
		}
	}
	return terms
}

func print_edges(ctx *TreeContext, node *Node) {
	if len(node.NodeList) == 0 {
		if node.skip {
			// Extra: Onzichtbare edge naar extra onzichtbare terminal
			ctx.graph.WriteString(fmt.Sprintf("    n%v -- e%v [style=invis];\n", node.Id, node.Id))
		}
		if ctx.marks[node.Id] || (ctx.refs[node.Index] && node.Word != "") {
			// Vette edge naar terminal
			ctx.graph.WriteString(fmt.Sprintf("    n%v -- t%v [color=black];\n", node.Id, node.Id))
		} else {
			// Gewone edge naar terminal
			ctx.graph.WriteString(fmt.Sprintf("    n%v -- t%v;\n", node.Id, node.Id))
		}
	} else {
		// Edges naar dochters
		for _, d := range node.NodeList {
			if ctx.marks[d.Id] {
				// Vette edge naar dochter
				ctx.graph.WriteString(fmt.Sprintf("    n%v -- n%v [color=black];\n", node.Id, d.Id))
			} else {
				// Gewone edge naar dochter
				ctx.graph.WriteString(fmt.Sprintf("    n%v -- n%v;\n", node.Id, d.Id))
			}
		}
		for _, d := range node.NodeList {
			print_edges(ctx, d)
		}
	}
}

func dotquote(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	return s
}

func dotquote2(s string) string {
	s = strings.Replace(s, "\\", "\\\\\\\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	return s
}

// Zet lijst van indexen (string met komma's) om in map[int]bool
func indexes(s string) map[int]bool {
	m := make(map[int]bool)
	for _, i := range strings.Split(s, ",") {
		j, err := strconv.Atoi(i)
		if err == nil {
			m[j] = true
		}
	}
	return m
}
