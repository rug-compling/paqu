
include ../Makefile.cfg

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go \
	../internal/file/file.go \
	../internal/ranges/ranges.go

SRC = dbxml_false.go dbxml_true.go pqbuild.go spod.go spoddata_generated.go xpath_generated.go spod2xpath_generated.go

../../bin/pqbuild: $(SRC) $(INTERNAL)
	$(GO) build -o $@ $(TAGS) $(OPTS) .

spoddata_generated.go: ../pqserve/spoddata.go
	rm -f spoddata_generated.go
	echo '// GENERATED FILE. DO NOT EDIT.' > spoddata_generated.go
	echo >> spoddata_generated.go
	cat ../pqserve/spoddata.go >> spoddata_generated.go
	chmod 444 spoddata_generated.go

xpath_generated.go: ../pqserve/xpath.go
	rm -f xpath_generated.go
	echo '// GENERATED FILE. DO NOT EDIT.' > xpath_generated.go
	echo >> xpath_generated.go
	cat ../pqserve/xpath.go >> xpath_generated.go
	chmod 444 xpath_generated.go

spod2xpath_generated.go: ../pqserve/spod2xpath.go
	rm -f spod2xpath_generated.go
	cp ../pqserve/spod2xpath.go spod2xpath_generated.go
	chmod 444 spod2xpath_generated.go

../internal/dir/default.go:
	make -C ../internal/dir

../internal/file/file.go: ../internal/file/data/*
	make -C ../internal/file/data
