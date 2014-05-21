package main

//. Imports

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

//. Main

func stats(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	err := r.ParseForm()
	if err != nil {
		writeErrHtml(false, w, err)
		return
	}

	download := false

	if first(r, "d") != "" {
		download = true
	}

	if !download {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
	}

	db, err := connect()
	if err != nil {
		writeErrHtml(!download, w, err)
		return
	}
	defer db.Close()

	option := make(map[string]string)
	for _, t := range []string{"word", "postag", "rel", "hpostag", "hword"} {
		option[t] = first(r, t)
	}
	if option["word"] == "" && option["postag"] == "" && option["rel"] == "" && option["hpostag"] == "" && option["hword"] == "" {
		writeErrHtml(!download, w, errors.New("Missing query"))
		return
	}

	prefix := first(r, "db")
	if prefix == "" && len(opt_db) > 0 {
		prefix = strings.Fields(opt_db[0])[0]
	}
	if prefix == "" {
		writeErrHtml(!download, w, errors.New("Missing database in query"))
		return
	}

	if !prefixes[prefix] {
		writeErrHtml(!download, w, fmt.Errorf("Invalid data set: %s", prefix))
		return
	}

	// BEGIN: Opstellen van de query
	// Deze code moet gelijk zijn aan die in het programma 'lassystats'
	parts := make([]string, 0, 6)
	for _, p := range []string{"", "h"} {
		if option[p+"word"] != "" {
			wrd := option[p+"word"]
			if wrd[0] == '+' {
				parts = append(parts, fmt.Sprintf("`"+p+"lemma` = %q", wrd[1:]))
			} else if wrd[0] == '@' {
				parts = append(parts, fmt.Sprintf("`"+p+"root` = %q", wrd[1:]))
			} else if wrd[0] == '=' {
				parts = append(parts, fmt.Sprintf("`"+p+"word` = %q COLLATE \"utf8_bin\"", wrd[1:]))
			} else if wrd[0] == '?' {
				parts = append(parts, fmt.Sprintf("`"+p+"word` = %q", wrd[1:]))
			} else if strings.Index(wrd, "%") >= 0 {
				parts = append(parts, fmt.Sprintf("`"+p+"word` LIKE %q", wrd))
			} else {
				var s string
				rows, err := db.Query(fmt.Sprintf("SELECT `lemma` FROM `%s_c_%s_word` WHERE `word` = %q", PRE, prefix, wrd))
				if err != nil {
					writeErrHtml(!download, w, err)
					return
				}
				lset := make(map[string]bool)
				for rows.Next() {
					err := rows.Scan(&s)
					if err != nil {
						writeErrHtml(!download, w, err)
						return
					}
					for _, i := range strings.Split(s, "\t") {
						lset[fmt.Sprintf("%q", i)] = true
					}
				}
				err = rows.Err()
				if err != nil {
					writeErrHtml(!download, w, err)
					return
				}
				if len(lset) > 0 {
					ll := make([]string, 0, len(lset))
					for key := range lset {
						ll = append(ll, key)
					}
					if len(ll) == 1 {
						parts = append(parts, "`"+p+"lemma` = "+ll[0])
					} else {
						parts = append(parts, "`"+p+"lemma` IN ("+strings.Join(ll, ", ")+")")
					}
				} else {
					parts = append(parts, fmt.Sprintf("`"+p+"word` = %q", wrd))
				}
			}
		}
	}
	if option["postag"] != "" {
		parts = append(parts, fmt.Sprintf("`postag` = %q", option["postag"]))
	}
	if option["rel"] != "" {
		parts = append(parts, fmt.Sprintf("`rel` = %q", option["rel"]))
	}
	if option["hpostag"] != "" {
		parts = append(parts, fmt.Sprintf("`hpostag` = %q", option["hpostag"]))
	}
	query := strings.Join(parts, " AND ")
	// EINDE: Opstellen van de query

	// DEBUG: HTML-uitvoer van de query
	if !download {
		fmt.Fprint(w, "<div style=\"font-family:monospace\">\n", html.EscapeString(query), "\n</div><p>\n")
	}

	// Aantal zinnen die matchen met de query
	rows, err := db.Query("SELECT 1 FROM `" + PRE + "_c_" + prefix + "_deprel` WHERE " + query + " GROUP BY `arch`,`file`")
	if err != nil {
		writeErrHtml(!download, w, err)
		return
	}
	counter := 0
	for rows.Next() {
		counter++
	}
	if err != nil {
		writeErrHtml(!download, w, err)
		return
	}
	if download {
		fmt.Fprintf(w, "# %d zinnen\t\n", counter)
	} else {
		fmt.Fprintln(w, "Aantal gevonden zinnen:", counter)
	}

	// Tellingen van onderdelen
	for i, ww := range []string{"word", "lemma", "postag", "rel", "hword", "hlemma", "hpostag"} {
		var j, count int
		var s, p, limit string
		if download {
			fmt.Fprintln(w, "# "+ww+"\t")
		} else {
			if i == 0 {
				fmt.Fprintln(w, "<p>"+YELLOW+"<b>word</b></span>: ")
			} else if i == 4 {
				fmt.Fprintln(w, "<p>"+GREEN+"<b>hword</b></span>: ")
			} else {
				fmt.Fprintln(w, "<p><b>"+ww+"</b>: ")
			}
			limit = " LIMIT " + fmt.Sprint(WRDMAX)
		}
		rows, err := db.Query("SELECT count(*), `" + ww + "` FROM `" + PRE + "_c_" + prefix + "_deprel` WHERE " + query +
			" GROUP BY `" + ww + "` COLLATE 'utf8_bin' ORDER BY 1 DESC, 2" + limit)
		if err != nil {
			writeErrHtml(!download, w, err)
			return
		}
		for rows.Next() {
			err := rows.Scan(&j, &s)
			if err != nil {
				writeErrHtml(!download, w, err)
				return
			}
			if s == "" {
				s = "\"\""
			}
			if download {
				fmt.Fprintf(w, "%d\t%s\n", j, s)
			} else {
				fmt.Fprint(w, p, j, "&times;&nbsp;", html.EscapeString(s))
				p = ", "
				count++
			}
		}
		err = rows.Err()
		if err != nil {
			writeErrHtml(!download, w, err)
			return
		}
		if count == WRDMAX {
			fmt.Fprint(w, ", ...")
		}
		if !download {
			fmt.Fprint(w, "\n<BR>\n")
		}
	}

	db.Close()

	if !download {
		fmt.Fprintf(w,
			"<hr>tijd: %s\n<p>\n<a href=\"/stats?%s&amp;d=1\" target=\"_blank\">download</a>\n",
			time.Now().Sub(now),
			r.URL.RawQuery)
	}
}
