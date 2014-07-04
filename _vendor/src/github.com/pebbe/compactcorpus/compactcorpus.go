/*
Package compactcorpus provides a reader and writer for corpora in the compact Alpino format.

See: http://github.com/rug-compling/alpinocorpus
*/
package compactcorpus

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Corpus struct {
	opened bool
	name   string
	names  []string
	xml    [][]byte
	idx    map[string]int
}

type Range struct {
	c   *Corpus
	idx int
}

// Return an iterator for the corpus
func (c *Corpus) NewRange() (r *Range, err error) {
	if !c.opened {
		err = fmt.Errorf("Corpus '%v' is closed", c.name)
		return
	}
	r = &Range{
		c:   c,
		idx: 0,
	}
	return
}

// Return true if there are more items available
func (r *Range) HasNext() bool {
	return r.idx < len(r.c.idx)
}

// Return the next available item
func (r *Range) Next() (name string, xml []byte) {
	if r.idx == len(r.c.idx) {
		return
	}
	name = r.c.names[r.idx]
	xml = r.c.xml[r.idx]
	r.idx++
	return
}

// Get a named item from the corpus
func (c *Corpus) Get(name string) (xml []byte, err error) {
	if !c.opened {
		err = fmt.Errorf("Corpus '%v' is closed", c.name)
		return
	}
	i, ok := c.idx[name]
	if !ok {
		err = fmt.Errorf("Item '%v' not found in corpus '%v'", c.name, name)
		return
	}
	xml = c.xml[i]
	return
}

// Open a compact corpus for reading in simple mode.
// Use this version if you want to retrieve most or all items from the corpus.
// See also: RaOpen()
func Open(name string) (corpus *Corpus, err error) {
	corpus = &Corpus{
		names: make([]string, 0),
		xml:   make([][]byte, 0),
		idx:   make(map[string]int),
	}

	name = root(name)
	corpus.name = name

	var fp2 *os.File
	curfile := name + ".data.dz"
	fp2, err = os.Open(curfile)
	if err != nil {
		return
	}
	var gz *gzip.Reader
	gz, err = gzip.NewReader(fp2)
	if err != nil {
		return
	}
	var data []byte
	data, err = ioutil.ReadAll(gz)
	if err != nil {
		return
	}

	var fp *os.File
	curfile = name + ".index"
	fp, err = os.Open(curfile)
	if err != nil {
		return
	}
	defer fp.Close()
	lineno := 0
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		var offset, size uint64
		line := scanner.Text()
		lineno++
		a := strings.Fields(line)
		if len(a) != 3 {
			err = fmt.Errorf("Invalid number of fields in file '%s.index', line %d", curfile, lineno)
			return
		}
		corpus.idx[a[0]] = len(corpus.names)
		corpus.names = append(corpus.names, a[0])
		offset, err = decode(a[1])
		if err != nil {
			err = fmt.Errorf("%v in file '%s.index', line %d", curfile, lineno)
			return
		}
		size, err = decode(a[2])
		if err != nil {
			err = fmt.Errorf("%v in file '%s.index', line %d", curfile, lineno)
			return
		}
		if offset+size > uint64(len(data)) {
			err = fmt.Errorf("Data in file '%s.data.dz' is too short", name)
			return
		}
		corpus.xml = append(corpus.xml, data[offset:offset+size])
	}
	if err = scanner.Err(); err != nil {
		return
	}

	corpus.opened = true

	return
}

