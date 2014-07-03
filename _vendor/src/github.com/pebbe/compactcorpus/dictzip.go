package compactcorpus

// dit is gecopieerd uit dictzip

//. Imports

import (
	"compress/flate"
	"fmt"
	"io"
	"sync"
)

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
	n, err := readfull(dz.fp, h)
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
		n, err := readfull(dz.fp, h)
		if err != nil {
			return nil, err
		}
		p += n

		xlen := int(h[0]) + 256*int(h[1])
		h = make([]byte, xlen)
		n, err = readfull(dz.fp, h)
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
				n, err := readfull(dz.fp, h)
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
		n, err := readfull(dz.fp, h)
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
	_, err = readfull(rd, data)
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

//. Helper function

func readfull(fp io.Reader, buf []byte) (int, error) {
	ln := len(buf)
	for p := 0; p < ln; {
		n, err := fp.Read(buf[p:])
		p += n
		if err != nil {
			if err != io.EOF || p < ln {
				return p, err
			} else {
				return p, nil
			}
		}
	}
	return ln, nil
}

