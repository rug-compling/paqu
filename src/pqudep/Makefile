
include ../Makefile.cfg

SRC = attributes_generated.go pqudep.go

all: ../../bin/pqudep

../../bin/pqudep: $(SRC)
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build -o $@ $(TAGS) $(OPTS) .

attributes_generated.go: ../pqserve/attributes_generated.go mkAttribs
	./mkAttribs < $< > $@

