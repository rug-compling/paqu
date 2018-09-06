
include ../Makefile.cfg

IS=${if ${shell sh -c 'case `go version` in go1.[01234].*) ;; *) echo yes; ;; esac'},=, }

BUILD=${if $(DEFAULT_PAQU_DIR),-ldflags "-X main.DefaultPaquDir$(IS)$(DEFAULT_PAQU_DIR)", }

../../bin/pqconfig: directories_generated.go pqconfig.go
	go build $(BUILD) -o $@ $(OPTS) .

directories_generated.go: ../pqserve/directories.go
	rm -f directories_generated.go
	echo '// GENERATED FILE. DO NOT EDIT.' > directories_generated.go
	echo >> directories_generated.go
	cat ../pqserve/directories.go >> directories_generated.go
	chmod 444 directories_generated.go
