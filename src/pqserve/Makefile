
include ../Makefile.cfg

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go


all: attributes_generated.go spod2xpath.go ../../bin/pqserve


spod2xpath.go: spoddata.go static/macros.txt spodmake.go
	$(GO) run -tags spodmake spodmake.go spoddata.go static/macros.txt > spod2xpath.tmp && mv spod2xpath.tmp spod2xpath.go || touch spod2xpath.go
	gofmt -w spod2xpath.go

attributes_generated.go: attributes.txt tools/mkattribs.go
	$(GO) run tools/mkattribs.go attributes.txt > attributes_generated.go

../../bin/pqserve: *.go static/*.* static/Makefile attributes.txt tools/mkattribs.go $(INTERNAL)
	make -C static
	$(GO) build -o $@ $(TAGS) $(OPTS) .

../internal/dir/default.go:
	make -C ../internal/dir
