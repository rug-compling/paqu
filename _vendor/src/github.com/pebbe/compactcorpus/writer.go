package compactcorpus

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	errEmptyName = errors.New("name is empty")
)

type Writer struct {
	size   uint64
	idx    *os.File
	dz     *io.PipeWriter
	chRet  chan error
	err    error
	opened bool
}

// Start creating a new corpus.
// The corpus will not be written completely until Close() is called.
func NewCorpus(filename string) (*Writer, error) {
	name := root(filename)

	fp, err := os.Create(name + ".index")
	if err != nil {
		return nil, err
	}

	rdr, wtr := io.Pipe()

	ret := make(chan error, 1)

	go func() {
		err := write(rdr, name+".data.dz", 9)
		ret <- err
	}()

	w := &Writer{
		idx:    fp,
		dz:     wtr,
		chRet:  ret,
		opened: true,
	}
	return w, nil
}

func (w *Writer) Close() error {
	if w.opened {
		w.opened = false
		w.idx.Close()
		w.dz.Close()
		w.err = <-w.chRet
	}
	return w.err
}

func (w *Writer) Write(name string, xml []byte) error {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return errEmptyName
	}

	w.dz.Write(xml)

	ln := uint64(len(xml))
	off := encode(w.size)
	size := encode(ln)
	fmt.Fprintf(w.idx, "%s\t%s\t%s\n", name, off, size)

	w.size += ln

	return nil
}

func (w *Writer) WriteString(name string, xml string) error {
	return w.Write(name, []byte(xml))
}

func (w *Writer) WriteFile(filename string) error {
	xml, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return w.Write(filename, xml)
}
