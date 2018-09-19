
include ../Makefile.cfg

SRC = attributes_generated.go udep_generated.go pqudep.go

../../bin/pqudep: $(SRC)
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build -o $@ $(TAGS) $(OPTS) .

attributes_generated.go: ../pqserve/attributes_generated.go mkAttribs
	./mkAttribs < $< > $@

udep_generated.go: udep.xq
	rm -f $@
	echo '// GENERATED FILE. DO NOT EDIT.' > $@
	echo >> $@
	echo package main >> $@
	echo >> $@
	echo 'var udep = `' >> $@
	cat $< >> $@
	echo '`' >> $@
