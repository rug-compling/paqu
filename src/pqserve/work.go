package main

import (
	"github.com/pebbe/util"

	"compress/gzip"
	"database/sql"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	LeegArchief = errors.New("Leeg zipbestand of tarbestand")
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
	isArch := false
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
		if strings.Contains(params, "-arch") {
			params = strings.Replace(params, "-arch", "", 1)
			isArch = true
		}
	}

	dirname := path.Join(paqudir, "data", task.id)
	data := path.Join(dirname, "data")
	xml := path.Join(dirname, "xml")
	dact := path.Join(dirname, "data.dact")
	stdout := path.Join(dirname, "stdout.txt")
	stderr := path.Join(dirname, "stderr.txt")
	summary := path.Join(dirname, "summary.txt")

	// gzip
	var fp *os.File
	fp, err = os.Open(data)
	if err != nil {
		return
	}
	b := make([]byte, 2)
	io.ReadFull(fp, b)
	fp.Close()
	if string(b) == "\x1F\x8B" {
		// gzip
		fpin, _ := os.Open(data)
		r, e := gzip.NewReader(fpin)
		if e != nil {
			fpin.Close()
			err = e
			return
		}
		fpout, _ := os.Create(data + ".tmp")
		_, err = io.Copy(fpout, r)
		fpout.Close()
		r.Close()
		fpin.Close()
		if err != nil {
			return
		}
		os.Rename(data+".tmp", data)
	}

	var ar *arch
	ar, err = NewArchReader(data)
	if err == nil {

		//    als eerste bestand alpino-xml is
		//       params is xmlzip
		//    anders uitpakken en alles aan elkaar plakken in data, met newline aan eind van elk deel

		e := ar.Next()
		if e == io.EOF {
			ar.Close()
			err = LeegArchief
			return
		}
		if e != nil {
			ar.Close()
			err = e
			return
		}

		var b []byte
		b, err = ar.ReadN(200)
		if err != nil {
			ar.Close()
			return
		}

		if strings.Contains(string(b), "<alpino_ds") {
			if !strings.HasPrefix(params, "xmlzip") {
				params = "xmlzip"
			}
			setinvoer(db, params, task.id, false)
			ar.Close()
		} else {
			isxml := strings.HasPrefix(string(b), "<?xml")
			if !isxml {
				isArch = true
			}
			ar.Close()
			var fp *os.File
			fp, err = os.Create(data + ".tmp")
			if err != nil {
				return
			}
			ar, err = NewArchReader(data)
			for {
				e := ar.Next()
				if e == io.EOF {
					break
				}
				if e != nil {
					ar.Close()
					fp.Close()
					err = e
					return
				}
				if !isxml {
					n := ar.Name()
					fmt.Fprintf(fp, "\n##PAQUFILE %s\n\n##PAQUMETA text paqu.filename = %s\n", n, n)
				}
				err = ar.Copy(fp)
				if err != nil {
					fp.Close()
					ar.Close()
					return
				}
				fmt.Fprintln(fp)
			}
			fp.Close()
			ar.Close()
			os.Rename(data+".tmp", data)
		}
	}

	if params == "auto" {
		params, err = invoersoort(db, data, task.id)
		if err != nil {
			return
		}
		if params == "dact" {
			os.Rename(data, dact)
		}
	}

	if isArch {
		setinvoer(db, params, task.id, true)
	}

	if params == "folia" {
		err = folia(data, data+".tmp")
		if err != nil {
			return
		}
		os.Rename(data+".tmp", data)
	}

	if params == "tei" {
		err = tei(data, data+".tmp")
		if err != nil {
			return
		}
		os.Rename(data+".tmp", data)
	}

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
		os.Remove(data + ".tmp")
		for _, f := range []string{data, data + ".lines", stdout, stderr, summary} {
			if f == data && (params == "dact" || strings.HasPrefix(params, "xmlzip")) {
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
		if strings.HasPrefix(params, "xmlzip") {
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
	} else if strings.HasPrefix(params, "xmlzip") {
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
			var has_tok, has_lbl bool
			if strings.Contains(params, "-lbl") || params == "folia" || params == "tei" {
				has_lbl = true
			}
			if strings.Contains(params, "-tok") || params == "folia" || params == "tei" {
				has_tok = true
			}

			//  pqtexter
			var pqtexter string
			if params == "run" {
				pqtexter = "-r"
			} else if has_lbl {
				pqtexter = "-l"
			}
			err = shell("pqtexter %s %s > %s.tmp 2>> %s", pqtexter, data, data, stderr).Run()
			if err != nil {
				return
			}
			os.Rename(data+".tmp", data)

			// tokenizer
			if !has_tok {
				var tok string
				if params == "run" {
					tok = "tokenize.sh"
				} else {
					tok = "tokenize_no_breaks.sh"
				}
				err = shell("$ALPINO_HOME/Tokenization/%s < %s > %s.tmp 2>> %s", tok, data, data, stderr).Run()
				if err != nil {
					return
				}
				os.Rename(data+".tmp", data)
			}

			var fp, fpin *os.File
			fpin, err = os.Open(data)
			if err != nil {
				return
			}
			fp, err = os.Create(data + ".lines")
			if err != nil {
				fpin.Close()
				return
			}

			rd := util.NewReader(fpin)
			var filename, lbl string
			var tokens, nlines, i int
			for {
				line, e := rd.ReadLineString()
				if e == io.EOF {
					break
				}
				if e != nil {
					err = e
					fp.Close()
					fpin.Close()
					return
				}
				if line == "" {
					continue
				}
				if strings.HasPrefix(line, "##PAQU") {
					a := strings.Fields(line)
					var val string
					if len(a) == 2 {
						b, e := hex.DecodeString(a[1])
						if e != nil {
							err = e
							fp.Close()
							fpin.Close()
							return
						}
						val = strings.TrimSpace(string(b))
					}
					if a[0] == "##PAQUFILE" {
						filename = val
						if strings.HasSuffix(filename, ".txt") || strings.HasSuffix(filename, ".doc") {
							filename = filename[:len(filename)-4]
						}
						i = 0
					} else if a[0] == "##PAQULBL" {
						lbl = val
					}
				} else if strings.HasPrefix(line, "##META") {
				} else {
					tokens += len(strings.Fields(line))
					if lbl == "" {
						if filename != "" {
							i++
							fmt.Fprintf(fp, "%04d/%04d-%s.%d|%s\n", nlines/10000, nlines%10000, encode_filename(filename), i, line)
						} else {
							fmt.Fprintf(fp, "%04d/%04d|%s\n", nlines/10000, nlines%10000, strings.TrimSpace(line))
						}
					} else {
						fmt.Fprintf(fp, "%04d/%04d-%s|%s\n", nlines/10000, nlines%10000, encode_filename(lbl), line)
						lbl = ""
					}
					if nlines%10000 == 0 {
						os.MkdirAll(path.Join(xml, fmt.Sprintf("%04d", nlines/10000)), 0777)
					}
					nlines++
				}
			}
			fp.Close()
			fpin.Close()

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

		//
		// TODO: ##META
		//

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
			a := strings.Split(line, "|")
			ln := len(a)
			n, err := strconv.Atoi(a[ln-3])
			if err == nil && n > 0 {
				continue
			}

			if ln > 5 {
				a[1] = strings.Join(a[1:ln-3], "|")
			}
			fname := a[0][2:]
			if params == "run" || strings.HasPrefix(params, "line") || params == "folia" || params == "tei" {
				fname = decode_filename(a[0][2:])
			}
			if strings.Contains(params, "-lbl") || params == "folia" || params == "tei" {
				fname = fname[1+strings.Index(fname, "-"):]
			}
			errlines = append(errlines, fname+"\t"+a[ln-3]+"\t"+a[ln-2]+"\t"+a[1]+"\n")
		}
	}
	fp, err = os.Create(summary)
	if err != nil {
		return
	}
	fmt.Fprintf(fp, "%d\t%d\t%s\n", nlines, len(errlines), invoertabel[params])
	for _, line := range errlines {
		fmt.Fprint(fp, line)
	}
	fp.Close()

	p := regexp.QuoteMeta(xml + "/")
	d := ""
	if strings.Contains(params, "-lbl") || params == "folia" || params == "tei" || isArch {
		p += "[0-9]+/[0-9]+-"
		d = "-d"
	} else if params == "dact" || strings.HasPrefix(params, "xmlzip") {
		p += "[0-9]+/"
		d = "-d"
	}

	cmd := shell(
		// optie -w i.v.m. recover()
		`find %s -name '*.xml' | sort | pqbuild -w -p %s %s -s %s %s %s 0 >> %s 2>> %s`,
		dirname,
		quote(p), d,
		path.Base(dirname), quote(title), quote(user), stdout, stderr)
	err = run(cmd, task.chKill, nil)
	if err != nil {
		return
	}

	if Cfg.Dact && params != "dact" {
		p := ""
		if strings.Contains(params, "-lbl") || params == "folia" || params == "tei" || isArch {
			p = "-"
		} else if strings.HasPrefix(params, "xmlzip") {
			p = "/"
		}
		err = makeDact(dact, xml, p, task.chKill)
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
		var params string
		rows, err := db.Query(fmt.Sprintf("SELECT `params` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, task.id))
		if err == nil {
			for rows.Next() {
				rows.Scan(&params)
			}
			rows.Close()
		}
		msg := ""
		switch params {
		case "xmlzip-d":
			msg = "afgeleid corpus"
		case "xmlzip-p":
			msg = "afgeleid corpus, beschermd"
		}
		db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"FINISHED\", `msg` = %q WHERE `id` = %q", Cfg.Prefix, msg, task.id))
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
	queuing := make([]string, 0)

	_, err = db.Exec(
		"UPDATE `" + Cfg.Prefix + "_info` SET `nword` = 0, `status` = \"QUEUED\" WHERE `status` = \"WORKING\"")
	util.CheckErr(err)

	rows, err := db.Query("SELECT `id`,`status` FROM `" +
		Cfg.Prefix + "_info` WHERE `status` = \"QUEUED\" OR `status` = \"QUEUING\"  ORDER BY `created`")
	util.CheckErr(err)
	for rows.Next() {
		var id, status string
		util.CheckErr(rows.Scan(&id, &status))
		if status == "QUEUED" {
			ids = append(ids, id)
		} else {
			queuing = append(queuing, id)
		}
	}
	util.CheckErr(rows.Err())

	if len(queuing) > 0 {
		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `status` = \"QUEUING\"", Cfg.Prefix))
		util.CheckErr(err)
	}
	for _, corpus := range queuing {
		util.CheckErr(os.RemoveAll(path.Join(paqudir, "data", corpus)))
		logf("QUEUING: rm -r %s: ok", path.Join(paqudir, "data", corpus))
	}

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

	ar, err := NewArchReader(data)
	if err != nil {
		return 0, 0, err
	}
	defer ar.Close()

	tokens = 0
	nline = 0
	nd := -1
	sdir := ""

	for {
		err := ar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, err
		}

		select {
		case <-chGlobalExit:
			return 0, 0, errGlobalExit
		case <-chKill:
			return 0, 0, errKilled
		default:
		}

		bindata, err := ar.Read()
		if len(bindata) == 0 {
			continue
		}

		nline++

		if nline%10000 == 1 {
			nd++
			sdir = fmt.Sprintf("%04d", nd)
			os.Mkdir(path.Join(xmldir, sdir), 0777)
		}

		name := ar.Name()
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
			return 0, 0, fmt.Errorf("Parsen van %q uit zip-bestand: %s", ar.Name(), err)
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
