package main

import (
	"fmt"
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
