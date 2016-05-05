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

	//
	// Begin verwerking van upload
	//
	if q.r.Method == "POST" {
		for {
			act := firstf(q.form, "act")
			var datadir string
			if act == "putdata" {
				datadir = filepath.Join(fdir, "data")
			} else if act == "putmeta" {
				datadir = filepath.Join(fdir, "meta")
			} else {
				break
			}
			uploadname := filepath.Base(q.form.File["data"][0].Filename)
			os.RemoveAll(datadir)
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
			settingsChanged = true
			break
		}
	}
	//
	// Einde verwerking van upload
	//

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

	fmt.Fprint(q.w, "<h1>Invoer en verwerking van FoLiA met metadata</h1>\nHuidige data: ")
	if settings.DataInfo == "" {
		fmt.Fprint(q.w, "<em>geen</em>")
	} else {
		fmt.Fprint(q.w, settings.DataInfo)
	}
	fmt.Fprint(q.w, `
<form class="foliafile" action="folia" method="post" enctype="multipart/form-data" accept-charset="utf-8">
<input type="hidden" name="act" value="putdata">
Nieuwe data (folia: .xml/.zip/.tar/.tar.gz):<br>
<input type="file" name="data">
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
Nieuwe metadata (cmdi/imdi/...: .xml/.zip/.tar/.tar.gz):<br>
<input type="file" name="data">
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

	html_footer(q)
}
