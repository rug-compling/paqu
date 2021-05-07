
include ../Makefile.cfg

all: ../../bin/pqudep

../../bin/pqudep: *.go ../go.mod
	$(GO) build -o $@ $(TAGS) $(OPTS) .
