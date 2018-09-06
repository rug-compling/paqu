package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func remove(q *Context) {

	if !q.auth {
		http.Error(q.w, "Je bent niet ingelogd", http.StatusUnauthorized)
		return
	}

	id := first(q.r, "id")

	// Kan myprefixes niet gebruiken omdat daar alleen corpora in staan die al klaar zijn

	rows, err := q.db.Query(fmt.Sprintf(
		"SELECT `description` FROM `%s_info` WHERE `id` = %q AND `owner` = %q",
		Cfg.Prefix, id, q.user))
	if err != nil {
		http.Error(q.w, err.Error(), http.StatusInternalServerError)
		logerr(err)
		return
	}

	my := false
	var desc string
	if rows.Next() {
		err := rows.Scan(&desc)
		if err == nil {
			rows.Close()
			my = true
		} else {
			http.Error(q.w, err.Error(), http.StatusInternalServerError)
			logerr(err)
			return
		}
	}

	if !my {
		http.Error(q.w, "Dat is niet je corpus", http.StatusUnauthorized)
		return
	}

	kill(id)

	go func() {
		p := filepath.Join(paqudatadir, "data", id)
		p2 := filepath.Join(paqudatadir, "data", "_invalid_"+id+"_"+rand16())
		err := os.Rename(p, p2)
		if err != nil {
			logf("os.Rename(%q, %q) error: %v", p, p2, err)
		}
		for i := 0; i < 12; i++ {
			err = os.RemoveAll(p2)
			if err == nil {
				break
			}
			time.Sleep(10 * time.Second)
		}
		if err != nil {
			logf("os.RemoveAll(%q) error: %v", p2, err)
		} else {
			logf("os.RemoveAll(%q) OK", p2)
		}
	}()

	_, err = q.db.Exec(fmt.Sprintf(
		"DROP TABLE IF EXISTS `%s_c_%s_deprel`, `%s_c_%s_sent`, `%s_c_%s_file`, `%s_c_%s_arch` , `%s_c_%s_word`, "+
			"`%s_c_%s_meta`, `%s_c_%s_midx`, `%s_c_%s_minf` , `%s_c_%s_mval`",
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id,
		Cfg.Prefix, id))
	logerr(err)
	_, err = q.db.Exec(fmt.Sprintf("DELETE FROM `%s_corpora` WHERE `prefix` = %q", Cfg.Prefix, id))
	logerr(err)
	_, err = q.db.Exec(fmt.Sprintf("DELETE FROM `%s_ignore` WHERE `prefix` = %q", Cfg.Prefix, id))
	logerr(err)
	_, err = q.db.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, id))
	logerr(err)

	logf("DELETED: %v", id)

	http.Redirect(q.w, q.r, urlJoin(Cfg.Url, "corpora"), http.StatusTemporaryRedirect)
}
