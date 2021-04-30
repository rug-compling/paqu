
include ../Makefile.cfg

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go \
	../internal/file/file.go \
	../internal/node/attribs.go \
	../internal/node/node.go


all: spod2xpath.go ../../bin/pqserve


spod2xpath.go: spoddata.go ../internal/file/data/macros.txt spodmake.go
	$(GO) run -tags spodmake spodmake.go spoddata.go ../internal/file/data/macros.txt > spod2xpath.tmp && mv spod2xpath.tmp spod2xpath.go || touch spod2xpath.go
	gofmt -w spod2xpath.go

../../bin/pqserve: *.go $(INTERNAL)
	$(GO) build -o $@ $(TAGS) $(OPTS) .

../internal/dir/default.go:
	make -C ../internal/dir

../internal/file/file.go: ../internal/file/data/*
	make -C ../internal/file/data

../internal/node/attribs.go: ../internal/node/attribs.txt ../internal/node/mkattribs.go
	make -C ../internal/node

