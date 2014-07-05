package main

import (
	"github.com/pebbe/util"

	"archive/zip"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
)

func quote(s string) string {
	return "'" + strings.Replace(s, "'", "'\\''", -1) + "'"
}

func dowork(db *sql.DB, task *Process) (user string, title string, err error) {
	logf("WORKING: " + task.id)

	processLock.Lock()
	if task.nr > taskWorkNr {
		taskWorkNr = task.nr
		if taskWaitNr == taskWorkNr {
			// queue is leeg: reset counters om (ooit) overflow te voorkomen
			taskWaitNr = 0
			taskWorkNr = 0
		}
	}
	processLock.Unlock()

	_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"WORKING\", `nword` = 0 WHERE `id` = %q",
		Cfg.Prefix, task.id))
	if err != nil {
		return
	}

	params := "unknown"
	var rows *sql.Rows
	rows, err = db.Query(fmt.Sprintf("SELECT `description`,`owner`,`params` FROM `%s_info` WHERE `id` = %q",
		Cfg.Prefix, task.id))
	if err != nil {
		return
	}
	if rows.Next() {
		err = rows.Scan(&title, &user, &params)
		rows.Close()
		if err != nil {
			return
		}
	}

	dirname := path.Join(paqudir, "data", task.id)
	data := path.Join(dirname, "data")
	xml := path.Join(dirname, "xml")
	dact := path.Join(dirname, "data.dact")
	stdout := path.Join(dirname, "stdout.txt")
	stderr := path.Join(dirname, "stderr.txt")
	summary := path.Join(dirname, "summary.txt")

	defer func() {
		select {
		case <-chGlobalExit:
			err = errGlobalExit
			return
		case <-task.chKill:
			err = errKilled
			return
		default:
		}
		os.Remove(data + ".lines.tmp")
		for _, f := range []string{data, data + ".lines", dact, stdout, stderr, summary} {
			if f == data && (params == "dact" || params == "xmlzip") {
				continue
			}
			if f == dact && !Cfg.Dact {
				continue
			}
			logerr(gz(f))
		}
		fnames, e := filenames2(xml)
		if logerr(e) {
			return
		}
		for _, fname := range fnames {
			select {
			case <-chGlobalExit:
				err = errGlobalExit
				return
			case <-task.chKill:
				err = errKilled
				return
			default:
			}
			logerr(gz(path.Join(xml, fname)))
		}
		if params == "dact" && !Cfg.Dact {
			os.Remove(dact)
		}
		if params == "xmlzip" {
			os.Remove(data)
		}
	}()

	select {
	case <-chGlobalExit:
		err = errGlobalExit
		return
	case <-task.chKill:
		err = errKilled
		return
	default:
	}

	if params == "dact" {
		var tokens, nlines int
		tokens, nlines, err = unpackDact(data, xml, dact, stderr, task.chKill)
		if err != nil {
			return
		}
		err = do_quotum(db, task.id, user, tokens, nlines)
		if err != nil {
			os.Remove(dact)
			os.Remove(data + ".lines")
			os.Remove(stderr)
			os.RemoveAll(xml)
			return
		}
	} else if params == "xmlzip" {
		var tokens, nlines int
		tokens, nlines, err = unpackXml(data, xml, stderr, task.chKill)
		if err != nil {
			return
		}
		err = do_quotum(db, task.id, user, tokens, nlines)
		if err != nil {
			os.Remove(data + ".lines")
			os.Remove(stderr)
			os.RemoveAll(xml)
			return
		}
	} else { // if params != (dact || xmlzip)
		reuse := false
		reuse_more := false
		if files, e := filenames2(xml); e == nil && len(files) > 0 {
			reuse = true
			done := make(map[string]bool)
			for _, f := range files {
				b, e := ioutil.ReadFile(path.Join(xml, f))
				if e != nil {
					err = e
					return
				}
				if strings.Index(string(b), "</alpino_ds>") > 0 {
					done[f] = true
				}
			}
			fpin, e := os.Open(data + ".lines")
			if e != nil {
				err = e
				return
			}
			fpout, e := os.Create(data + ".lines.tmp")
			if e != nil {
				fpin.Close()
				err = e
				return
			}
			nword := 0
			r := util.NewReader(fpin)
			for {
				line, e := r.ReadLineString()
				if e != nil {
					fpout.Close()
					fpin.Close()
					if e != io.EOF {
						err = e
						return
					}
					break
				}
				nword += len(strings.Fields(line))
				key := line[:strings.Index(line, "|")]
				if !done[key+".xml"] {
					fmt.Fprintln(fpout, line)
					reuse_more = true
				}
			}

			_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `nword` = %d WHERE `id` = %q",
				Cfg.Prefix, nword, task.id))
			if err != nil {
				return
			}

		} else { // if !reuse
			var pqtexter, tok string
			switch params {
			case "run":
				tok = "tokenize.sh"
				pqtexter = "-r"
			case "line":
				tok = "tokenize_no_breaks.sh"
			default:
				err = errors.New("Niet geïmplementeerd: " + params)
				return
			}
			cmd := shell("pqtexter %s %s | $ALPINO_HOME/Tokenization/%s 2>> %s", pqtexter, data, tok, stderr)
			chPipe := make(chan string)
			chTokens := make(chan int, 2)
			var fp *os.File
			fp, err = os.Create(data + ".lines")
			if err != nil {
				return
			}
			go func() {
				var tokens, lineno int
				for line := range chPipe {
					tokens += len(strings.Fields(line))
					fmt.Fprintf(fp, "%04d/%04d|%s\n", lineno/10000, lineno%10000, line)
					if lineno%10000 == 0 {
						os.MkdirAll(path.Join(xml, fmt.Sprintf("%04d", lineno/10000)), 0777)
					}
					lineno++
				}
				fp.Close()
				chTokens <- tokens
				chTokens <- lineno
			}()
			err = run(cmd, task.chKill, chPipe)
			if err != nil {
				return
			}
			tokens := <-chTokens
			nlines := <-chTokens

			err = do_quotum(db, task.id, user, tokens, nlines)
			if err != nil {
				os.Remove(data)
				os.Remove(data + ".lines")
				return
			}

			if err != nil {
				return
			}

		} // end if !reuse

		if !reuse || reuse_more {

			var ext string
			if reuse {
				ext = ".tmp"
			}

			var timeout string
			if Cfg.Timeout > 0 {
				timeout = fmt.Sprint("-t ", Cfg.Timeout)
			}
			cmd := shell(
				`pqalpino -a %s -d %s %s %s.lines%s >> %s 2>> %s`,
				Cfg.Alpino, xml, timeout, data, ext, stdout, stderr)
			err = run(cmd, task.chKill, nil)
			if err != nil {
				return
			}
		}
	} // end if params != (dact || xmlzip)

	// TODO: inlezen uit xml-bestanden als er een Alpino-server wordt gebruikt
	nlines := 0
	errlines := make([]string, 0)
	fp, e := os.Open(stderr)
	if e != nil {
		err = e
		return
	}
	rd := util.NewReader(fp)
	for {
		line, e := rd.ReadLineString()
		if e != nil {
			fp.Close()
			if e == io.EOF {
				break
			} else {
				err = e
				return
			}
		}
		if strings.HasPrefix(line, "Q#") {
			nlines++
			if strings.Index(line, "|1|1|") < 0 {
				errlines = append(errlines, line)
			}
		}
	}

	fp, err = os.Create(summary)
	if err != nil {
		return
	}

	s := params
	switch params {
	case "run":
		s = "doorlopende tekst"
	case "line":
		s = "een zin per regel"
	case "dact":
		s = "Dact-bestand"
	case "xmlzip":
		s = "Alpino XML-bestanden in zipbestand"
	}
	fmt.Fprintf(fp, "Bron: %s\n\n", s)

	if nlines > 0 {
		if len(errlines) == 0 {
			fmt.Fprintf(fp, "Alle %d regels zijn met succes geparst.\n", nlines)
		} else {
			fmt.Fprintf(
				fp,
				`%d van de %d regels zijn met succes geparst.

Er waren problemen met de %d regels hieronder. Misschien was er bij
die regels een time-out waardoor geen volledige parse gedaan kon worden.

`,
				nlines-len(errlines),
				nlines,
				len(errlines))
			for _, line := range errlines {
				a := strings.Split(line, "|")
				if len(a) > 5 {
					a[1] = strings.Join(a[1:len(a)-3], "|")
				}
				fmt.Fprintf(fp, "%s\t%s\n", decode_filename(a[0][2:]), a[1])
			}
		}
	}
	fp.Close()

	cmd := shell(
		// optie -w i.v.m. recover()
		`find %s -name '*.xml' | pqbuild -w -s %s %s %s 0 >> %s 2>> %s`,
		dirname,
		path.Base(dirname), quote(title), quote(user), stdout, stderr)
	err = run(cmd, task.chKill, nil)
	if err != nil {
		return
	}

	if Cfg.Dact && params != "dact" {
		err = makeDact(dact, xml, params == "xmlzip", task.chKill)
		if err != nil {
			return
		}
	}

	return
}

// Run command, maar onderbreek het als chKill of chGlobalExit gesloten is
func run(cmd *exec.Cmd, chKill chan bool, chPipe chan string) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	chRet := make(chan error, 2)
	// deze functie schrijft twee keer op chRet
	go func() {
		if chPipe == nil {
			cmd.Start()
			chRet <- nil
			err := cmd.Wait()
			chRet <- err
			return
		}

		pipe, err := cmd.StdoutPipe()
		if err != nil {
			chRet <- nil
			chRet <- err
			return
		}
		cmd.Start()
		chRet <- nil

		var err1, err2 error

		rd := util.NewReader(pipe)
		for {
			var line string
			line, err1 = rd.ReadLineString()
			if err1 == io.EOF {
				err1 = nil
				break
			}
			if err1 != nil {
				break
			}
			chPipe <- line
		}
		close(chPipe)

		err2 = cmd.Wait()

		if err1 != nil {
			chRet <- err1
		} else {
			chRet <- err2
		}
	}()

	<-chRet // commando is gestart

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		// misschien betekent een fout alleen maar dat het process net klaar is
		logf("BIG TROUBLE: syscall.Getpgid(cmd.Process.Pid) error: %v", err)
		pgid = 0
	}

FORSELECT:
	for {
		select {
		case err = <-chRet:
			return err
		case <-chGlobalExit:
			break FORSELECT
		case <-chKill:
			break FORSELECT
		}
	}
	err = syscall.Kill(-pgid, 9)
	if err != nil {
		logf("syscall.Kill(-pgid, 9) error: %v", err)
	}
	err = <-chRet
	return err
}

func work(task *Process) {

	// taak niet uitvoeren als ie al gekilld is
	done := false
	task.lock.Lock()
	task.queued = false
	if task.killed {
		done = true
	}
	task.lock.Unlock()
	if done {
		return
	}

	defer func() {
		processLock.Lock()
		delete(processes, task.id)
		processLock.Unlock()
	}()

	db, err := dbopen()
	if err != nil {
		logerr(err)
		return
	}
	defer db.Close()
	user, title, err := dowork(db, task)

	select {
	case <-chGlobalExit:
		return
	default:
	}

	if err == nil {
		logf("FINISHED: " + task.id)
		sendmail(user, "Corpus klaar", fmt.Sprintf("Je corpus \"%s\" staat klaar op %s", title, urlJoin(Cfg.Url, "/?db="+task.id)))
		db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"FINISHED\" WHERE `id` = %q", Cfg.Prefix, task.id))
	} else {
		logf("FAILED: %v, %v", task.id, err)
		if !task.killed {
			sendmail(user, "Corpus fout", fmt.Sprintf("Er ging iets fout met je corpus \"%s\": %v", title, err))
		}
		db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"FAILED\", `msg` = %q WHERE `id` = %q", Cfg.Prefix, err.Error(), task.id))
	}
}

func kill(id string) {
	processLock.RLock()
	task, ok := processes[id]
	processLock.RUnlock()
	if ok {
		done := false

		task.lock.Lock()
		if task.killed {
			task.lock.Unlock()
			return
		}
		task.killed = true
		if task.queued {
			done = true
		}
		task.lock.Unlock()

		if done {
			logf("UNQUEUED: %v", id)
			processLock.Lock()
			delete(processes, task.id)
			processLock.Unlock()
			return
		}

		close(task.chKill)
		for {
			time.Sleep(500 * time.Millisecond)
			processLock.RLock()
			_, ok := processes[id]
			processLock.RUnlock()
			if !ok {
				logf("KILLED: %v", id)
				return
			}
		}
	}
}

func recover() {
	db, err := dbopen()
	util.CheckErr(err)

	ids := make([]string, 0)

	_, err = db.Exec(
		"UPDATE `" + Cfg.Prefix + "_info` SET `nword` = 0, `status` = \"QUEUED\" WHERE `status` = \"WORKING\"")
	util.CheckErr(err)

	rows, err := db.Query("SELECT `id` FROM `" + Cfg.Prefix + "_info` WHERE `status` = \"QUEUED\" ORDER BY `created`")
	util.CheckErr(err)
	for rows.Next() {
		var id string
		util.CheckErr(rows.Scan(&id))
		ids = append(ids, id)
	}
	util.CheckErr(rows.Err())

	db.Close()

	for _, id := range ids {
		p := &Process{
			id:     id,
			chKill: make(chan bool),
			queued: true,
		}
		processLock.Lock()
		taskWaitNr++
		p.nr = taskWaitNr
		processes[id] = p
		processLock.Unlock()
		chWork <- p
	}
}

func do_quotum(db *sql.DB, id, user string, tokens, nlines int) error {
	quotumLock.Lock()
	defer quotumLock.Unlock()

	quotum := 0
	gebruikt := 0
	rows, err := db.Query(fmt.Sprintf("SELECT `quotum` FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, user))
	if err != nil {
		return err
	}
	if rows.Next() {
		err = rows.Scan(&quotum)
		rows.Close()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("MySQL: Kan quotum niet vinden")
	}
	if quotum > 0 {
		rows, err = db.Query(fmt.Sprintf("SELECT `nword` FROM `%s_info` WHERE `owner` = %q", Cfg.Prefix, user))
		if err != nil {
			return err
		}
		for rows.Next() {
			var n int
			err = rows.Scan(&n)
			if err != nil {
				return err
			}
			gebruikt += n
		}
	}

	if quotum > 0 && quotum-gebruikt < tokens {
		return fmt.Errorf("Ruimte voor %d tokens. Nieuw corpus bevat %d tokens.", quotum-gebruikt, tokens)
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `nword` = %d, `nline` = %d WHERE `id` = %q",
		Cfg.Prefix, tokens, nlines, id))
	return err
}

func unpackXml(data, xmldir, stderr string, chKill chan bool) (tokens, nline int, err error) {

	os.Mkdir(xmldir, 0777)

	fperr, err := os.Create(stderr)
	if err != nil {
		return 0, 0, err
	}
	defer fperr.Close()

	fplines, err := os.Create(data + ".lines")
	if err != nil {
		return 0, 0, err
	}
	defer fplines.Close()

	z, err := zip.OpenReader(data)
	if err != nil {
		return 0, 0, err
	}
	defer z.Close()

	tokens = 0
	nline = 0
	nd := -1
	sdir := ""

	for _, f := range z.File {

		select {
		case <-chGlobalExit:
			return 0, 0, errGlobalExit
		case <-chKill:
			return 0, 0, errKilled
		default:
		}

		rc, err := f.Open()
		if err != nil {
			return 0, 0, err
		}
		bindata, err := ioutil.ReadAll(rc)
		rc.Close()
		if len(bindata) == 0 {
			continue
		}

		nline++

		if nline%10000 == 1 {
			nd++
			sdir = fmt.Sprintf("%04d", nd)
			os.Mkdir(path.Join(xmldir, sdir), 0777)
		}

		name := f.Name
		if strings.HasSuffix(strings.ToLower(name), ".xml") {
			name = name[:len(name)-4]
		}
		encname := encode_filename(name)

		fp, err := os.Create(path.Join(xmldir, sdir, encname+".xml"))
		if err != nil {
			return 0, 0, err
		}
		_, err = fp.Write(bindata)
		fp.Close()
		if err != nil {
			return 0, 0, err
		}

		alpino := Alpino_ds_no_node{}
		err = xml.Unmarshal(bindata, &alpino)
		if err != nil {
			return 0, 0, fmt.Errorf("Parsen van %q uit zip-bestand: %s", f.Name, err)
		}
		tokens += len(strings.Fields(alpino.Sentence))
		fmt.Fprintf(fplines, "%s|%s\n", name, strings.TrimSpace(alpino.Sentence))
		for _, c := range alpino.Comments {
			if strings.HasPrefix(c.Comment, "Q#") {
				a := strings.SplitN(c.Comment, "|", 2)
				if len(a) == 2 {
					fmt.Fprintf(fperr, "Q#%s|%s\n", name, strings.TrimSpace(a[1]))
					break
				}
			}
		}
	}

	return tokens, nline, nil
}
