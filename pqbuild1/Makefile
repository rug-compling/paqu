
include ../Makefile.cfg

ifeq ($(shell go list github.com/pebbe/dbxml), github.com/pebbe/dbxml)
 TAGS =
else
 TAGS = -tags nodbxml
endif

pqbuild1: *.go
	$(VARS) go build $(TAGS) $(OPTS) .
