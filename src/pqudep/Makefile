
include ../Makefile.cfg

SRC = attributes_generated.go pqudep.go

all: ../../bin/pqudep

../../bin/pqudep: $(SRC)
	$(GO) build -o $@ $(TAGS) $(OPTS) .

attributes_generated.go: ../pqserve/attributes_generated.go mkAttribs
	./mkAttribs < $< > $@

