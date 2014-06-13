package main

import (
	"github.com/pebbe/util"

	"database/sql"
	"errors"
	"fmt"
	"io"
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
	stdout := path.Join(dirname, "stdout.txt")
	stderr := path.Join(dirname, "stderr.txt")

	select {
	case <-task.chKill:
		err = fmt.Errorf("Killed")
		return
	default:
	}

	var prepare, tok string
	switch params {
	case "run":
		tok = "tokenize.sh"
		prepare = "-r"
	case "line":
		tok = "tokenize_no_breaks.sh"
	default:
		err = errors.New("Not implemented: " + params)
		return
	}
	cmd := shell("prepare %s %s | $ALPINO_HOME/Tokenization/%s 2> %s", prepare, data, tok, stderr)
	chPipe := make(chan string)
	chTokens := make(chan int, 1)
	var fp *os.File
	fp, err = os.Create(data + ".lines")
	if err != nil {
		return
	}
	go func() {
		var tokens, lineno int
		for line := range chPipe {
			tokens += len(strings.Fields(line))
			lineno++
			fmt.Fprintf(fp, "%08d|%s\n", lineno, line)
		}
		fp.Close()
		chTokens <- tokens
	}()
	err = run(cmd, task.chKill, chPipe)
	if err != nil {
		return
	}
	tokens := <-chTokens

	quotumLock.Lock()
	quotum := 0
	gebruikt := 0
	rows, err = db.Query(fmt.Sprintf("SELECT `quotum` FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, user))
	if err != nil {
		quotumLock.Unlock()
		return
	}
	if rows.Next() {
		err = rows.Scan(&quotum)
		rows.Close()
		if err != nil {
			quotumLock.Unlock()
			return
		}
	} else {
		err = fmt.Errorf("MySQL: Can't find quotum")
		quotumLock.Unlock()
		return
	}
	if quotum > 0 {
		rows, err = db.Query(fmt.Sprintf("SELECT `nword` FROM `%s_info` WHERE `owner` = %q", Cfg.Prefix, user))
		if err != nil {
			quotumLock.Unlock()
			return
		}
		for rows.Next() {
			var n int
			err = rows.Scan(&n)
			if err != nil {
				quotumLock.Unlock()
				return
			}
			gebruikt += n
		}
	}

	if quotum > 0 && quotum-gebruikt < tokens {
		err = fmt.Errorf("Ruimte voor %d tokens. Nieuw corpus bevat %d tokens.", quotum-gebruikt, tokens)
		os.Remove(data)
		os.Remove(data + ".lines")
		quotumLock.Unlock()
		return
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `nword` = %d WHERE `id` = %q",
		Cfg.Prefix, tokens, task.id))
	quotumLock.Unlock()

	if err != nil {
		return
	}

	// kan er nog staan na onderbroken run
	os.RemoveAll(xml)

	var timeout string
	if Cfg.Timeout > 0 {
		timeout = fmt.Sprint("-t ", Cfg.Timeout)
	}
	cmd = shell(
		`alpino -a %s -d %s %s %s.lines > %s 2>> %s`,
		Cfg.Alpino, xml, timeout, data, stdout, stderr)
	err = run(cmd, task.chKill, nil)
	if err != nil {
		return
	}

	cmd = shell(
		// optie -w i.v.m. recover()
		`find %s -name '*.xml' | pqbuild -w %s %s %s 0 >> %s 2>> %s`,
		dirname,
		path.Base(dirname), quote(title), quote(user), stdout, stderr)
	err = run(cmd, task.chKill, nil)
	if err != nil {
		return
	}

	return
}

// Run command, maar onderbreek het als chKill gesloten is
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

	for {
		select {
		case err := <-chRet:
			return err
		case <-chKill:
			chKill = nil // niet opnieuw van lezen
			err = syscall.Kill(-pgid, 15)
			if err != nil {
				logf("syscall.Kill(-pgid, 15) error: %v", err)
			}
		}
	}
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
	if err == nil {
		logf("FINISHED: " + task.id)
		sendmail(user, "Corpus Ready", fmt.Sprintf("Your corpus \"%s\" is ready at %s", title, Cfg.Url))
	} else {
		logf("FAILED: %v, %v", task.id, err)
		if !task.killed {
			sendmail(user, "Corpus error", fmt.Sprintf("There was an error with your corpus \"%s\": %v", title, err))
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
	defer db.Close()

	ids := make([]string, 0)

	rows, err := db.Query("SELECT `id` FROM `" + Cfg.Prefix + "_info` WHERE `status` = \"QUEUED\" OR `status` = \"WORKING\" ORDER BY `created`")
	util.CheckErr(err)
	for rows.Next() {
		var id string
		util.CheckErr(rows.Scan(&id))
		ids = append(ids, id)
	}
	util.CheckErr(rows.Err())

	for _, id := range ids {
		_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `nword` = 0 WHERE `id` = %q", Cfg.Prefix, id))
		util.CheckErr(err)
	}

	for _, id := range ids {
		p := &Process{
			id:     id,
			chKill: make(chan bool),
			queued: true,
		}
		processLock.Lock()
		processes[id] = p
		processLock.Unlock()
		go func() {
			chWork <- p
		}()
	}
}
