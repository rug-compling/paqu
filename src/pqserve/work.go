package main

import (
	"github.com/pebbe/util"
	"github.com/rug-compling/alud/v2"

	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	LeegArchief = errors.New("Leeg zipbestand of tarbestand")

	reTrailingSpace = regexp.MustCompile(">( |\t)*(\r\n|\n\r|\r|\n)")

	isTrue = map[string]bool{
		"true": true,
		"yes":  true,
		"ja":   true,
		"1":    true,
		"t":    true,
		"y":    true,
		"j":    true,
	}
)

func quote(s string) string {
	return "'" + strings.Replace(s, "'", "'\\''", -1) + "'"
}

func dowork(task *Process) (user string, title string, err error) {
	logf("WORKING: " + task.id)

	_, err = sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"WORKING\", `nword` = 0 WHERE `id` = %q",
		Cfg.Prefix, task.id))
	if err != nil {
		return
	}

	params := "unknown"
	isArch := false
	var rows *sql.Rows
	rows, err = sqlDB.Query(fmt.Sprintf("SELECT `description`,`owner`,`params` FROM `%s_info` WHERE `id` = %q",
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

	dirname := filepath.Join(paqudatadir, "data", task.id)
	data := filepath.Join(dirname, "data")
	xml := filepath.Join(dirname, "xml")
	dact := filepath.Join(dirname, "data.dact")
	stdout := filepath.Join(dirname, "stdout.txt")
	stderr := filepath.Join(dirname, "stderr.txt")
	summary := filepath.Join(dirname, "summary.txt")
	conllu := filepath.Join(dirname, "conllu")

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
		b, err = ar.ReadN(1000)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			ar.Close()
			return
		}

		if strings.Contains(string(b), "<alpino_ds") {
			if !strings.HasPrefix(params, "xmlzip") {
				params = "xmlzip"
			}
			setinvoer(params, task.id, false)
			ar.Close()
		} else {
			if strings.Contains(string(b), "<FoLiA") {
				params = "folia-arch"
			} else if strings.Contains(string(b), "<TEI") {
				params = "tei-arch"
			}
			isArch = true
			ar.Close()
			var fp *os.File
			fp, err = os.Create(data + ".unzip")
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
				if params == "folia-arch" || params == "tei-arch" {
					fmt.Fprintf(fp, "##PAQUFILE %s\n", hex.EncodeToString([]byte(ar.Name())))
				} else {
					fmt.Fprintf(fp, "\n##PAQUFILE %s\n", ar.Name())
				}
				var buf bytes.Buffer
				if params == "folia-arch" || params == "tei-arch" {
					err = ar.Copy(&buf)
				} else {
					err = ar.Copy(fp)
				}
				if err != nil {
					fp.Close()
					ar.Close()
					return
				}
				if params == "folia-arch" {
					folia(ar.Name()+".", &buf, fp)
				} else if params == "tei-arch" {
					tei(ar.Name()+".", &buf, fp)
				} else {
					fmt.Fprintln(fp)
				}
			}
			fp.Close()
			ar.Close()
			if params == "folia-arch" || params == "tei-arch" {
				os.Rename(data+".unzip", data+".tmp")
			}
		}
	}

	if params == "auto" {
		if isArch {
			params, err = invoersoort(data+".unzip", task.id)
		} else {
			params, err = invoersoort(data, task.id)
		}
		if err != nil {
			return
		}
		if params == "dact" {
			os.Rename(data, dact)
		}
	}

	if isArch {
		setinvoer(params, task.id, true)
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
		os.Remove(data + ".tmp2")
		for _, f := range []string{data, data + ".lines", stdout, stderr, summary} {
			if f == data && (params == "dact" || strings.HasPrefix(params, "xmlzip")) {
				continue
			}
			logerr(gz(f))
		}
		fnames, e := filenames2(xml, false)
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
			logerr(gz(filepath.Join(xml, fname)))
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
		tokens, nlines, err = unpackDact(data, xml, dact, conllu, stderr, task.chKill)
		if err != nil {
			return
		}
		err = do_quotum(task.id, user, tokens, nlines)
		if err != nil {
			os.Remove(dact)
			os.Remove(dact + "x")
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
		err = do_quotum(task.id, user, tokens, nlines)
		if err != nil {
			os.Remove(data + ".lines")
			os.Remove(stderr)
			os.RemoveAll(xml)
			return
		}
	} else { // if params != (dact || xmlzip)
		reuse := false
		reuse_more := false
		if files, e := filenames2(xml, false); e == nil && len(files) > 0 {
			reuse = true
			if isArch {
				os.Remove(data + ".unzip")
			}
			done := make(map[string]bool)
			for _, f := range files {
				b, e := ioutil.ReadFile(filepath.Join(xml, f))
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

			_, err = sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `nword` = %d WHERE `id` = %q",
				Cfg.Prefix, nword, task.id))
			if err != nil {
				return
			}

		} else { // if !reuse
			var has_tok, has_lbl, is_xml, is_arch bool
			if strings.Contains(params, "-arch") {
				is_arch = true
			}
			if strings.HasPrefix(params, "folia") || strings.HasPrefix(params, "tei") {
				is_xml = true
				has_lbl = true
				has_tok = true
			}
			if strings.Contains(params, "-lbl") {
				has_lbl = true
			}
			if strings.Contains(params, "-tok") {
				has_tok = true
			}

			if is_xml {
				if !is_arch {
					var fpin, fpout *os.File
					fpin, err = os.Open(data)
					if err != nil {
						return
					}
					fpout, err = os.Create(data + ".tmp")
					if err != nil {
						fpin.Close()
						return
					}
					if params == "folia" {
						err = folia("", fpin, fpout)
					} else {
						err = tei("", fpin, fpout)
					}
					fpout.Close()
					fpin.Close()
					if err != nil {
						return
					}
				}
			} else {
				// pqtexter
				var pqtexter, unzip string
				if params == "run" {
					pqtexter = "-r"
				} else if has_lbl {
					pqtexter = "-l"
				}
				if isArch {
					unzip = ".unzip"
				}
				err = shell("pqtexter %s %s%s > %s.tmp 2>> %s", pqtexter, data, unzip, data, stderr).Run()
				if err != nil {
					return
				}
				if isArch {
					os.Remove(data + ".unzip")
				}
			}

			// verwijderen van commentaren en labels zonder tekst
			if params == "run" || strings.HasPrefix(params, "line") {
				var fpin, fpout *os.File
				fpin, err = os.Open(data + ".tmp")
				if err != nil {
					return
				}
				fpout, err = os.Create(data + ".tmp2")
				if err != nil {
					fpin.Close()
					return
				}
				rd := util.NewReader(fpin)
				for {
					line, e := rd.ReadLineString()
					if e != nil {
						break
					}
					if strings.TrimSpace(line) == "" || line[0] == '%' || reRunLabel.MatchString(line) {
						continue
					}
					fmt.Fprintln(fpout, line)
				}
				fpout.Close()
				fpin.Close()
				os.Remove(data + ".tmp")
				os.Rename(data+".tmp2", data+".tmp")
			}

			// ontdubbelen van labels
			if strings.HasPrefix(params, "folia") || strings.HasPrefix(params, "tei") || strings.Contains(params, "-lbl") {
				var fpin, fpout *os.File
				var dubbel bool
				fpin, err = os.Open(data + ".tmp")
				if err != nil {
					return
				}
				fpout, err = os.Create(data + ".tmp2")
				if err != nil {
					fpin.Close()
					return
				}
				rd := util.NewReader(fpin)
				seen := make(map[string]bool)
				for {
					line, e := rd.ReadLineString()
					if e != nil {
						break
					}
					if strings.HasPrefix(line, "##PAQULBL") {
						var b []byte
						b, err = hex.DecodeString(strings.Fields(line)[1])
						if err != nil {
							fpout.Close()
							fpin.Close()
							return
						}
						lbl := string(b)
						if seen[lbl] {
							dubbel = true
							for i := 1; true; i++ {
								lbl2 := fmt.Sprintf("%s.dup.%d", lbl, i)
								if !seen[lbl2] {
									lbl = lbl2
									break
								}
							}
							line = "##PAQULBL " + hex.EncodeToString([]byte(lbl))
						}
						seen[lbl] = true
					}
					fmt.Fprintln(fpout, line)
				}
				fpout.Close()
				fpin.Close()
				if dubbel {
					os.Remove(data + ".tmp")
					os.Rename(data+".tmp2", data+".tmp")
				} else {
					os.Remove(data + ".tmp2")
				}
			}

			// tokenizer
			if !has_tok {
				var tok string
				if params == "run" {
					tok = "tokenize.sh"
				} else {
					tok = "tokenize_no_breaks.sh"
				}
				err = shell("$ALPINO_HOME/Tokenization/%s < %s.tmp > %s.tmp2 2>> %s", tok, data, data, stderr).Run()
				if err != nil {
					return
				}
				os.Rename(data+".tmp2", data+".tmp")
			}

			var fp, fpin *os.File
			fpin, err = os.Open(data + ".tmp")
			if err != nil {
				return
			}
			fp, err = os.Create(data + ".lines")
			if err != nil {
				fpin.Close()
				return
			}

			metalines := make([]string, 0)
			inmeta := false

			rd := util.NewReader(fpin)
			var filename, lbl string
			var tokens, nlines, i int
			var metaseen map[string]bool
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
						metalines = metalines[0:0]
						pp := strings.Split(filename, "/")
						for i, p := range pp {
							metalines = append(metalines, fmt.Sprintf("text\tpaqu.path%d\t%s", i+1, p))
						}
						metaseen = make(map[string]bool)
						inmeta = true
					} else if a[0] == "##PAQULBL" {
						lbl = val
					}
				} else if strings.HasPrefix(line, "##META") {
					if !inmeta {
						metaseen = make(map[string]bool)
						inmeta = true
					}
					a := strings.Fields(line)
					if len(a) == 2 {
						b, e := hex.DecodeString(a[1])
						if e != nil {
							err = e
							fp.Close()
							fpin.Close()
							return
						}
						f := strings.Fields(string(b))
						if len(f) > 2 && f[2] == "=" {
							// als deze voor het eerst, dan alle oude met dezelfde naam wegdoen
							if !metaseen[f[1]] {
								metaseen[f[1]] = true
								for i := 0; i < len(metalines); i++ {
									if i == 0 && filename != "" {
										continue
									}
									if a := strings.Split(metalines[i], "\t"); a[1] == f[1] {
										metalines = append(metalines[:i], metalines[i+1:]...)
										i--
									}
								}
							}
							if len(f) > 3 {
								value := strings.Join(f[3:], " ")
								if f[0] == "bool" {
									if isTrue[strings.ToLower(value)] {
										value = "true"
									} else {
										value = "false"
									}
								}
								metalines = append(metalines, fmt.Sprintf("%s\t%s\t%s", f[0], f[1], value))
							}
						}
					}
				} else {
					inmeta = false
					tokens += len(strings.Fields(line))
					var fname string
					if lbl == "" {
						if filename != "" {
							i++
							fname = fmt.Sprintf("%04d/%04d-%s.%d", nlines/10000, nlines%10000, encode_filename(filename), i)
						} else {
							fname = fmt.Sprintf("%04d/%04d", nlines/10000, nlines%10000)
						}
					} else {
						fname = fmt.Sprintf("%04d/%04d-%s", nlines/10000, nlines%10000, encode_filename(lbl))
						lbl = ""
					}
					fmt.Fprintf(fp, "%s|%s\n", fname, strings.TrimSpace(line))
					nlines++
					if len(metalines) > 0 {
						name := filepath.Join(xml, fname+".meta")
						dir := filepath.Dir(name)
						os.MkdirAll(dir, 0777)
						fpm, e := os.Create(filepath.Join(xml, fname+".meta"))
						if e != nil {
							err = e
							fp.Close()
							fpin.Close()
							return
						}
						for _, m := range metalines {
							fmt.Fprintln(fpm, m)
						}
						fpm.Close()
					}
				}
			}
			fp.Close()
			fpin.Close()
			os.Remove(data + ".tmp")

			err = do_quotum(task.id, user, tokens, nlines)
			if err != nil {
				os.Remove(data)
				os.Remove(data + ".lines")
				os.Remove(stderr)
				os.RemoveAll(xml)
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

			var server, timeout string
			if Cfg.Alpinoserver != "" {
				server = "-s " + Cfg.Alpinoserver
			}
			if Cfg.Timeout > 0 {
				timeout = fmt.Sprint("-t ", Cfg.Timeout)
			}
			ud1 := ""
			ud2 := ""
			if Cfg.Conllu {
				ud1 = "-u " + conllu + ".err"
				ud2 = "; echo " + alud.VersionID() + " > " + conllu + ".version"
			}
			cmd := shell(
				`pqalpino %s -e half -l -T -q -n %d -d %s %s %s %s.lines%s >> %s 2>> %s%s`,
				ud1, Cfg.Maxtokens, xml, server, timeout, data, ext, stdout, stderr, ud2)
			err = run(cmd, task.chKill, nil)
			if err != nil {
				return
			}
		}
	} // end if params != (dact || xmlzip)

	nlines := 0
	fp, e := os.Open(data + ".lines")
	if e != nil {
		err = e
		return
	}
	rd := util.NewReader(fp)
	for {
		_, e := rd.ReadLine()
		if e != nil {
			fp.Close()
			if e == io.EOF {
				break
			} else {
				err = e
				return
			}
		}
		nlines++
	}

	errlines := make([]string, 0)
	fp, e = os.Open(stderr)
	if e != nil {
		err = e
		return
	}
	rd = util.NewReader(fp)
	errseen := make(map[string]bool)
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
			if params == "run" || strings.HasPrefix(params, "line") || strings.HasPrefix(params, "folia") || strings.HasPrefix(params, "tei") {
				fname = decode_filename(a[0][2:])
			}
			if strings.Contains(params, "-lbl") || strings.HasPrefix(params, "folia") || strings.HasPrefix(params, "tei") {
				fname = fname[1+strings.Index(fname, "-"):]
			}
			// bij herstart worden mislukte zinnen opnieuw geparst, en mislukken dan opnieuw
			errline := fname + "\t" + a[ln-3] + "\t" + a[ln-2] + "\t" + a[1] + "\n"
			if !errseen[errline] {
				errlines = append(errlines, errline)
				errseen[errline] = true
			}
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
	if strings.Contains(params, "-lbl") || strings.HasPrefix(params, "folia") || strings.HasPrefix(params, "tei") || isArch {
		p += "[0-9]+/[0-9]+-"
		d = "-d"
	} else if params == "dact" || strings.HasPrefix(params, "xmlzip") {
		p += "[0-9]+/"
		d = "-d"
	}

	filenames, e := filenames2(xml, true)
	if e != nil {
		err = e
		return
	}
	for _, filename := range filenames {
		m := filepath.Join(xml, filename)
		x := m[:len(m)-4] + "xml"
		var xb, mb []byte
		xb, err = ioutil.ReadFile(x)
		if err != nil {
			os.Remove(m)
			continue
		}
		mb, err = ioutil.ReadFile(m)
		if err != nil {
			return
		}
		fp, err = os.Create(x + ".tmp")
		if err != nil {
			return
		}
		xt := string(xb)
		mt := strings.Split(strings.TrimSpace(string(mb)), "\n")
		i := strings.Index(xt, "<parser")
		if i < 0 {
			i = strings.Index(xt, "<node")
		}
		fmt.Fprint(fp, xt[:i], "<metadata>\n")
		for _, m := range mt {
			mm := strings.Split(m, "\t")
			fmt.Fprintf(fp,
				"    <meta type=%q name=%q value=%q/>\n",
				html.EscapeString(strings.ToLower(mm[0])),
				html.EscapeString(mm[1]),
				html.EscapeString(mm[2]))
		}
		fmt.Fprint(fp, "  </metadata>\n  ", xt[i:])
		fp.Close()
		os.Rename(x+".tmp", x)
		os.Remove(m)
	}

	if Cfg.Conllu && strings.HasPrefix(params, "xmlzip") {
		cmd := shell(
			`find %s -name '*.xml' | sort > %s.list; pqudep -p %s/data/ -l %s.list -o > /dev/null 2> %s.err; rm %s.list ; pqudep -v > %s.version`,
			dirname, conllu, paqudatadir, conllu, conllu, conllu, conllu)
		err = run(cmd, task.chKill, nil)
		if err != nil {
			return
		}
		if cu, _ := os.Stat(conllu + ".err"); cu.Size() != 0 {
			sysErr(fmt.Errorf("CONLLU error(s) in %s.err", conllu))
		}
	}

	cmd := shell(
		// optie -w i.v.m. recover_work()
		`find %s -name '*.xml' | sort | pqbuild -w -p %s %s -s -m %s %s %s %s 0 >> %s 2>> %s`,
		dirname,
		quote(p), d, quote(task.info),
		filepath.Base(dirname), quote(title), quote(user), stdout, stderr)
	err = run(cmd, task.chKill, nil)
	if err != nil {
		return
	}

	if Cfg.Dact && params != "dact" {
		p := ""
		if strings.Contains(params, "-lbl") || strings.HasPrefix(params, "folia") || strings.HasPrefix(params, "tei") || isArch {
			p = "-"
		} else if strings.HasPrefix(params, "xmlzip") {
			p = "/"
		}
		err = makeDact(dact, conllu, xml, p, task.chKill)
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
	for _, sig := range []int{15, 9} {
		err = syscall.Kill(-pgid, syscall.Signal(sig))
		if err != nil {
			logf("syscall.Kill(-pgid, %d) error: %v", sig, err)
		}
		if sig != 9 {
			time.Sleep(2 * time.Second)
		}
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

	user, title, err := dowork(task)

	select {
	case <-chGlobalExit:
		return
	default:
	}

	if err == nil {
		logf("FINISHED: " + task.id)
		sendmail(user, "Corpus klaar", fmt.Sprintf("Je corpus \"%s\" staat klaar op %s", title, urlJoin(Cfg.Url, "/?db="+task.id)))
		var params string
		rows, err := sqlDB.Query(fmt.Sprintf("SELECT `params` FROM `%s_info` WHERE `id` = %q", Cfg.Prefix, task.id))
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
		sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"FINISHED\", `msg` = %q WHERE `id` = %q", Cfg.Prefix, msg, task.id))
	} else {
		logf("FAILED: %v, %v", task.id, err)
		if !task.killed {
			sendmail(user, "Corpus fout", fmt.Sprintf("Er ging iets fout met je corpus \"%s\": %v", title, err))
		}
		sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `status` = \"FAILED\", `msg` = %q WHERE `id` = %q", Cfg.Prefix, err.Error(), task.id))
	}
}

func kill(id string) {
	chDelete <- id
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

func recover_work() {

	ids := make([][2]string, 0)
	queuing := make([]string, 0)

	_, err := sqlDB.Exec(
		"UPDATE `" + Cfg.Prefix + "_info` SET `nword` = 0, `status` = \"QUEUED\" WHERE `status` = \"WORKING\" AND `owner` LIKE \"%@%\"")
	util.CheckErr(err)

	rows, err := sqlDB.Query("SELECT `id`,`status`,`owner` FROM `" +
		Cfg.Prefix + "_info` WHERE `status` = \"QUEUED\" OR `status` = \"QUEUING\"  ORDER BY `created`")
	util.CheckErr(err)
	for rows.Next() {
		var id, status, owner string
		util.CheckErr(rows.Scan(&id, &status, &owner))
		if status == "QUEUED" {
			ids = append(ids, [2]string{id, owner})
		} else {
			queuing = append(queuing, id)
		}
	}
	util.CheckErr(rows.Err())

	if len(queuing) > 0 {
		_, err = sqlDB.Exec(fmt.Sprintf("DELETE FROM `%s_info` WHERE `status` = \"QUEUING\"", Cfg.Prefix))
		util.CheckErr(err)
	}
	for _, corpus := range queuing {
		util.CheckErr(os.RemoveAll(filepath.Join(paqudatadir, "data", corpus)))
		logf("QUEUING: rm -r %s: ok", filepath.Join(paqudatadir, "data", corpus))
	}

	for _, id := range ids {
		p := &Process{
			id:     id[0],
			owner:  id[1],
			chKill: make(chan bool),
			queued: true,
		}
		processLock.Lock()
		processes[id[0]] = p
		processLock.Unlock()
		chWork <- p
	}
}

func do_quotum(id, user string, tokens, nlines int) error {
	quotumLock.Lock()
	defer quotumLock.Unlock()

	quotum := 0
	gebruikt := 0
	rows, err := sqlDB.Query(fmt.Sprintf("SELECT `quotum` FROM `%s_users` WHERE `mail` = %q", Cfg.Prefix, user))
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
		rows, err = sqlDB.Query(fmt.Sprintf("SELECT `nword` FROM `%s_info` WHERE `owner` = %q", Cfg.Prefix, user))
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

	_, err = sqlDB.Exec(fmt.Sprintf("UPDATE `%s_info` SET `nword` = %d, `nline` = %d WHERE `id` = %q",
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

		// sanitize (bug on https://webservices-lst.science.ru.nl/portal/ )
		//bindata = reTrailingSpace.ReplaceAll(bindata, []byte(">\000"))
		bindata = bytes.Replace(bindata, []byte("\r"), []byte(""), -1)
		//bindata = bytes.Replace(bindata, []byte("\n"), []byte(""), -1)
		//bindata = bytes.Replace(bindata, []byte(">\000"), []byte(">\n"), -1)

		nline++

		if nline%10000 == 1 {
			nd++
			sdir = fmt.Sprintf("%04d", nd)
			os.Mkdir(filepath.Join(xmldir, sdir), 0777)
		}

		name := ar.Name()
		if strings.HasSuffix(strings.ToLower(name), ".xml") {
			name = name[:len(name)-4]
		}
		encname := encode_filename(name)

		fp, err := os.Create(filepath.Join(xmldir, sdir, encname+".xml"))
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
		if alpino.Comments != nil {
			for _, c := range alpino.Comments {
				if strings.HasPrefix(c, "Q#") {
					a := strings.SplitN(c, "|", 2)
					if len(a) == 2 {
						fmt.Fprintf(fperr, "Q#%s|%s\n", name, strings.TrimSpace(a[1]))
						break
					}
				}
			}
		}
	}

	return tokens, nline, nil
}
