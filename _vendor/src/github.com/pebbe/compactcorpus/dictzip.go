package compactcorpus

// dit is gecopieerd uit dictzip

//. Imports

import (
	"compress/flate"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

//. Writer

/*
Levels range from 1 (BestSpeed) to 9 (BestCompression), Level 0 (NoCompression), -1 (DefaultCompression)
*/
func write(r io.Reader, filename string, level int) error {

	const blocksize = 58315

	dirname := path.Dir(filename)
	root := path.Base(filename)
	fptmp, err := ioutil.TempFile(dirname, root)
	if err != nil {
		return err
	}
	defer func(name string) {
		fptmp.Close()
		os.Remove(name)
	}(fptmp.Name())

	crc := crc32.NewIEEE()
	isize := 0

	fw, err := flate.NewWriter(fptmp, level)
	if err != nil {
		return err
	}
	sizes := make([]int64, 0)
	b := make([]byte, blocksize)
	total := int64(0)
	eof := false
	for !eof {
		n, err := io.ReadFull(r, b)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			} else {
				eof = true
			}
		}
		if n > 0 {
			crc.Write(b[:n])
			isize += n

			fw.Write(b[:n])
			fw.Flush()
			fw.Reset(fptmp)

			s, _ := fptmp.Stat()
			l := s.Size()
			sizes = append(sizes, l-total)
			total = l
		}
	}
	fw.Close()

	fp, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	xfl := byte(0)
	if level == flate.BestCompression {
		xfl = 2
	} else if level == flate.BestSpeed {
		xfl = 4
	}
	now := time.Now().Unix()
	_, err = fp.Write([]byte{
		31, 139, 8, 4,
		byte(now & 255), byte((now >> 8) & 255), byte((now >> 16) & 255), byte((now >> 24) & 255),
		xfl, 255})
	if err != nil {
		return err
	}

	xlen := 10 + 2*len(sizes)
	ln := 6 + 2*len(sizes)
	_, err = fp.Write([]byte{
		byte(xlen & 255), byte((xlen >> 8) & 255),
		'R', 'A', byte(ln & 255), byte((ln >> 8) & 255),
		1, 0,
		byte(blocksize & 255), byte((blocksize >> 8) & 255),
		byte(len(sizes) & 255), byte((len(sizes) >> 8) & 255)})
	if err != nil {
		return err
	}
	for _, o := range sizes {
		_, err = fp.Write([]byte{byte(o & 255), byte((o >> 8) & 255)})
		if err != nil {
			return err
		}
	}

	fptmp.Seek(0, 0)
	io.Copy(fp, fptmp)

	c := crc.Sum32()
	_, err = fp.Write([]byte{
		byte(c & 255), byte((c >> 8) & 255), byte((c >> 16) & 255), byte((c >> 24) & 255),
		byte(isize & 255), byte((isize >> 8) & 255), byte((isize >> 16) & 255), byte((isize >> 24) & 255),
	})
	if err != nil {
		return err
	}

	return nil
}

//. Reader

type reader struct {
	fp        io.ReadSeeker
	offsets   []int64
	blocksize int64
	lock      sync.Mutex
}

func newReader(rs io.ReadSeeker) (*reader, error) {

	dz := &reader{fp: rs}

	_, err := dz.fp.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	metadata := []byte{}

	p := 0

	h := make([]byte, 10)
	n, err := io.ReadFull(dz.fp, h)
	if err != nil {
		return nil, err
	}
	p += n

	if h[0] != 31 || h[1] != 139 {
		return nil, fmt.Errorf("Invalid header: %02X %02X\n", h[0], h[1])
	}

	if h[2] != 8 {
		return nil, fmt.Errorf("Unknown compression method: %v", h[2])
	}

	flg := h[3]

	if flg&4 != 0 {
		h := make([]byte, 2)
		n, err := io.ReadFull(dz.fp, h)
		if err != nil {
			return nil, err
		}
		p += n

		xlen := int(h[0]) + 256*int(h[1])
		h = make([]byte, xlen)
		n, err = io.ReadFull(dz.fp, h)
		if err != nil {
			return nil, err
		}
		p += n

		for q := 0; q < len(h); {
			si1 := h[q]
			si2 := h[q+1]
			ln := int(h[q+2]) + 256*int(h[q+3])

			if si1 == 'R' && si2 == 'A' {
				metadata = h[q+4 : q+4+ln]
			}

			q += 4 + ln
		}

	}

	// skip file name (8), file comment (16)
	for _, f := range []byte{8, 16} {
		if flg&f != 0 {
			h := make([]byte, 1)
			for {
				n, err := io.ReadFull(dz.fp, h)
				if err != nil {
					return nil, err
				}
				p += n
				if h[0] == 0 {
					break
				}
			}
		}
	}

	if flg&2 != 0 {
		h := make([]byte, 2)
		n, err := io.ReadFull(dz.fp, h)
		if err != nil {
			return nil, err
		}
		p += n
	}

	if len(metadata) < 6 {
		return nil, fmt.Errorf("Missing dictzip metadata")
	}

	version := int(metadata[0]) + 256*int(metadata[1])

	if version != 1 {
		return nil, fmt.Errorf("Unknown dictzip version: %v", version)
	}

	dz.blocksize = int64(metadata[2]) + 256*int64(metadata[3])
	blockcnt := int(metadata[4]) + 256*int(metadata[5])

	dz.offsets = make([]int64, blockcnt+1)
	dz.offsets[0] = int64(p)
	for i := 0; i < blockcnt; i++ {
		dz.offsets[i+1] = dz.offsets[i] + int64(metadata[6+2*i]) + 256*int64(metadata[7+2*i])
	}

	return dz, nil

}

func (dz *reader) get(start, size int64) ([]byte, error) {

	if size == 0 {
		return []byte{}, nil
	}

	if start < 0 || size < 0 {
		return nil, fmt.Errorf("Negative start or size")
	}

	if int(start/dz.blocksize) >= len(dz.offsets) {
		return nil, fmt.Errorf("Start passed end of archive")
	}

	start1 := dz.blocksize * (start / dz.blocksize)
	size1 := size + (start - start1)

	dz.lock.Lock()
	defer dz.lock.Unlock()

	_, err := dz.fp.Seek(dz.offsets[start/dz.blocksize], 0)
	if err != nil {
		return nil, err
	}
	rd := flate.NewReader(dz.fp)

	data := make([]byte, size1)
	_, err = io.ReadFull(rd, data)
	if err != nil {
		return nil, err
	}

	return data[start-start1:], nil
}

// Start and size in base64 notation, such as used by the `dictunzip` program.
func (dz *reader) getB64(start, size string) ([]byte, error) {
	start2, err := decode(start)
	if err != nil {
		return nil, err
	}
	size2, err := decode(size)
	if err != nil {
		return nil, err
	}
	return dz.get(int64(start2), int64(size2))
}

