
include ../Makefile.cfg

IS=${if ${shell sh -c 'case `go version` in go1.[01234].*) ;; *) echo yes; ;; esac'},=, }

BUILD=${if $(DEFAULT_PAQU_DIR),-ldflags "-X main.DefaultPaquDir$(IS)$(DEFAULT_PAQU_DIR)", }

all: attributes_generated.go ../../bin/pqserve

attributes_generated.go: attributes.txt tools/mkattribs.go
	go run tools/mkattribs.go attributes.txt > attributes_generated.go

../../bin/pqserve: *.go static/*.* static/Makefile attributes.txt tools/mkattribs.go
	make -C static
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build $(BUILD) -o $@ $(TAGS) $(OPTS) .
