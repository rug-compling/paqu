package main

//. Imports

import (
	"github.com/pebbe/util"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

//. Main

func main() {
	if len(os.Args) != 3 {
		fmt.Printf(`
Syntax: %s login_file port_number
`, os.Args[0])
		return
	}

	data, err := ioutil.ReadFile(os.Args[1])
	util.CheckErr(err)
	login = strings.TrimSpace(string(data))

	port := os.Args[2]

	db, err := connect()
	util.CheckErr(err)
	for _, tag := range []string{"postag", "hpostag", "rel"} {
		rows, err := db.Query("SELECT tag FROM " + PRE + "_" + tag + ";")
		util.CheckErr(err)
		for rows.Next() {
			var val string
			util.CheckErr(rows.Scan(&val))
			switch tag {
			case "postag":
				opt_postag = append(opt_postag, val)
			case "hpostag":
				opt_hpostag = append(opt_hpostag, val)
			case "rel":
				opt_rel = append(opt_rel, val)
			}
		}
		util.CheckErr(rows.Err())
	}
	rows, err := db.Query("SELECT id, description FROM " + PRE + "_info")
	util.CheckErr(err)
	for rows.Next() {
		var id, desc string
		util.CheckErr(rows.Scan(&id, &desc))
		opt_dbc = append(opt_dbc, &Corpus{Id: id, Title: desc})
		prefixes[id] = true
	}
	util.CheckErr(rows.Err())
	sort.Sort(Corpora(opt_dbc))
	for _, dbc := range opt_dbc {
		opt_db = append(opt_db, fmt.Sprintf("%s %s", dbc.Id, dbc.Title))
	}

	db.Close()

	for _, c := range opt_db {
		fmt.Printf("Corpus %s\n", strings.Replace(c, " ", ": ", 1))
	}

	handleFunc("info.html", static_info_html)
	handleFunc("", form)
	handleFunc("tree", tree)
	handleFunc("stats", stats)
	handleFunc("clarinnl.png", static_clarinnl_png)
	handleFunc("jquery.js", static_jquery_js)
	handleFunc("wordrel.css", static_wordrel_css)
	handleFunc("tooltip.css", static_tooltip_css)
	handleFunc("tooltip.js", static_tooltip_js)

	fmt.Print("Serving on http://127.0.0.1:", port, "/\n")

	util.CheckErr(http.ListenAndServe(":"+port, Log(http.DefaultServeMux)))
}

func handleFunc(url string, handler func(http.ResponseWriter, *http.Request)) {
	url = path.Join("/", url)
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url {
				http.NotFound(w, r)
			} else {
				handler(w, r)
			}
		})
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
