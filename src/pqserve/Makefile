
include ../Makefile.cfg

INTERNAL = \
	../internal/node/attribs.go \
	../internal/node/node.go \
	../internal/dir/default.go \
	../internal/dir/dir.go


all: spod2xpath.go ../../bin/pqserve


spod2xpath.go: spoddata.go static/macros.txt spodmake.go
	$(GO) run -tags spodmake spodmake.go spoddata.go static/macros.txt > spod2xpath.tmp && mv spod2xpath.tmp spod2xpath.go || touch spod2xpath.go
	gofmt -w spod2xpath.go

../../bin/pqserve: *.go static/*.* static/Makefile $(INTERNAL)
	make -C static
	$(GO) build -o $@ $(TAGS) $(OPTS) .

../internal/node/attribs.go: ../internal/node/attribs.txt ../internal/node/mkattribs.go
	make -C ../internal/node

../internal/dir/default.go:
	make -C ../internal/dir
