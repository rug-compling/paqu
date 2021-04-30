package main

import (
	"github.com/rug-compling/paqu/internal/dir"

	"github.com/pebbe/util"

	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func download(q *Context) {
	id := first(q.r, "id")
	dl := first(q.r, "dl")

	params := q.params[id]
	if !q.myprefixes[id] {
		// misschien een corpus dat mislukt is
		rows, err := sqlDB.Query(
			fmt.Sprintf("SELECT `params` FROM `%s_info` WHERE `id` = %q AND `owner` = %q",
				Cfg.Prefix,
				id,
				q.user))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			err := rows.Scan(&params)
			rows.Close()
			if err != nil {
				http.Error(q.w, err.Error(), http.StatusInternalServerError)
				logerr(err)
				return
			}
		} else {
			http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
			return
		}
	}

	if q.protected[id] && (dl == "xml" || dl == "dact") {
		http.Error(q.w, "Dat is afgeleid van een corpus dat niet van jou is", http.StatusUnauthorized)
		return
	}

	datadir := filepath.Join(dir.Data, "data", id)
	var filename string
	switch dl {
	case "summary":
		filename = "summary.txt"
	case "stdout":
		filename = "stdout.txt"
	case "stderr":
		filename = "stderr.txt"
	case "zinnen":
		if !strings.Contains(params, "-lbl") && !strings.Contains(params, "-arch") && !strings.HasPrefix(params, "folia") && !strings.HasPrefix(params, "tei") {
			filename = "data.lines"
		}
	case "dact":
	case "xml":
	default:
		http.Error(q.w, "Ongeldige selectie: "+dl, http.StatusUnauthorized)
		return
	}

	if filename != "" {
		fp, err := os.Open(filepath.Join(datadir, filename+".gz"))
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
		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.Copy(q.w, r)
		r.Close()
		fp.Close()
		return
	}

	// dact
	if dl == "dact" {
		fp, err := os.Open(filepath.Join(datadir, "data.dact"))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		q.w.Header().Set("Content-Type", "application/octet-stream")
		q.w.Header().Set("Content-Disposition", "attachment; filename="+id+".dact")
		io.Copy(q.w, fp)
		fp.Close()
		return
	}

	// data.lines met verkeerde labels
	if dl == "zinnen" {
		fp, err := os.Open(filepath.Join(datadir, "data.lines.gz"))
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

		q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rd := util.NewReader(r)
		for {
			line, err := rd.ReadLineString()
			if err != nil {
				break
			}
			a := strings.SplitN(line, "|", 2)
			lbl := decode_filename(a[0][1+strings.Index(a[0], "-"):])
			fmt.Fprintf(q.w, "%s|%s\n", lbl, a[1])
		}
		r.Close()
		fp.Close()

		return
	}

	// xml
	datadir = filepath.Join(datadir, "xml")
	files, err := filenames2(datadir, false)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	q.w.Header().Set("Content-Type", "application/zip")
	q.w.Header().Set("Content-Disposition", "attachment; filename="+id+".zip")

	w := zip.NewWriter(q.w)
	for _, gzname := range files {
		fullgzname := filepath.Join(datadir, gzname)
		file, err := os.Stat(fullgzname)
		name := decode_filename(gzname[:len(gzname)-3])
		if params == "dact" || strings.HasPrefix(params, "xmlzip") {
			name = name[1+strings.Index(name, "/"):]
		} else if strings.Contains(params, "-lbl") || strings.Contains(params, "-arch") || strings.HasPrefix(params, "folia") || strings.HasPrefix(params, "tei") {
			name = name[1+strings.Index(name, "-"):]
		}
		if err != nil {
			logerr(err)
			return
		}
		fh, err := zip.FileInfoHeader(file)
		if err != nil {
			logerr(err)
			return
		}
		fh.Name = filepath.Join(id, name)
		f, err := w.CreateHeader(fh)
		if err != nil {
			logerr(err)
			return
		}

		fp, err := os.Open(filepath.Join(datadir, gzname))
		if err != nil {
			logerr(err)
			return
		}
		r, err := gzip.NewReader(fp)
		if err != nil {
			fp.Close()
			logerr(err)
			return
		}

		io.Copy(f, r)
		r.Close()
		fp.Close()
	}
	err = w.Close()
	if err != nil {
		logerr(err)
		return
	}
}
