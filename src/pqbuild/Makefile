
include ../Makefile.cfg

SRC = dbxml_false.go dbxml_true.go ranges_generated.go directories_generated.go pqbuild.go spod.go spoddata_generated.go xpath_generated.go spod2xpath_generated.go static_generated.go

IS=${if ${shell sh -c 'case `go version` in go1.[01234].*) ;; *) echo yes; ;; esac'},=, }

BUILD=${if $(DEFAULT_PAQU_DIR),-ldflags "-X main.DefaultPaquDir$(IS)$(DEFAULT_PAQU_DIR)", }

../../bin/pqbuild: $(SRC)
	$(GO) build $(BUILD) -o $@ $(TAGS) $(OPTS) .

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

spoddata_generated.go: ../pqserve/spoddata.go
	rm -f spoddata_generated.go
	echo '// GENERATED FILE. DO NOT EDIT.' > spoddata_generated.go
	echo >> spoddata_generated.go
	cat ../pqserve/spoddata.go >> spoddata_generated.go
	chmod 444 spoddata_generated.go

xpath_generated.go: ../pqserve/xpath.go
	rm -f xpath_generated.go
	echo '// GENERATED FILE. DO NOT EDIT.' > xpath_generated.go
	echo >> xpath_generated.go
	cat ../pqserve/xpath.go >> xpath_generated.go
	chmod 444 xpath_generated.go

spod2xpath_generated.go: ../pqserve/spod2xpath.go
	rm -f spod2xpath_generated.go
	cp ../pqserve/spod2xpath.go spod2xpath_generated.go
	chmod 444 spod2xpath_generated.go

static_generated.go: ../pqserve/static_generated.go
	rm -f static_generated.go
	cp ../pqserve/static_generated.go static_generated.go
	chmod 444 static_generated.go
