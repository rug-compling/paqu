package main

import (
	"github.com/BurntSushi/toml"

	"archive/zip"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FoliaSettings struct {
	DataInfo  string
	DataCount int
	DataMulti bool
	MetaInfo  string

	Tokenized bool

	LabelMeta    string
	UseLabelMeta bool

	Items []FoliaItem
}

type FoliaItem struct {
	Label  string
	Type   string
	Source string
	Value  string
	Use    bool
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

	fdir := foliadir(q)
	if doErr(q, os.MkdirAll(fdir, 0700)) {
		return
	}

	settingsChanged := false
	settings := FoliaSettings{
		Tokenized:    true,
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

	action := ""
	idxdel := 0
	if q.r.Method == "POST" {
		action = firstf(q.form, "act")
		switch action {
		case "putdata", "putmeta":
			settingsChanged = foliaputfile(q, &settings, fdir)
		case "save":
			settingsChanged = foliasave(q, &settings)
		case "delete":
			settingsChanged, idxdel = foliadelete(q, &settings)
		case "new":
			settingsChanged = folianew(q, &settings)
		case "test":
			var err error
			settingsChanged, err = foliatest(q, &settings, fdir)
			if doErr(q, err) {
				return
			}
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
function markeer(i) {
    var x = document.forms["settings"];
    if (x["label" + i].value == "" || x["value" + i].value == "") {
        x["use" + i].checked = false;
    } else {
        x["use" + i].checked = true;
    }
}
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
function testen() {
    var x = document.forms["settings"];
    x["act"].value = "test";
    x.submit();
}
//--></script>
<h1>Invoer en verwerking van FoLiA met metadata</h1>`)

	if action == "submit" && firstf(q.form, "naam") != "" {
		foliasubmit(q, &settings, fdir)
		return
	}

	fmt.Fprint(q.w, "Huidige data: ")
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

	if settings.DataInfo == "" {
		html_footer(q)
		return
	}

	checked := " checked"

	ch := ""
	if settings.Tokenized {
		ch = checked
	}
	fmt.Fprintf(q.w, `
<form name="settings" action="folia#a" method="post" enctype="multipart/form-data" accept-charset="utf-8">
<input type="hidden" name="act" value="save">
<input type="hidden" name="index" value="0">
<input type="hidden" name="len" value="%d">
Soort invoer:
<div class="foliaform">
<input type="checkbox" name="tokenized" value="true"%s> Getokeniseerd
</div>
`, len(settings.Items), ch)

	ch = ""
	if settings.UseLabelMeta {
		ch = checked
	}
	fmt.Fprintf(q.w, `
Label in uitvoer:
<div class="foliaform">
<input type="checkbox" name="usemeta"%s> Label voor metadatabestand<br>
<input type="text" name="labelmeta" value="%s">
</div>
`,
		ch, html.EscapeString(settings.LabelMeta))

	fmt.Fprintln(q.w, `Metadata (zie <a href="foliavb">voorbeelden</a>):`)

	if idxdel > 0 && idxdel >= len(settings.Items) {
		idxdel = len(settings.Items) - 1
	}
	for i, item := range settings.Items {
		var ch string
		var chs [5]string
		var che [2]string
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
		n = 0
		if item.Source == "id" {
			n = 1
		}
		che[n] = " selected"
		ids := ""
		if action == "delete" && i == idxdel || action == "new" && i == len(settings.Items)-1 {
			ids = ` id="a"`
		}
		fmt.Fprintf(q.w, `
<div  class="foliaform"%s>
<input type="checkbox" name="use%d"%s><br>
Label:
<input type="text" name="label%d" value="%s" onchange="markeer('%d')"><br>
Soort:
<select name="type%d">
  <option%s>text</option>
  <option%s>int</option>
  <option%s>float</option>
  <option%s>date</option>
  <option%s>datetime</option>
</select><br>
<select name="source%d">
  <option value="xpath"%s>Extern XPath:</option>
  <option value="id"%s>Intern ID:</option>
</select>
<input type="text" name="value%d" value="%s" size="80" onchange="markeer('%d')">
<div class="right">
<button type="button" onclick="verwijder(%d)">Verwijderen</button>
</div>
</div>
`,
			ids,
			i, ch,
			i, html.EscapeString(item.Label), i,
			i, chs[0], chs[1], chs[2], chs[3], chs[4],
			i, che[0], che[1],
			i, html.EscapeString(item.Value), i,
			i)
	}

	ids := ""
	if action == "save" {
		ids = ` id="a"`
	}
	fmt.Fprintf(q.w, `
<button type="button" onclick="toevoegen()">Toevoegen</button>
<hr class="wide">
<button type="button" onclick="opslaan()"%s>Opslaan</button>
<button type="button" onclick="testen()">Testen</button>
`, ids)

	fmt.Fprint(q.w, `
</form>
`)

	if action == "test" || action == "submit" {
		n, m := foliamax(&settings)

		fmt.Fprint(q.w, `
<hr class="wide">
<span id="a">Test:</span> bestandsnamen en eventuele fouten`)
		if settings.DataCount > n {
			fmt.Fprintf(q.w, ", maximaal %d bestanden", n)
		}
		fmt.Fprint(q.w, `
<div class="output">
`)
		data, err := ioutil.ReadFile(filepath.Join(fdir, "test.err"))
		if doErr(q, err) {
			return
		}
		fmt.Fprint(q.w, html.EscapeString(strings.Replace(string(data), filepath.Join(fdir, "data"), ".", -1)))

		if settings.DataCount == 1 {
			fmt.Fprintf(q.w, `
</div>
Test: uitvoer, maximaal %d zinnen
<div class="output">
`, m)
		} else if settings.DataCount <= n {
			fmt.Fprintf(q.w, `
</div>
Test: uitvoer, maximaal %d zinnen per bestand
<div class="output">
`, m)
		} else {
			fmt.Fprintf(q.w, `
</div>
Test: uitvoer, maximaal %d bestanden, %d zinnen per bestand
<div class="output">
`, n, m)
		}
		data, err = ioutil.ReadFile(filepath.Join(fdir, "test.out"))
		if doErr(q, err) {
			return
		}
		fmt.Fprint(q.w, html.EscapeString(string(data)))
		fmt.Fprint(q.w, `
</div>
Is bovenstaande uitvoer naar wens? Geen onverwachte foutmeldingen? Dan kun je het corpus met deze metadata invoeren.
<div class="foliaform">
<form action="folia#a" method="post" enctype="multipart/form-data" accept-charset="utf-8">
<input type="hidden" name="act" value="submit">
Naam voor corpus: <input type="text" name="naam"><p>
<input type="submit" value="Invoeren">
</form>
</div>
`)
	}

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
			settings.DataCount = 1
			settings.DataMulti = false
		}
		st, err := os.Stat(datafile)
		if doErr(q, err) {
			return
		}
		info = fmt.Sprintf("%s &mdash; %d Kb", html.EscapeString(uploadname), (st.Size()+512)/1024)
		os.Rename(datafile, filepath.Join(datadir, uploadname))
	} else {
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
		if act == "putdata" {
			settings.DataCount = filecount
			settings.DataMulti = true
		}
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

	settings.LabelMeta = firstf(q.form, "labelmeta")
	settings.UseLabelMeta = firstf(q.form, "usemeta") != ""

	settings.Items = settings.Items[0:0]
	n, _ := strconv.Atoi(firstf(q.form, "len"))
	for i := 0; i < n; i++ {
		s := fmt.Sprint(i)
		item := FoliaItem{
			Label:  firstf(q.form, "label"+s),
			Type:   firstf(q.form, "type"+s),
			Source: firstf(q.form, "source"+s),
			Value:  firstf(q.form, "value"+s),
			Use:    firstf(q.form, "use"+s) != "",
		}
		if item.Label == "" || item.Value == "" {
			item.Use = false
		}
		settings.Items = append(settings.Items, item)
	}
	if len(settings.Items) == 0 {
		settings.Items = append(settings.Items, FoliaItem{Type: "text"})
	}

	return
}

func foliadelete(q *Context, settings *FoliaSettings) (settingsChanged bool, index int) {
	settingsChanged = foliasave(q, settings)
	index, _ = strconv.Atoi(firstf(q.form, "index"))
	if index < len(settings.Items) {
		settings.Items = append(settings.Items[:index], settings.Items[index+1:]...)
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

func foliatest(q *Context, settings *FoliaSettings, fdir string) (settingsChanged bool, err error) {
	settingsChanged = foliasave(q, settings)

	os.Remove(filepath.Join(fdir, "error.txt"))

	var fp *os.File
	fp, err = os.Create(filepath.Join(fdir, "config.toml"))
	if err != nil {
		return
	}

	fmt.Fprintln(fp, `File_src = ""`)
	fmt.Fprintln(fp, `File_path = ""`)
	if settings.UseLabelMeta {
		fmt.Fprintf(fp, "Meta_src = %q\n", settings.LabelMeta)
	} else {
		fmt.Fprintln(fp, `Meta_src = ""`)
	}

	fmt.Fprintln(fp, `Output_dir = ""`)

	fmt.Fprint(fp, "Item_list = [")
	p := ""
	for _, item := range settings.Items {
		if item.Use {
			fmt.Fprintf(fp, "%s\n    %q", p, item.Label)
			p = ","
		}
	}
	fmt.Fprintln(fp, "\n]")

	fmt.Fprintf(fp, "Data_dir = \"%s\"\n", filepath.Join(fdir, "data"))
	fmt.Fprintf(fp, "Meta_dir = \"%s\"\n", filepath.Join(fdir, "meta"))

	fmt.Fprintf(fp, "Tokenized = %v\n", settings.Tokenized)

	for _, item := range settings.Items {
		if item.Use {
			p := "XPath"
			if item.Source == "id" {
				p = "ID"
			}
			fmt.Fprintf(fp, `
[Items.%q]
Type = %q
%s = %q
`, item.Label, item.Type, p, item.Value)
		}
	}

	fp.Close()

	n, m := foliamax(settings)

	shell("pqfolia -n %d -m %d %s > %s 2> %s",
		n,
		m,
		filepath.Join(fdir, "config.toml"),
		filepath.Join(fdir, "test.out"),
		filepath.Join(fdir, "test.err")).Run()

	return
}

func foliadir(q *Context) string {
	return filepath.Join(paqudir, "folia", hex.EncodeToString([]byte(q.user)))
}

func foliamax(settings *FoliaSettings) (int, int) {
	n := settings.DataCount
	if n > 20 {
		n = 20
	}
	m := 1000 / n
	return n, m
}

func foliasubmit(q *Context, settings *FoliaSettings, fdir string) {
	fmt.Fprint(q.w, `
<span id="a">Het corpus wordt verwerkt.</span>
<p>
Het kan even duren voordat je corpus verschijnt in je lijst met corpora.
<p>
Ga naar: <a href="corpora">corpora</a>
`)
	html_footer(q)

	go func() {
		foliaMu.Lock()
		if _, ok := foliaUsers[q.user]; !ok {
			foliaUsers[q.user] = new(sync.Mutex)
		}
		foliaMu.Unlock()

		foliaUsers[q.user].Lock()
		defer foliaUsers[q.user].Unlock()

		fp, err := os.Create(filepath.Join(fdir, "config.toml"))
		if foliaErr(q, err) {
			return
		}

		fmt.Fprintln(fp, `File_src = ""`)
		fmt.Fprintln(fp, `File_path = ""`)
		if settings.UseLabelMeta {
			fmt.Fprintf(fp, "Meta_src = %q\n", settings.LabelMeta)
		} else {
			fmt.Fprintln(fp, `Meta_src = ""`)
		}

		outdir := filepath.Join(fdir, "out")
		if settings.DataMulti {
			os.RemoveAll(outdir)
			os.MkdirAll(outdir, 0700)
			fmt.Fprintf(fp, "Output_dir = \"%s\"\n", outdir)
		} else {
			fmt.Fprintln(fp, `Output_dir = ""`)
		}

		fmt.Fprint(fp, "Item_list = [")
		pre := ""
		for _, item := range settings.Items {
			if item.Use {
				fmt.Fprintf(fp, "%s\n    %q", pre, item.Label)
				pre = ","
			}
		}
		fmt.Fprintln(fp, "\n]")

		fmt.Fprintf(fp, "Data_dir = \"%s\"\n", filepath.Join(fdir, "data"))
		fmt.Fprintf(fp, "Meta_dir = \"%s\"\n", filepath.Join(fdir, "meta"))

		fmt.Fprintf(fp, "Tokenized = %v\n", settings.Tokenized)

		for _, item := range settings.Items {
			if item.Use {
				p := "XPath"
				if item.Source == "id" {
					p = "ID"
				}
				fmt.Fprintf(fp, `
[Items.%q]
Type = %q
%s = %q
`, item.Label, item.Type, p, item.Value)
			}
		}

		fp.Close()

		outfile := filepath.Join(fdir, "outfile")

		o := ""
		if !settings.DataMulti {
			o = " > " + outfile
		}
		errfile := filepath.Join(fdir, "pqfolia.err")
		err = shell("pqfolia %s%s 2> %s", filepath.Join(fdir, "config.toml"), o, errfile).Run()
		if foliaErr(q, err) {
			os.Rename(errfile, filepath.Join(foliadir(q), "error.txt"))
			os.RemoveAll(outdir)
			return
		}
		os.Remove(errfile)

		if settings.DataMulti {
			fp, err = os.Create(outfile)
			if foliaErr(q, err) {
				os.RemoveAll(outdir)
				return
			}
			zf := zip.NewWriter(fp)
			ok := foliazipdir(q, zf, outdir, "")
			err = zf.Close()
			fp.Close()
			os.RemoveAll(outdir)
			if !ok {
				return
			}
			if foliaErr(q, err) {
				return
			}
		}

		db, err := dbopen()
		if foliaErr(q, err) {
			return
		}
		defer db.Close()

		title := firstf(q.form, "naam")

		dirname, fulldirname, ok := beginNewCorpus(q, db, title, foliaErr)
		if !ok {
			return
		}

		os.Rename(outfile, filepath.Join(fulldirname, "data"))

		tok := ""
		if settings.Tokenized {
			tok = "-tok"
		}
		newCorpus(q, db, dirname, title, "line-lbl"+tok, 0, foliaErr, false)

	}()
}

func foliazipdir(q *Context, zf *zip.Writer, fdir, subdir string) (ok bool) {
	dir := filepath.Join(fdir, subdir)
	files, err := ioutil.ReadDir(dir)
	if foliaErr(q, err) {
		return false
	}
	for _, file := range files {
		fname := path.Join(subdir, file.Name()) // in zip alleen forward slashes toegestaan
		if file.IsDir() {
			if !foliazipdir(q, zf, fdir, fname) {
				return false
			}
		} else {
			data, err := ioutil.ReadFile(filepath.Join(fdir, fname))
			if foliaErr(q, err) {
				return false
			}
			fh, err := zip.FileInfoHeader(file)
			if foliaErr(q, err) {
				return false
			}
			fh.Name = fname
			f, err := zf.CreateHeader(fh)
			if foliaErr(q, err) {
				return false
			}
			_, err = f.Write(data)
			if foliaErr(q, err) {
				return false
			}
		}
	}
	return true
}

// system error zonder user -> alleen log
func foliaErr(q *Context, err error) bool {
	if err == nil {
		return false
	}

	s := err.Error()

	var s1 string
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		s1 = fmt.Sprintf("FOUT: %v:%v: %v", filepath.Base(filename), lineno, s)
	} else {
		s1 = "FOUT: " + s
	}
	chLog <- s1

	fp, e := os.Create(filepath.Join(foliadir(q), "error.txt"))
	if e == nil {
		defer fp.Close()
		fmt.Fprintln(fp, s)
	}
	return true
}

func foliaclean() {
	fdir := filepath.Join(paqudir, "folia")
	for {
		// clean up
		then := time.Now().AddDate(0, 0, -Cfg.Foliadays)

		files, err := ioutil.ReadDir(fdir)
		if sysErr(err) {
			goto SLEEP
		}
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			if then.After(file.ModTime()) {
				fname := file.Name()
				user, err := hex.DecodeString(fname)
				if sysErr(err) {
					user = []byte(fname)
				}
				chLog <- "Verwijderen FoLiA-bestanden van gebruiker " + string(user)
				sysErr(os.RemoveAll(filepath.Join(fdir, fname)))
			}
		}

	SLEEP:
		// sleep tot na vier uur 's ochtends
		time.Sleep(time.Duration((28 - time.Now().Hour())) * time.Hour)

	}
}
