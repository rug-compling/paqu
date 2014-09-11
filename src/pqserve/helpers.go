package main

import (
	_ "github.com/go-sql-driver/mysql"

	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	reFilechars = regexp.MustCompile("[^-._a-zA-Z0-9]+")
	reFilecodes = regexp.MustCompile("_[0-9A-F][0-9A-F]|__")
)

func doErr(q *Context, err error) bool {
	if err == nil {
		return false
	}

	s := err.Error()

	var s1 string
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		s1 = fmt.Sprintf("FOUT: %v:%v: %v", path.Base(filename), lineno, s)
	} else {
		s1 = "FOUT: " + s
	}
	chLog <- s1

	fmt.Fprintln(q.w, "<pre>")
	fmt.Fprintln(q.w, html.EscapeString(s))
	fmt.Fprintln(q.w, "</pre>\n</body>\n</html>")

	return true
}

func writeHead(q *Context, title string, tab int) {
	q.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	q.w.Header().Set("Cache-Control", "no-cache")
	q.w.Header().Add("Pragma", "no-cache")
	if title == "" {
		title = "PaQu"
	} else {
		title = "PaQu -- " + title
	}
	fmt.Fprintf(q.w, `<!DOCTYPE html>
<html>
<head>
<title>%s</title>
<link rel="stylesheet" type="text/css" href="paqu.css">
<meta name="robots" content="noindex,nofollow">
</head>
<body>
<div id="login">
`, title)
	if q.auth {
		fmt.Fprintf(q.w, "<form action=\"logout\">%s &nbsp; <input type=\"submit\" value=\"Log uit\"></form>\n", html.EscapeString(q.user))
	} else {
		fmt.Fprintln(q.w, "<form action=\"login1\">E-mail: <input type=\"text\" name=\"mail\"> <input type=\"submit\" value=\"Log in\"></form>")
	}

	var t [5]string
	t[tab] = " class=\"selected\""
	fmt.Fprintln(q.w, "</div>\n<div id=\"topmenu\">\n<a href=\".\""+t[1]+">Zoeken</a>")
	if has_dbxml {
		fmt.Fprintln(q.w, "<a href=\"xpath\""+t[2]+">XPath</a>")
	}
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpora\""+t[3]+">Corpora</a>")
	}
	fmt.Fprintln(q.w, "<a href=\"info.html\""+t[4]+">Info</a>\n</div>\n")
}

func writeHtml(q *Context, title, msg string) {
	writeHead(q, title, 0)
	fmt.Fprintf(q.w, `
<h1>%s</h1>
%s
</body>
</html>
`, title, msg)
}

func first(r *http.Request, opt string) string {
	if len(r.Form[opt]) > 0 {
		return strings.TrimSpace(r.Form[opt][0])
	}
	return ""
}

func firstf(form *multipart.Form, opt string) string {
	if len(form.Value[opt]) > 0 {
		return strings.TrimSpace(form.Value[opt][0])
	}
	return ""
}

func urlencode(s string) string {
	var buf bytes.Buffer
	for _, b := range []byte(s) {
		if b >= 'a' && b <= 'z' ||
			b >= 'A' && b <= 'Z' ||
			b >= '0' && b <= '9' {
			buf.WriteByte(b)
		} else {
			buf.WriteString(fmt.Sprintf("%%%02x", b))
		}
	}
	return buf.String()
}

func iformat(i int) string {
	s1 := fmt.Sprint(i)
	s2 := ""
	for n := len(s1); n > 3; n = len(s1) {
		s2 = "." + s1[n-3:n] + s2
		s1 = s1[0 : n-3]
	}
	return s1 + s2
}

func shell(format string, a ...interface{}) *exec.Cmd {
	cmd := exec.Command(Cfg.Sh, "-c", fmt.Sprintf(format, a...))
	cmd.Env = []string{
		"PAQU=" + paqudir,
		"ALPINO_HOME=" + Cfg.Alpino,
		"PATH=" + Cfg.Path,
		"LANG=en_US.utf8",
		"LANGUAGE=en_US.utf8",
		"LC_ALL=en_US.utf8",
	}
	return cmd
}

func minversion(major, minor, patch int) bool {
	if version[0] > major {
		return true
	}
	if version[0] == major {
		if version[1] > minor {
			return true
		}
		if version[1] == minor {
			if version[2] >= patch {
				return true
			}
		}
	}
	return false
}

func dbopen() (*sql.DB, error) {
	return sql.Open("mysql", Cfg.Login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam")
}

func urlJoin(elem ...string) string {
	p := elem[0]
	for _, e := range elem[1:] {
		if strings.HasSuffix(p, "/") {
			if strings.HasPrefix(e, "/") {
				p += e[1:]
			} else {
				p += e
			}
		} else {
			if strings.HasPrefix(e, "/") {
				p += e
			} else {
				p += "/" + e
			}
		}
	}
	return p
}

func maxtitlelen(s string) string {
	if utf8.RuneCountInString(s) <= MAXTITLELEN {
		return s
	}
	r := make([]rune, MAXTITLELEN)
	p := 0
	for _, c := range s {
		r[p] = c
		p++
		if p == MAXTITLELEN {
			break
		}
	}
	return string(r)
}

func gz(filename string) error {
	fpin, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fpin.Close()
	fpout, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer fpout.Close()
	w := gzip.NewWriter(fpout)
	defer w.Close()
	_, err = io.Copy(w, fpin)
	if err != nil {
		return err
	}
	os.Remove(filename)
	return nil
}

func repl_filechar(s string) string {
	a := make([]string, 0, 5)
	for _, b := range []byte(s) {
		a = append(a, fmt.Sprintf("_%2X", b))
	}
	return strings.Join(a, "")
}

func encode_filename(s string) string {

	if s == "" {
		return "_"
	}

	s = strings.Replace(s, "_", "__", -1)

	s = reFilechars.ReplaceAllStringFunc(s, repl_filechar)

	if s[0] == '.' {
		s = "_2E" + s[1:]
	}
	if s[0] == '-' {
		s = "_2D" + s[1:]
	}
	return s
}

func repl_filecode(s string) string {
	if s == "__" {
		return "_"
	}
	i, _ := strconv.ParseInt(s[1:], 16, 0)
	b := []byte{byte(i)}
	return string(b)
}

func decode_filename(s string) string {
	if s == "_" {
		return ""
	}

	return reFilecodes.ReplaceAllStringFunc(s, repl_filecode)
}

// alle bestandsnamen van all subdirectories van de gegeven directory)
func filenames2(dirname string) ([]string, error) {
	fnames := make([]string, 0)
	dirs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fnames, err
	}
	for _, dir := range dirs {
		dname := dir.Name()
		files, err := ioutil.ReadDir(path.Join(dirname, dname))
		if err != nil {
			return fnames, err
		}
		for _, file := range files {
			fnames = append(fnames, path.Join(dname, file.Name()))
		}
	}
	return fnames, nil
}
