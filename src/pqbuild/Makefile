
include ../Makefile.cfg

../../bin/pqbuild: *.go
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build -o $@ $(TAGS) $(OPTS) .
