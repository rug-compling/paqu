/*
Package compactcorpus is a reader for corpora in the compact Alpino format.

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

	i := len(name)
	if strings.HasSuffix(name, ".data.dz") {
		name = name[:i-8]
	} else if strings.HasSuffix(name, ".index") {
		name = name[:i-6]
	}
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

////////////////////////////////////////////////////////////////

var (
	list = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	index = []uint64{
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 62, 99, 99, 99, 63,
		52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 99, 99, 99, 99, 99, 99,
		99, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
		15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 99, 99, 99, 99, 99,
		99, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99,
	}
)

func encode(val uint64) string {
	chunks := 1
	if val != 0 {
		val2 := val
		for {
			val2 >>= 6
			if val2 == 0 {
				break
			}
			chunks++
		}
	}
	result := make([]uint8, chunks)
	for i := 0; i < chunks; i++ {
		shift := uint64(i * 6)
		mask := uint64(0x3f) << shift
		result[chunks-i-1] = list[(val&mask)>>shift]
	}
	return string(result)
}

func decode(val string) (uint64, error) {
	var result uint64
	var offset uint64

	for i := len(val) - 1; i >= 0; i-- {
		tmp := index[val[i]]
		if tmp == 99 {
			return 0, fmt.Errorf("Illegal character in base64 value: %v", val[i:i+1])
		}

		if (tmp<<offset)>>offset != tmp {
			return 0, fmt.Errorf("Type uint64 cannot store decoded base64 value: %v", val)
		}

		result |= tmp << offset
		offset += 6
	}
	return result, nil
}
