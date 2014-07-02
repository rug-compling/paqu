package main

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func download(q *Context) {
	id := first(q.r, "id")
	dl := first(q.r, "dl")

	if !q.prefixes[id] {
		// misschien een corpus dat mislukt is
		rows, err := q.db.Query(
			fmt.Sprintf("SELECT 1 FROM `%s_info` WHERE `id` = %q AND `owner` = %q",
				Cfg.Prefix,
				id,
				q.user))
		if err != nil {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
		if rows.Next() {
			rows.Close()
		} else {
			http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
			return
		}
	}

	datadir := path.Join(paqudir, "data", id)
	var filename string
	switch dl {
	case "summary":
		filename = "summary.txt"
	case "stdout":
		filename = "stdout.txt"
	case "stderr":
		filename = "stderr.txt"
	case "zinnen":
		filename = "data.lines"
	case "dact":
		filename = "data.dact"
	case "xml":
	default:
		http.Error(q.w, "Ongeldige selectie: "+dl, http.StatusUnauthorized)
		return
	}

	if filename != "" {
		fp, err := os.Open(path.Join(datadir, filename+".gz"))
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
		if filename == "data.dact" {
			q.w.Header().Set("Content-Type", "application/octet-stream")
			q.w.Header().Set("Content-Disposition", "attachment; filename="+id+".dact")
		} else {
			q.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		io.Copy(q.w, r)
		r.Close()
		fp.Close()
		return
	}

	// xml
	datadir = path.Join(datadir, "xml")
	files, err := ioutil.ReadDir(datadir)
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	q.w.Header().Set("Content-Type", "application/zip")
	q.w.Header().Set("Content-Disposition", "attachment; filename="+id+".zip")

	w := zip.NewWriter(q.w)
	for _, file := range files {
		gzname := file.Name()
		name := gzname[:len(gzname)-3]
		fh, err := zip.FileInfoHeader(file)
		if err != nil {
			logerr(err)
			return
		}
		fh.Name = path.Join(id, decode_filename(name))
		f, err := w.CreateHeader(fh)
		if err != nil {
			logerr(err)
			return
		}

		fp, err := os.Open(path.Join(datadir, gzname))
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
