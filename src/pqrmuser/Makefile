
include ../Makefile.cfg

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go

../../bin/pqrmuser: *.go $(INTERNAL)
	$(GO) build $(BUILD) -o $@ $(OPTS) .

../internal/dir/default.go:
	make -C ../internal/dir
