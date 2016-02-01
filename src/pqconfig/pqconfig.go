package main

//. Imports

import (
	"github.com/BurntSushi/toml"
	"github.com/pebbe/util"

	"os"
	"path"
)

type Config struct {
	Contact string

	Port int
	Url  string

	Default string

	Mailfrom string
	Smtpserv string
	Smtpuser string
	Smtppass string

	Login  string
	Prefix string

	Maxjob int
	Maxwrd int
	Maxdup int
	Dact   bool

	Secret string

	Sh       string
	Path     string
	Alpino   string
	Alpino15 bool
	Timeout  int

	Https     bool
	Httpdual  bool
	Remote    bool
	Forwarded bool

	Querytimeout int // in secondes

	Directories []string

	View   []ViewType
	Access []AccessType
}

type ViewType struct {
	Allow bool
	Addr  []string
}

type AccessType struct {
	Allow bool
	Mail  []string
}

//. Main

func main() {

	paqudir := os.Getenv("PAQU")
	if paqudir == "" {
		paqudir = path.Join(os.Getenv("HOME"), ".paqu")
	}

	cfg := Config{}
	_, err := toml.DecodeFile(path.Join(paqudir, "setup.toml"), &cfg)
	util.CheckErr(err)

	e := toml.NewEncoder(os.Stdout)
	e.Encode(cfg)
}
