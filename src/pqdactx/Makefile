
include ../Makefile.cfg

INTERNAL = \
	../internal/node/attribs.go \
	../internal/node/node.go

../../bin/pqdactx: *.go $(INTERNAL)
	$(GO) build -o $@ $(TAGS) $(OPTS) .

../internal/node/attribs.go: ../internal/node/attribs.txt ../internal/node/mkattribs.go
	make -C ../internal/node
