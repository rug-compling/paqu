
include ../Makefile.cfg

ifeq ($(shell $(VARS) go list github.com/pebbe/dbxml), github.com/pebbe/dbxml)
 TAGS =
else
 TAGS = -tags nodbxml
endif

pqserve: *.go static/*.*
	make -C static
	$(VARS) go build $(TAGS) $(OPTS) .
