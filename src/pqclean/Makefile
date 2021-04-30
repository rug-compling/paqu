
include ../Makefile.cfg

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go

../../bin/pqclean: *.go $(INTERNAL)
	$(GO) build -o $@ $(OPTS) .

../internal/dir/default.go:
	make -C ../internal/dir
