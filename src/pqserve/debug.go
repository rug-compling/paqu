package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type Info struct {
	HasDbXML     bool
	QueueLen     int
	QueueCap     int
	NumCPU       int
	NumCgoCall   int64
	NumGoroutine int
	Uptime       time.Duration
	UptimeString string
	Version      string
}

func GetInfo() interface{} {
	d := time.Now().Sub(started)
	return Info{
		HasDbXML:     has_dbxml,
		QueueLen:     len(chWork),
		QueueCap:     cap(chWork),
		NumCPU:       runtime.NumCPU(),
		NumCgoCall:   runtime.NumCgoCall(),
		NumGoroutine: runtime.NumGoroutine(),
		Uptime:       d,
		UptimeString: d.String(),
		Version:      runtime.Version(),
	}
}

func environment(q *Context) {
	q.w.Header().Set("Content-Type", "text/plain")
	cmd := shell("env | sort")
	cmd.Stdout = q.w
	e := cmd.Run()
	if e != nil {
		fmt.Fprintln(q.w, e)
	}
}

func stacktrace(q *Context) {
	q.w.Header().Set("Content-Type", "text/plain")

	cmd := shell(fmt.Sprintf("top -b -n 1 -p %d", os.Getpid()))
	cmd.Stdout = q.w
	cmd.Stderr = q.w
	cmd.Run()

	fmt.Fprintf(q.w, "\n\n%d goroutines\n\n", runtime.NumGoroutine())

	stacktrace := make([]byte, 50000)
	length := runtime.Stack(stacktrace, true)
	fmt.Fprintln(q.w, string(stacktrace[:length]))
}
