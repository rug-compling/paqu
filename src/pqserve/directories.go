package main

import (
	"os"
	"path/filepath"
)

var (
	DefaultPaquDir string
	paqudatadir    string
	paquconfigdir  string
)

func init() {

	var d string

	d = os.Getenv("PAQU")
	if d != "" {
		paqudatadir = d
		paquconfigdir = d
		return
	}

	if DefaultPaquDir != "" {
		paqudatadir = DefaultPaquDir
		paquconfigdir = DefaultPaquDir
		return
	}

	home := os.Getenv("HOME")

	d = os.Getenv("XDG_DATA_HOME")
	if d == "" {
		if home == "" {
			d = filepath.Join("/", "usr", "local", "share")
		} else {
			d = filepath.Join(home, ".local", "share")
		}
	}
	paqudatadir = filepath.Join(d, "paqu")

	d = os.Getenv("XDG_CONFIG_HOME")
	if d == "" {
		if home == "" {
			d = filepath.Join("/", "etc", "xdg")
		} else {
			d = filepath.Join(home, ".config")
		}
	}
	paquconfigdir = filepath.Join(d, "paqu")

}
