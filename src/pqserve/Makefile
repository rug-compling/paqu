
include ../Makefile.cfg

all: attributes_generated.go ../../bin/pqserve

attributes_generated.go: attributes.txt tools/mkattribs.go
	go run tools/mkattribs.go attributes.txt > attributes_generated.go

../../bin/pqserve: *.go static/*.* static/Makefile attributes.txt tools/mkattribs.go
	make -C static
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build -o $@ $(TAGS) $(OPTS) .
