package main

import (
	"github.com/pebbe/util"

	"fmt"
	"html"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func logf(format string, v ...interface{}) {
	chLog <- fmt.Sprintf(format, v...)
}

func logerr(err error) bool {
	if err == nil {
		return false
	}
	var msg string
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		msg = fmt.Sprintf("%v:%v: %v", filepath.Base(filename), lineno, err.Error())
	} else {
		msg = err.Error()
	}
	chLog <- msg
	return true
}

func logerrfrag(q *Context, err error) bool {
	if err == nil {
		return false
	}
	var msg string
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		msg = fmt.Sprintf("%v:%v: %v", filepath.Base(filename), lineno, err.Error())
	} else {
		msg = err.Error()
	}
	chLog <- msg
	fmt.Fprintf(q.w, "<div class=\"warning\">%s</div>\n", html.EscapeString(err.Error()))
	return true
}

func logger() {

	logfile := filepath.Join(paqudir, "pqserve.log")

	rotate := func() {
		for i := 4; i > 1; i-- {
			os.Rename(
				fmt.Sprintf("%s%d", logfile, i-1),
				fmt.Sprintf("%s%d", logfile, i))
		}
		os.Rename(logfile, logfile+"1")
	}

	rotate()
	fp, err := os.Create(logfile)
	util.CheckErr(err)

	n := 0
	for {
		select {
		case msg := <-chLog:
			now := time.Now()
			s := fmt.Sprintf("%04d-%02d-%02d %d:%02d:%02d %s", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), msg)
			fmt.Fprintln(fp, s)
			fp.Sync()
			if verbose {
				fmt.Println(s)
			}
			n++
			if n == 10000 {
				fp.Close()
				rotate()
				fp, _ = os.Create(logfile)
				n = 0
			}
		case <-chLoggerExit:
			fp.Close()
			return
		}
	}
}
