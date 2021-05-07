
include ../Makefile.cfg

INTERNAL = \
	../internal/dir/default.go \
	../internal/dir/dir.go

../../bin/pqstatus: *.go ../go.mod $(INTERNAL)
	$(GO) build $(BUILD) -o $@ $(OPTS) .

../internal/dir/default.go:
	make -C ../internal/dir
