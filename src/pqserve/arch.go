package main

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

var (
	ArchClosed = errors.New("Archive is closed")
)

type arch struct {
	opened bool

	isZip bool
	zi    int
	zr    *zip.ReadCloser

	isTar   bool
	tfp     *os.File
	tr      *tar.Reader
	theader *tar.Header
	terr    error
	tstart  bool
}

func NewArchReader(filename string) (*arch, error) {
	a := &arch{}
	var err error
	a.zr, err = zip.OpenReader(filename)
	if err == nil {
		a.isZip = true
		a.opened = true
		return a, nil
	}

	a.tfp, err = os.Open(filename)
	if err != nil {
		return a, err
	}
	a.tr = tar.NewReader(a.tfp)
	a.theader, a.terr = a.tr.Next()
	if a.terr != nil && a.terr != io.EOF {
		return a, a.terr
	}
	a.isTar = true
	a.opened = true
	return a, nil

}

func (a *arch) Next() error {
	if !a.opened {
		return ArchClosed
	}
	if a.isZip {
		for {
			if a.zi == len(a.zr.File) {
				if a.opened {
					a.opened = false
					a.zr.Close()
				}
				return io.EOF
			}
			a.zi++
			if !a.zr.File[a.zi-1].FileInfo().IsDir() {
				return nil
			}
		}
	}

	// bij de setup is a.tr.Next() al een keer gedaan, dus vandaar deze vreemde constructie
	for {
		if a.tstart {
			a.theader, a.terr = a.tr.Next()
		} else {
			a.tstart = true
		}
		if a.terr != nil || !a.theader.FileInfo().IsDir() {
			break
		}
	}

	if a.terr == io.EOF {
		a.opened = false
		a.tfp.Close()
	}

	return a.terr
}

func (a *arch) ReadN(n uint) ([]byte, error) {
	if !a.opened {
		return []byte{}, ArchClosed
	}

	b := make([]byte, n)

	if a.isZip {
		rc, err := a.zr.File[a.zi-1].Open()
		if err != nil {
			return []byte{}, err
		}
		_, err = io.ReadFull(rc, b)
		rc.Close()
		return b, err
	}

	_, err := io.ReadFull(a.tr, b)
	return b, err
}

func (a *arch) Read() ([]byte, error) {
	if !a.opened {
		return []byte{}, ArchClosed
	}

	if a.isZip {
		rc, err := a.zr.File[a.zi-1].Open()
		if err != nil {
			return []byte{}, err
		}
		b, err := ioutil.ReadAll(rc)
		rc.Close()
		return b, err
	}

	b, err := ioutil.ReadAll(a.tr)
	return b, err
}

func (a *arch) Copy(fp io.Writer) error {
	if !a.opened {
		return ArchClosed
	}

	if a.isZip {
		rc, err := a.zr.File[a.zi-1].Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(fp, rc)
		rc.Close()
		return err
	}

	_, err := io.Copy(fp, a.tr)
	return err
}

func (a *arch) Name() string {
	if !a.opened {
		return ""
	}

	if a.isZip {
		return a.zr.File[a.zi-1].Name
	}

	return a.theader.Name
}

func (a *arch) Close() {
	if a.opened {
		a.opened = false
		if a.isZip {
			a.zr.Close()
		} else {
			a.tfp.Close()
		}
	}
}

func (a *arch) IsZip() bool {
	return a.isZip
}

func (a *arch) IsTar() bool {
	return a.isTar
}
