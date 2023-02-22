
include ../../Makefile.cfg

all: attribs.go

attribs.go: mkattribs.go ../../go.mod
	$(GO) run mkattribs.go > attribs.tmp
	mv attribs.tmp attribs.go
