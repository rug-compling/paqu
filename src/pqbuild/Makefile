
include ../Makefile.cfg

SRC = dbxml_false.go dbxml_true.go ranges.go pqbuild.go

IS=${if ${shell sh -c 'case `go version` in go1.[01234].*) ;; *) echo yes; ;; esac'},=, }

BUILD=${if $(DEFAULT_PAQU_DIR),-ldflags "-X main.DefaultPaquDir$(IS)$(DEFAULT_PAQU_DIR)", }

../../bin/pqbuild: $(SRC)
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build $(BUILD) -o $@ $(TAGS) $(OPTS) .

ranges.go: ../pqserve/ranges.go
	rm -f ranges.go
	echo '// GENERATED FILE. DO NOT EDIT.' > ranges.go
	echo >> ranges.go
	cat ../pqserve/ranges.go >> ranges.go
	chmod 444 ranges.go
