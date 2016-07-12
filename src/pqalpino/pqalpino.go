package main

import (
	"github.com/pebbe/util"

	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Response struct {
	Code     int
	Status   string
	Message  string
	Id       string
	Interval int
	Finished bool
	Batch    []Line
}

type Line struct {
	Status   string
	Lineno   int
	Label    string
	Sentence string
	Xml      string
	Log      string
}

var (
	opt_a = flag.String("a", "", "ALPINO_HOME")
	opt_d = flag.String("d", "xml", "Directory voor uitvoer")
	opt_s = flag.String("s", "", "Alpino-server")
	opt_t = flag.Int("t", 900, "Time-out in seconden per regel")

	x = util.CheckErr
)

func usage() {
	fmt.Fprintf(os.Stderr, `
Syntax: %s [opties] datafile

Verplichte optie:
  -a directory : ALPINO_HOME

Overige opties:
  -d directory : Directory waar uitvoer wordt geplaatst (default: xml)
  -s server    : Alpino-server, zie: https://github.com/rug-compling/alpino-api
                 Als deze ontbreekt wordt een lokale versie van Alpino gebruikt
  -t seconden  : Time-out per regel (default: 900)

`, os.Args[0])
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		return
	}
	filename := flag.Arg(0)

	if *opt_a == "" {
		fmt.Println("Optie -a ontbreekt")
		return
	}

	os.Mkdir(*opt_d, 0777)

	// PARSEN

	if *opt_s == "" {
		var fpin, fpout *os.File
		var errval error
		tmpfile := filename + ".part"
		defer func() {
			if fpin != nil {
				fpin.Close()
			}
			if fpout != nil {
				fpout.Close()
			}
			os.Remove(tmpfile)
			if errval != io.EOF {
				x(errval)
			}
		}()
		fpin, errval = os.Open(filename)
		if errval != nil {
			return
		}
		rd := util.NewReaderSize(fpin, 5000)
		n := 0
		for {
			line, err := rd.ReadLineString()
			if err != nil && err != io.EOF {
				errval = err
				return
			}
			if err == nil && n == 0 {
				fpout, errval = os.Create(tmpfile)
				if errval != nil {
					return
				}
			}
			if err == nil {
				fmt.Fprintln(fpout, line)
				n++
			}
			if (err == io.EOF && n > 0) || n == 10000 {
				n = 0
				fpout.Close()
				fpout = nil
				cmd := exec.Command(
					"/bin/bash",
					"-c",
					fmt.Sprintf(
						"$ALPINO_HOME/bin/Alpino -veryfast -flag treebank %s debug=1 end_hook=xml user_max=%d -parse < %s",
						*opt_d, *opt_t*1000, tmpfile))
				cmd.Env = []string{
					"ALPINO_HOME=" + *opt_a,
					"PATH=" + os.Getenv("PATH"),
					"LANG=en_US.utf8",
					"LANGUAGE=en_US.utf8",
					"LC_ALL=en_US.utf8",
				}
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				errval = cmd.Run()
				if errval != nil {
					return
				}
			}
			if err == io.EOF {
				break
			}
		}
	} else {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, `{"request":"parse", "lines":true, "tokens":true, "escape":"none", "timeout":%d}`, *opt_t)
		fp, err := os.Open(filename)
		x(err)
		_, err = io.Copy(&buf, fp)
		fp.Close()
		x(err)
		resp, err := http.Post(*opt_s, "application/json", &buf)
		util.CheckErr(err)
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		util.CheckErr(err)
		var response Response
		err = json.Unmarshal(data, &response)
		util.CheckErr(err)
		if response.Code > 299 {
			x(fmt.Errorf("%d %s -- %s", response.Code, response.Status, response.Message))
		}
		maxinterval := response.Interval
		id := response.Id

		go func() {
			chSignal := make(chan os.Signal, 1)
			signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
			sig := <-chSignal
			fmt.Fprintf(os.Stderr, "Signal: %v\n", sig)

			var buf bytes.Buffer
			fmt.Fprintf(&buf, `{"request":"cancel", "id":%q}`, id)
			resp, err := http.Post(*opt_s, "application/json", &buf)
			util.CheckErr(err)
			_, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			util.CheckErr(err)

			os.Exit(0)
		}()

		interval := 2
		for {
			if interval > maxinterval {
				interval = maxinterval
			}
			if interval > 120 {
				interval = 120
			}
			time.Sleep(time.Duration(interval) * time.Second)

			var buf bytes.Buffer
			fmt.Fprintf(&buf, `{"request":"output", "id":%q}`, id)
			resp, err := http.Post(*opt_s, "application/json", &buf)
			util.CheckErr(err)
			data, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			util.CheckErr(err)
			var response Response
			err = json.Unmarshal(data, &response)
			util.CheckErr(err)
			if response.Code > 299 {
				x(fmt.Errorf("%d %s -- %s", response.Code, response.Status, response.Message))
			}
			for _, line := range response.Batch {
				if line.Status == "ok" {
					fp, err := os.Create(filepath.Join(*opt_d, line.Label+".xml"))
					x(err)
					fmt.Fprintln(fp, line.Xml)
					fp.Close()
				} else {
					fmt.Fprintf(os.Stderr, `**** parsing %s (line number %d)
%s
Q#%s|%s|%s|??|????
**** parsed %s (line number %d)
`,
						line.Label, line.Lineno,
						line.Log,
						line.Label, line.Sentence, line.Status,
						line.Label, line.Lineno)
				}
			}

			if response.Finished {
				break
			}
			interval = (3 * interval) / 2
		}
	}
}
