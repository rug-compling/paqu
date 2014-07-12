package compactcorpus

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type RaCorpus struct {
	opened bool
	name   string
	names  []string
	values [][2]int64
	idx    map[string]int
	fp     *os.File
	dz     *reader
}

type RaRange struct {
	c   *RaCorpus
	idx int
}

// Return an iterator for the corpus
func (c *RaCorpus) NewRange() (r *RaRange, err error) {
	if !c.opened {
		err = fmt.Errorf("Corpus '%v' is closed", c.name)
		return
	}
	r = &RaRange{
		c:   c,
		idx: 0,
	}
	return
}

// Return true if there are more items available
func (r *RaRange) HasNext() bool {
	return r.idx < len(r.c.idx)
}

// Return the next available item
func (r *RaRange) Next() (name string, xml []byte, err error) {
	if r.idx == len(r.c.idx) {
		return
	}
	name = r.c.names[r.idx]
	v := r.c.values[r.idx]
	xml, err = r.c.dz.get(v[0], v[1])
	r.idx++
	return
}

// Get a named item from the corpus
func (c *RaCorpus) Get(name string) (xml []byte, err error) {
	if !c.opened {
		err = fmt.Errorf("Corpus '%v' is closed", c.name)
		return
	}
	i, ok := c.idx[name]
	if !ok {
		err = fmt.Errorf("Item '%v' not found in corpus '%v'", c.name, name)
		return
	}
	v := c.values[i]
	xml, err = c.dz.get(v[0], v[1])
	return
}

// Close the corpus. This is necessary for freeing files
func (c *RaCorpus) Close() {
	if c.opened {
		c.opened = false
		c.fp.Close()
	}
}

// Open a compact corpus for reading in random access mode.
//
// You need to call Close() when you are done.
//
// Use this version if you want to retrieve only a few items from the
// corpus, or if you don't have enough memory to unpack the complete corpus
// into memory.
// Otherwise, use: Open()
func RaOpen(name string) (corpus *RaCorpus, err error) {
	corpus = &RaCorpus{
		names:  make([]string, 0),
		values: make([][2]int64, 0),
		idx:    make(map[string]int),
	}

	name = root(name)
	corpus.name = name

	var fp *os.File
	curfile := name + ".index"
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
		corpus.values = append(corpus.values, [2]int64{int64(offset), int64(size)})
	}
	if err = scanner.Err(); err != nil {
		return
	}

	curfile = name + ".data.dz"
	corpus.fp, err = os.Open(curfile)
	if err != nil {
		return
	}
	corpus.dz, err = newReader(corpus.fp)
	if err != nil {
		corpus.fp.Close()
		return
	}

	corpus.opened = true

	runtime.SetFinalizer(corpus, (*RaCorpus).Close)

	return
}
