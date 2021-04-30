package dir

import (
	"os"
	"path/filepath"
)

var (
	Data   string
	Config string
)

func init() {

	var d string

	d = os.Getenv("PAQU")
	if d != "" {
		Data = d
		Config = d
		return
	}

	if DefaultDir != "" {
		Data = DefaultDir
		Config = DefaultDir
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
	Data = filepath.Join(d, "paqu")

	d = os.Getenv("XDG_CONFIG_HOME")
	if d == "" {
		if home == "" {
			d = filepath.Join("/", "etc", "xdg")
		} else {
			d = filepath.Join(home, ".config")
		}
	}
	Config = filepath.Join(d, "paqu")

}
