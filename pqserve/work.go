package main

import (
	"github.com/pebbe/util"

	"database/sql"
	"errors"
	"fmt"
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

	_, err = db.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"WORKING\" WHERE `id` = %q",
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

	dirname := path.Join(Cfg.Data, task.id)
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

	// kan er nog staan na onderbroken run
	os.RemoveAll(xml)

	var cmd *exec.Cmd
	var f string
	switch params {
	case "run":
		f = "run"
	case "line":
		f = "lines"
	default:
		err = errors.New("Not implemented: " + params)
		logerr(err)
		return
	}
	cmd = shell(
		`alpino -a %s -d %s -f %s %s > %s 2> %s`,
		Cfg.Alpino,	xml, f, data, stdout, stderr)
	err = run(cmd, task.chKill)
	if err != nil {
		return
	}

	cmd = shell(
		// optie -w i.v.m. revocer()
		`find %s -name '*.xml' | pqbuild -w %s %s %s %s 0 >> %s 2>> %s`,
		dirname,
		os.Args[1], path.Base(dirname), quote(title), quote(user), stdout, stderr)
	err = run(cmd, task.chKill)
	if err != nil {
		return
	}

	return
}

// Run command, maar onderbreek het als input via chKill
func run(cmd *exec.Cmd, chKill chan bool) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	chRet := make(chan error)
	go func() {
		cmd.Start()
		chRet <- nil
		err := cmd.Wait()
		chRet <- err
	}()
	<-chRet
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
			err = syscall.Kill(-pgid, 15)
			if err != nil {
				logf("syscall.Kill(-pgid, 15) error: %v", err)
			}
		}
	}
	panic("Can't reach")
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

	defer delete(processes, task.id)

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
	if task, ok := processes[id]; ok {
		done := false

		task.lock.Lock()
		task.killed = true
		if task.queued {
			done = true
		}
		task.lock.Unlock()

		if done {
			logf("UNQUEUED: %v", id)
			delete(processes, task.id)
			return
		}

		task.chKill <- true
		for {
			time.Sleep(500 * time.Millisecond)
			if _, ok := processes[id]; !ok {
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
	rows, err := db.Query("SELECT `id` FROM `" + Cfg.Prefix + "_info` WHERE `status` = \"QUEUED\" OR `status` = \"WORKING\" ORDER BY `created`")
	util.CheckErr(err)
	for rows.Next() {
		var id string
		util.CheckErr(rows.Scan(&id))
		processes[id] = &Process{
			id:     id,
			chKill: make(chan bool, 10),
			queued: true,
		}
		go func() {
			chWork <- processes[id]
		}()
	}
	util.CheckErr(rows.Err())
}
