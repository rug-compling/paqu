
include ../Makefile.cfg

../../bin/pqserve: *.go static/*.* static/Makefile
	make -C static
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build -o $@ $(TAGS) $(OPTS) .
