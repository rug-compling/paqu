package main

import (
	"github.com/BurntSushi/toml"

	"compress/gzip"
	"encoding/hex"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type FoliaSettings struct {
	DataInfo  string
	DataMulti bool
	MetaInfo  string
	WasTested bool
	WasBuild  bool

	Tokenized bool

	OutputZip bool

	LabelFile    string
	UseLabelFile bool
	LabelPath    string
	UseLabelPath bool
	LabelMeta    string
	UseLabelMeta bool

	Items []FoliaItem
}

type FoliaItem struct {
	Label string
	Type  string
	XPath string
	Use   bool
}

var (
	foliaMu    sync.Mutex
	foliaUsers = make(map[string]*sync.Mutex)
)

func foliatool(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	writeHead(q, "FoLiA met metadata", 0)

	foliaMu.Lock()
	if _, ok := foliaUsers[q.user]; !ok {
		foliaUsers[q.user] = new(sync.Mutex)
	}
	foliaMu.Unlock()

	foliaUsers[q.user].Lock()
	defer foliaUsers[q.user].Unlock()

	fdir := filepath.Join(paqudir, "folia", hex.EncodeToString([]byte(q.user)))
	if doErr(q, os.MkdirAll(fdir, 0700)) {
		return
	}

	settingsChanged := false
	settings := FoliaSettings{
		Tokenized:    true,
		LabelFile:    "File.Src",
		UseLabelFile: true,
		LabelPath:    "File.Path.",
		UseLabelPath: true,
		LabelMeta:    "Meta.Src",
		UseLabelMeta: true,
		Items:        []FoliaItem{FoliaItem{Type: "text"}},
	}

	settingsFile := filepath.Join(fdir, "settings.toml")
	if _, err := os.Stat(settingsFile); err == nil {
		_, err := toml.DecodeFile(settingsFile, &settings)
		if doErr(q, err) {
			return
		}
	}

	if q.r.Method == "POST" {
		switch firstf(q.form, "act") {
		case "putdata", "putmeta":
			settingsChanged = foliaputfile(q, &settings, fdir)
		case "save":
			settingsChanged = foliasave(q, &settings)
		case "delete":
			settingsChanged = foliadelete(q, &settings)
		case "new":
			settingsChanged = folianew(q, &settings)
		}
	}

	if settingsChanged {
		fp, err := os.Create(settingsFile)
		if doErr(q, err) {
			return
		}
		defer fp.Close()
		err = toml.NewEncoder(fp).Encode(settings)
		if doErr(q, err) {
			return
		}
	}

	// pagina maken
	fmt.Fprint(q.w, `
<script type="text/javascript"><!--
function verwijder(i) {
    var x = document.forms["settings"];
    x["act"].value = "delete";
    x["index"].value = i;
    x.submit();
}
function toevoegen() {
    var x = document.forms["settings"];
    x["act"].value = "new";
    x.submit();
}
function opslaan() {
    var x = document.forms["settings"];
    x["act"].value = "save";
    x.submit();
}
//--></script>
<h1>Invoer en verwerking van FoLiA met metadata</h1>
Huidige data: `)
	if settings.DataInfo == "" {
		fmt.Fprint(q.w, "<em>geen</em>")
	} else {
		fmt.Fprint(q.w, settings.DataInfo)
	}
	fmt.Fprint(q.w, `
<form class="foliafile" action="folia" method="post" enctype="multipart/form-data" accept-charset="utf-8">
<input type="hidden" name="act" value="putdata">
Nieuwe data (folia: .xml/.xml.gz/.zip/.tar/.tar.gz/.tgz):<br>
<input type="file" name="data" accept=".xml,.xml.gz,.zip,.tar,.tar.gz,.tgz">
<input type="submit">
`)
	if settings.DataInfo != "" {
		fmt.Fprint(q.w, "<BR>LET OP: Upload vervangt oude data!\n")
	}
	fmt.Fprint(q.w, `
</form>
Huidige metadata: `)
	if settings.MetaInfo == "" {
		fmt.Fprint(q.w, "<em>geen</em>")
	} else {
		fmt.Fprint(q.w, settings.MetaInfo)
	}
	fmt.Fprint(q.w, `
<form class="foliafile" action="folia" method="post" enctype="multipart/form-data" accept-charset="utf-8">
<input type="hidden" name="act" value="putmeta">
Nieuwe metadata (cmdi/imdi/...: .xml/.xml.gz/.zip/.tar/.tar.gz/.tgz):<br>
<input type="file" name="data" accept=".xml,.xml.gz,.zip,.tar,.tar.gz,.tgz">
<input type="submit">
`)
	if settings.MetaInfo != "" {
		fmt.Fprint(q.w, "<BR>LET OP: Upload vervangt oude data!\n")
	}
	fmt.Fprint(q.w, `
</form>
`)

	if settings.DataInfo == "" || settings.MetaInfo == "" {
		html_footer(q)
		return
	}

	checked := " checked"

	ch := ""
	if settings.Tokenized {
		ch = checked
	}
	fmt.Fprintf(q.w, `
<form name="settings" action="folia" method="post" enctype="multipart/form-data" accept-charset="utf-8">
<input type="hidden" name="act" value="save">
<input type="hidden" name="index" value="0">
<input type="hidden" name="len" value="%d">
Soort invoer:
<div class="foliaform">
<input type="checkbox" name="tokenized" value="true"%s> Getokeniseerd
</div>
`, len(settings.Items), ch)

	if settings.DataMulti {
		var ch1, ch2 string
		if settings.OutputZip {
			ch2 = checked
		} else {
			ch1 = checked
		}
		fmt.Fprintf(q.w, `
Uitvoer:
<div class="foliaform">
<input type="radio" name="outzip" value="false"%s> Alles in één bestand (tekst)<br>
<input type="radio" name="outzip" value="true"%s> Eén uitvoerbestand per invoerbestand (zip)
</div>
`, ch1, ch2)
	}

	var ch1, ch2, ch3 string
	if settings.UseLabelFile {
		ch1 = checked
	}
	if settings.UseLabelPath {
		ch2 = checked
	}
	if settings.UseLabelMeta {
		ch3 = checked
	}
	fmt.Fprintf(q.w, `
Labels in uitvoer:
<div class="foliaform">
<input type="checkbox" name="usefile"%s> Label voor invoerbestand, zonder path<br>
<input type="text" name="labelfile" value="%s">
</div>

<div class="foliaform">
<input type="checkbox" name="usepath"%s> Prefix van delen van het path van het invoerbestand<br>
<input type="text" name="labelpath" value="%s">
</div>

<div class="foliaform">
<input type="checkbox" name="usemeta"%s> Label voor metadatabestand<br>
<input type="text" name="labelmeta" value="%s">
</div>
`,
		ch1, html.EscapeString(settings.LabelFile),
		ch2, html.EscapeString(settings.LabelPath),
		ch3, html.EscapeString(settings.LabelMeta))

	fmt.Fprintln(q.w, `Metadata (zie <a href="foliavb">voorbeelden</a>):`)

	for i, item := range settings.Items {
		var ch string
		var chs [5]string
		if item.Use {
			ch = checked
		}
		var n int
		switch item.Type {
		case "text":
			n = 0
		case "int":
			n = 1
		case "float":
			n = 2
		case "date":
			n = 3
		case "datetime":
			n = 4
		}
		chs[n] = " selected"
		fmt.Fprintf(q.w, `
<div  class="foliaform">
<input type="checkbox" name="use%d"%s><br>
Label:
<input type="text" name="label%d" value="%s"><br>
Soort:
<select name="type%d">
  <option%s>text</option>
  <option%s>int</option>
  <option%s>float</option>
  <option%s>date</option>
  <option%s>datetime</option>
</select><br>
XPath:
<input type="text" name="xpath%d" value="%s" size="80">
<div class="right">
<button type="button" onclick="verwijder(%d)">Verwijderen</button>
</div>
</div>
`,
			i, ch,
			i, html.EscapeString(item.Label),
			i, chs[0], chs[1], chs[2], chs[3], chs[4],
			i, html.EscapeString(item.XPath),
			i)
	}

	fmt.Fprintln(q.w, `
<button type="button" onclick="toevoegen()">Toevoegen</button>
<hr>
<button type="button" onclick="opslaan()">Opslaan</button>
`)

	fmt.Fprint(q.w, `
</form>
`)

	html_footer(q)
}

func foliaputfile(q *Context, settings *FoliaSettings, fdir string) (settingsChanged bool) {
	act := firstf(q.form, "act")
	var datadir string
	if act == "putdata" {
		datadir = filepath.Join(fdir, "data")
		settings.DataInfo = ""
	} else if act == "putmeta" {
		datadir = filepath.Join(fdir, "meta")
		settings.MetaInfo = ""
	} else {
		return
	}
	settingsChanged = true
	os.RemoveAll(datadir)
	if len(q.form.File["data"]) < 1 {
		return
	}
	uploadname := filepath.Base(q.form.File["data"][0].Filename)
	os.MkdirAll(datadir, 0700)
	datafile := filepath.Join(fdir, "datafile")
	fpout, err := os.Create(datafile)
	if doErr(q, err) {
		return
	}
	fpin, err := q.form.File["data"][0].Open()
	if doErr(q, err) {
		fpout.Close()
		return
	}
	_, err = io.Copy(fpout, fpin)
	fpin.Close()
	fpout.Close()
	if doErr(q, err) {
		return
	}

	// gzip
	var fp *os.File
	fp, err = os.Open(datafile)
	if doErr(q, err) {
		return
	}
	b := make([]byte, 2)
	io.ReadFull(fp, b)
	fp.Close()
	if string(b) == "\x1F\x8B" {
		// gzip
		fpin, _ := os.Open(datafile)
		r, err := gzip.NewReader(fpin)
		if doErr(q, err) {
			fpin.Close()
			return
		}
		fpout, _ := os.Create(datafile + ".tmp")
		_, err = io.Copy(fpout, r)
		fpout.Close()
		r.Close()
		fpin.Close()
		if doErr(q, err) {
			return
		}
		os.Rename(datafile+".tmp", datafile)
	}

	var ar *arch
	ar, err = NewArchReader(datafile)
	var info string
	if err != nil {
		if act == "putdata" {
			settings.DataMulti = false
			settings.UseLabelFile = false
			settings.UseLabelPath = false
		}
		st, err := os.Stat(datafile)
		if doErr(q, err) {
			return
		}
		info = fmt.Sprintf("%s &mdash; %d Kb", html.EscapeString(uploadname), (st.Size()+512)/1024)
		os.Rename(datafile, filepath.Join(datadir, uploadname))
	} else {
		if act == "putdata" {
			settings.DataMulti = true
			settings.UseLabelFile = true
			settings.UseLabelPath = true
		}
		filecount := 0
		var filesize int64
		for {
			err := ar.Next()
			if err == io.EOF {
				break
			}
			if doErr(q, err) {
				ar.Close()
				fp.Close()
				return
			}
			filename := ar.Name()
			newfile := filepath.Join(datadir, filename)
			os.MkdirAll(filepath.Join(datadir, filepath.Dir(filename)), 0700)
			fp, err := os.Create(newfile)
			if doErr(q, err) {
				return
			}
			err = ar.Copy(fp)
			fp.Close()
			if doErr(q, err) {
				return
			}
			filecount++
			st, err := os.Stat(newfile)
			if doErr(q, err) {
				return
			}
			filesize += st.Size()
		}
		if filecount == 1 {
			info = fmt.Sprintf("%s &mdash; 1 bestand &mdash; %d Kb", html.EscapeString(uploadname), (filesize+512)/1024)
		} else {
			info = fmt.Sprintf("%s &mdash; %d bestanden &mdash; %d Kb", html.EscapeString(uploadname), filecount, (filesize+512)/1024)
		}
		os.Remove(datafile)
	}
	if act == "putdata" {
		settings.DataInfo = info
	} else {
		settings.MetaInfo = info
	}
	return
}

func foliasave(q *Context, settings *FoliaSettings) (settingsChanged bool) {
	settingsChanged = true

	settings.Tokenized = firstf(q.form, "tokenized") != ""

	settings.OutputZip = firstf(q.form, "outzip") == "true"

	settings.LabelFile = firstf(q.form, "labelfile")
	settings.UseLabelFile = firstf(q.form, "usefile") != ""
	settings.LabelPath = firstf(q.form, "labelpath")
	settings.UseLabelPath = firstf(q.form, "usepath") != ""
	settings.LabelMeta = firstf(q.form, "labelmeta")
	settings.UseLabelMeta = firstf(q.form, "usemeta") != ""

	settings.Items = settings.Items[0:0]
	n, _ := strconv.Atoi(firstf(q.form, "len"))
	for i := 0; i < n; i++ {
		s := fmt.Sprint(i)
		settings.Items = append(settings.Items, FoliaItem{
			Label: firstf(q.form, "label"+s),
			Type:  firstf(q.form, "type"+s),
			XPath: firstf(q.form, "xpath"+s),
			Use:   firstf(q.form, "use"+s) != "",
		})
	}
	if len(settings.Items) == 0 {
		settings.Items = append(settings.Items, FoliaItem{Type: "text"})
	}

	return
}

func foliadelete(q *Context, settings *FoliaSettings) (settingsChanged bool) {
	settingsChanged = foliasave(q, settings)
	n, _ := strconv.Atoi(firstf(q.form, "index"))
	if n < len(settings.Items) {
		settings.Items = append(settings.Items[:n], settings.Items[n+1:]...)
		settingsChanged = true
	}
	if len(settings.Items) == 0 {
		settings.Items = append(settings.Items, FoliaItem{Type: "text"})
		settingsChanged = true
	}
	return
}

func folianew(q *Context, settings *FoliaSettings) (settingsChanged bool) {
	settingsChanged = foliasave(q, settings)
	settings.Items = append(settings.Items, FoliaItem{Type: "text"})
	return
}
