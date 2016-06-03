package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	"github.com/pebbe/util"

	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Contact string `toml:"contact"`

	Url  string `toml:"url"`
	Port int    `toml:"port"`

	Default string `toml:"default"`

	Mailfrom string `toml:"mailfrom"`
	Smtpserv string `toml:"smtpserv"`
	Smtpuser string `toml:"smtpuser"`
	Smtppass string `toml:"smtppass"`

	Login  string `toml:"login"`
	Prefix string `toml:"prefix"`

	Maxjob int  `toml:"maxjob"`
	Maxwrd int  `toml:"maxwrd"`
	Maxdup int  `toml:"maxdup"`
	Dact   bool `toml:"dact"`

	Sh      string `toml:"sh"`
	Path    string `toml:"path"`
	Alpino  string `toml:"alpino"`
	Timeout int    `toml:"timeout"`

	Secret string `toml:"secret"`

	Https     bool `toml:"https"`
	Httpdual  bool `toml:"httpdual"`
	Remote    bool `toml:"remote"`
	Forwarded bool `toml:"forwarded"`

	Querytimeout int `toml:"querytimeout"` // in secondes

	Loginurl string `toml:"loginurl"`

	Foliadays int `toml:"foliadays"`

	View   []ViewType   `toml:"view"`
	Access []AccessType `toml:"access"`
}

type ViewType struct {
	Allow bool
	Addr  []string
}

type AccessType struct {
	Allow bool
	Mail  []string
}

var (
	DefaultPaquDir string
)

//. Main

func main() {

	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		if DefaultPaquDir != "" {
			paqudir = DefaultPaquDir
		} else {
			paqudir = filepath.Join(os.Getenv("HOME"), ".paqu")
		}
	}

	cfg := Config{}
	md, err := toml.DecodeFile(filepath.Join(paqudir, "setup.toml"), &cfg)
	util.CheckErr(err)

	for _, un := range md.Undecoded() {
		fmt.Println("UNDEFINED:", un)
	}

	e := toml.NewEncoder(os.Stdout)
	e.Encode(cfg)

}
