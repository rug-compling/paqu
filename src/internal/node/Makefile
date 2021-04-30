
include ../../Makefile.cfg

all: attribs.go

attribs.go: attribs.txt mkattribs.go
	$(GO) run mkattribs.go attribs.txt > attribs.tmp
	mv attribs.tmp attribs.go
