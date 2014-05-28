
include ../Makefile.cfg

ifeq ($(shell $(VARS) go list github.com/pebbe/dbxml), github.com/pebbe/dbxml)
 TAGS =
else
 TAGS = -tags nodbxml
endif

pqserve1: *.go static/*.*
	make -C static
	$(VARS) go build $(TAGS) $(OPTS) .
