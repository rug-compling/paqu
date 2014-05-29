package main

import (
	"github.com/pebbe/util"

	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

func logf(format string, v ...interface{}) {
	chLog <- fmt.Sprintf(format, v...)
}

func logerr(err error) {
	if err == nil {
		return
	}
	var msg string
	_, filename, lineno, ok := runtime.Caller(1)
	if ok {
		msg = fmt.Sprintf("%v:%v: %v", path.Base(filename), lineno, err.Error())
	} else {
		msg = err.Error()
	}
	chLog <- msg
	return
}

func logger() {

	logfile := path.Join(paqudir, "pqserve.log")

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
	for msg := range chLog {
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
	}
}
