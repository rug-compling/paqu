
include ../Makefile.cfg

INTERNAL = \
	../internal/node/attribs.go \
	../internal/node/node.go

../../bin/pqdactx: *.go ../go.mod $(INTERNAL)
	$(GO) build -o $@ $(TAGS) $(OPTS) .

../internal/node/attribs.go: ../go.mod ../internal/node/mkattribs.go
	make -C ../internal/node
