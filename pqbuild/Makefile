
ifeq ($(dbxml), true)
  DBXML = ok
  TAGS =
else ifeq ($(dbxml), false)
  DBXML = ok
  TAGS = -tags nodbxml
else
  DBXML = fout
endif

all: check pqbuild

pqbuild: *.go
	LIBRARY_PATH=/net/aps/64/lib go build $(TAGS) .

check:
	@if [ $(DBXML) != ok ]; then echo "\nGebruik:\n\n  dbxml=true make\n  dbxml=false make\n"; exit 1; fi
