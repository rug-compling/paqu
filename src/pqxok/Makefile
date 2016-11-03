
include ../Makefile.cfg

../../bin/pqxok: *.go
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build -o $@ $(TAGS) $(OPTS) .
