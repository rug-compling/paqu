
include ../Makefile.cfg

SRC = dbxml_false.go dbxml_true.go ranges_generated.go directories_generated.go pqbuild.go

IS=${if ${shell sh -c 'case `go version` in go1.[01234].*) ;; *) echo yes; ;; esac'},=, }

BUILD=${if $(DEFAULT_PAQU_DIR),-ldflags "-X main.DefaultPaquDir$(IS)$(DEFAULT_PAQU_DIR)", }

../../bin/pqbuild: $(SRC)
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go build $(BUILD) -o $@ $(TAGS) $(OPTS) .

ranges_generated.go: ../pqserve/ranges.go
	rm -f ranges_generated.go
	echo '// GENERATED FILE. DO NOT EDIT.' > ranges_generated.go
	echo >> ranges_generated.go
	cat ../pqserve/ranges.go >> ranges_generated.go
	chmod 444 ranges_generated.go

directories_generated.go: ../pqserve/directories.go
	rm -f directories_generated.go
	echo '// GENERATED FILE. DO NOT EDIT.' > directories_generated.go
	echo >> directories_generated.go
	cat ../pqserve/directories.go >> directories_generated.go
	chmod 444 directories_generated.go
