
include ../Makefile.cfg

all: ../../bin/pqudep

../../bin/pqudep: *.go
	$(GO) build -o $@ $(TAGS) $(OPTS) .
