
include ../Makefile.cfg

../../bin/pqconllu: *.go
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build $(BUILD) -o $@ $(TAGS) $(OPTS) .
