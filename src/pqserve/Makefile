
include ../Makefile.cfg

ifeq ($(shell $(VARS) go list github.com/pebbe/dbxml), github.com/pebbe/dbxml)
 TAGS =
else
 TAGS = -tags nodbxml
endif

../../bin/pqserve: *.go static/*.* static/Makefile
	make -C static
	$(VARS) go build -o $@ $(TAGS) $(OPTS) .
