
include ../Makefile.cfg

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go

../../bin/pqsetquota: *.go $(INTERNAL)
	$(GO) build $(BUILD) -o $@ $(OPTS) .

../internal/dir/default.go:
	make -C ../internal/dir
