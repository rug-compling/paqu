
include ../Makefile.cfg

SRC = dbxml_false.go dbxml_true.go ranges.go pqbuild.go

../../bin/pqbuild: $(SRC)
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build -o $@ $(TAGS) $(OPTS) .

ranges.go: ../pqserve/ranges.go
	rm -f ranges.go
	echo '// GENERATED FILE. DO NOT EDIT.' > ranges.go
	echo >> ranges.go
	cat ../pqserve/ranges.go >> ranges.go
	chmod 444 ranges.go
