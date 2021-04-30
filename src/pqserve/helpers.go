package main

import (
	pqnode "github.com/rug-compling/paqu/internal/node"

	"github.com/BurntSushi/toml"

	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	reFilechars = regexp.MustCompile("[^-._a-zA-Z0-9]+")
	reFilecodes = regexp.MustCompile("_[0-9A-F][0-9A-F]|__")
)

// fout zonder gebruiker
func sysErr(err error) bool {
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

	return true
}

// system error -> log
func doErr(q *Context, err error) bool {
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

	fmt.Fprintln(q.w, "<pre>")
	fmt.Fprintln(q.w, html.EscapeString(s))
	fmt.Fprintln(q.w, "</pre>\n</body>\n</html>")

	return true
}

// user error -> geen log
func userErr(q *Context, err error) bool {
	if err == nil {
		return false
	}

	fmt.Fprintln(q.w, "<pre>")
	fmt.Fprintln(q.w, html.EscapeString(err.Error()))
	fmt.Fprintln(q.w, "</pre>\n</body>\n</html>")

	return true
}

// system error -> log
func hErr(q *Context, err error) bool {
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

	http.Error(q.w, s, http.StatusInternalServerError)

	return true
}

// user error -> geen log
func uhErr(q *Context, err error) bool {
	if err == nil {
		return false
	}

	http.Error(q.w, err.Error(), http.StatusPreconditionFailed)

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
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="noindex,nofollow">
<link rel="icon" href="favicon.ico" type="image/x-icon">
<link rel="stylesheet" type="text/css" href="paqu.css">
<!--[if gte IE 10]> -->
<style type="text/css">span.ie { display:none; }</style>
<!-- <![endif]-->
</head>
<body>
<div id="login">
`, title)
	if q.auth {
		fmt.Fprintf(q.w, "<form action=\"logout\">%s &nbsp; <input type=\"submit\" value=\"Log uit\"></form>\n", html.EscapeString(q.user[:strings.Index(q.user, "@")]))
	} else {
		if u := strings.TrimSpace(Cfg.Loginurl); u == "" {
			fmt.Fprintln(q.w, "<form action=\"login1\"><span class=\"ie\">E-mail: </span><input type=\"email\" name=\"mail\" placeholder=\"E-mail\"> <input type=\"submit\" value=\"Log in\"></form>")
		} else {
			fmt.Fprintf(q.w, "<form action=\"%s\"><input type=\"submit\" value=\"Log in\"></form>\n", u)
		}
	}

	n := 7 + len(localMenu)
	t := make([]string, n)
	if tab < n {
		t[tab] = " class=\"selected\""
	}
	fmt.Fprintln(q.w, "</div>\n<div id=\"topmenu\">\n<a href=\".\""+t[1]+">Zoeken</a>")
	if has_dbxml {
		fmt.Fprintln(q.w, "<a href=\"xpath\""+t[2]+">XPath</a>")
	}
	if len(q.hasmeta) > 0 {
		fmt.Fprintln(q.w, "<a href=\"metadata\""+t[3]+">Metadata</a>")
	}
	if q.auth {
		fmt.Fprintln(q.w, "<a href=\"corpora\""+t[4]+">Corpora</a>")
	}
	fmt.Fprintln(q.w, "<a href=\"spod\""+t[5]+">SPOD</a>")
	fmt.Fprintln(q.w, "<a href=\"info.html\""+t[6]+">Info</a>")
	for i, item := range localMenu {
		if q.auth || !item.needAuth {
			fmt.Fprintf(q.w, "<a href=%q %s>%s</a>\n", item.path, t[7+i], item.text)
		}
	}
	fmt.Fprintln(q.w, "</div>\n")
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
	if form == nil {
		return ""
	}
	if len(form.File[opt]) > 0 {
		file, err := form.File[opt][0].Open()
		if logerr(err) {
			return ""
		}
		data, err := ioutil.ReadAll(file)
		file.Close()
		if logerr(err) {
			return ""
		}
		return string(data)
	}

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
		// U+202F = NARROW NO-BREAK SPACE
		s2 = "\u202F" + s1[n-3:n] + s2
		s1 = s1[0 : n-3]
	}
	return s1 + s2
}

func iformat0(i int) string {
	if i == 0 {
		return ""
	}
	return iformat(i)
}

func shell(format string, a ...interface{}) *exec.Cmd {
	cmd := exec.Command(Cfg.Sh, "-c", fmt.Sprintf(format, a...))
	cmd.Env = []string{
		"ALPINO_HOME=" + Cfg.Alpino,
		"PATH=" + Cfg.Path,
		"LANG=en_US.utf8",
		"LANGUAGE=en_US.utf8",
		"LC_ALL=en_US.utf8",
	}
	for _, e := range []string{"HOME", "PAQU", "XDG_DATA_HOME", "XDG_CONFIG_HOME", "LD_LIBRARY_PATH"} {
		if p := os.Getenv(e); p != "" {
			cmd.Env = append(cmd.Env, e+"="+p)
		}
	}
	if Cfg.Login[0] == '$' {
		cmd.Env = append(cmd.Env, Cfg.Login[1:]+"="+os.Getenv(Cfg.Login[1:]))
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
	login := Cfg.Login
	if login[0] == '$' {
		login = os.Getenv(login[1:])
	}
	return sql.Open("mysql", login+"?charset=utf8&parseTime=true&loc=Europe%2FAmsterdam&sql_mode=''")
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
	fpout, err := os.Create(filename + ".gz")
	if err != nil {
		fpin.Close()
		return err
	}
	defer fpout.Close()
	w := gzip.NewWriter(fpout)
	defer w.Close()
	_, err = io.Copy(w, fpin)
	if err != nil {
		fpin.Close()
		return err
	}
	fpin.Close()
	return os.Remove(filename)
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
func filenames2(dirname string, meta bool) ([]string, error) {
	fnames := make([]string, 0)
	dirs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fnames, err
	}
	for _, dir := range dirs {
		dname := dir.Name()
		files, err := ioutil.ReadDir(filepath.Join(dirname, dname))
		if err != nil {
			return fnames, err
		}
		for _, file := range files {
			if name := file.Name(); strings.HasSuffix(name, ".meta") == meta {
				fnames = append(fnames, filepath.Join(dname, name))
			}
		}
	}
	return fnames, nil
}

var reTijd = regexp.MustCompile(`\.[0-9]+([a-z]+)$`)

func tijd(t time.Duration) string {
	return reTijd.ReplaceAllString(t.String(), "$1")
}

func getMeta(q *Context, prefix string) []MetaType {
	result := make([]MetaType, 0)
	rows, err := sqlDB.Query(fmt.Sprintf(
		"SELECT `id`,`name`,`type`,`indexed`,`dtype`,`istep` FROM `%s_c_%s_midx` LEFT JOIN `%s_c_%s_minf` USING (`id`) ORDER BY 2",
		Cfg.Prefix, prefix,
		Cfg.Prefix, prefix))
	if logerr(err) {
		return result
	}
	var i int
	var n, t string
	var indexed sql.NullBool
	var dtype, istep sql.NullInt64
	for rows.Next() {
		if logerr(rows.Scan(&i, &n, &t, &indexed, &dtype, &istep)) {
			continue
		}
		v := "interval"
		switch t {
		case "TEXT":
			v = "waarde"
		case "INT":
			if !indexed.Bool || istep.Int64 < 2 {
				v = "waarde"
			}
		case "FLOAT":
			if !indexed.Bool {
				v = "waarde"
			}
		case "DATE", "DATETIME":
			switch int(dtype.Int64) {
			case dr_hour:
				v = "datum + uur"
			case dr_day:
				v = "datum"
			case dr_month:
				v = "maand"
			case dr_year:
				v = "jaar"
			case dr_dec:
				v = "decennium"
			case dr_cent:
				v = "eeuw"
			}
		}
		result = append(result, MetaType{i, unHigh(n), t, v})
	}
	logerr(rows.Err())
	return result
}

func TomlDecodeFile(fpath string, v interface{}) (toml.MetaData, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return toml.MetaData{}, err
	}
	// skip BOM (berucht op Windows)
	if bytes.HasPrefix(bs, []byte{239, 187, 191}) {
		bs = bs[3:]
	}
	return toml.Decode(string(bs), v)
}

func format(alpino Alpino_ds_complete) (string, error) {

	minifyNode(alpino.Node0)

	b, err := xml.MarshalIndent(&alpino, "", "  ")
	if err != nil {
		return "", nil
	}
	s := "<?xml version=\"1.0\"?>\n" + string(b)

	// shorten
	for _, v := range []string{"meta", "parser", "node", "dep"} {
		s = strings.Replace(s, "></"+v+">", "/>", -1)
	}

	return s, nil
}

func minifyNode(node *pqnode.Node) {
	if node == nil {
		return
	}
	if node.Ud != nil {
		if node.Ud.Id == "" {
			node.Ud = nil
		} else {
			if len(node.Ud.Dep) == 0 {
				node.Ud.Dep = nil
			}
		}
	}
	if node.NodeList != nil {
		for _, n := range node.NodeList {
			minifyNode(n)
		}
	}
}

func ifelse(v bool, yes interface{}, no interface{}) interface{} {
	if v {
		return yes
	}
	return no
}
