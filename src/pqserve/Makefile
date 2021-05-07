
include ../Makefile.cfg

# internal/file is nodig voor internal/spod

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go \
	../internal/file/file.go \
	../internal/node/attribs.go \
	../internal/node/node.go \
	../internal/ranges/ranges.go \
	../internal/spod/spod2xpath.go \
	../internal/spod/spoddata.go \
	../internal/spod/util.go \
	../internal/spod/xpath.go


all: ../../bin/pqserve

../../bin/pqserve: *.go ../go.mod $(INTERNAL)
	$(GO) build -o $@ $(TAGS) $(OPTS) .

../internal/dir/default.go:
	make -C ../internal/dir

../internal/file/file.go: ../internal/file/data/*
	make -C ../internal/file/data

../internal/node/attribs.go: ../internal/node/attribs.txt ../internal/node/mkattribs.go
	make -C ../internal/node

../internal/spod/spod2xpath.go: ../internal/spod/spoddata.go ../internal/file/data/macros.txt ../internal/spod/spodmake/spodmake.go
	make -C ../internal/spod
