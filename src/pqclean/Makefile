
include ../Makefile.cfg

IS=${if ${shell sh -c 'case `go version` in go1.[01234].*) ;; *) echo yes; ;; esac'},=, }

BUILD=${if $(DEFAULT_PAQU_DIR),-ldflags "-X main.DefaultPaquDir$(IS)$(DEFAULT_PAQU_DIR)", }

../../bin/pqclean: *.go
	go build $(BUILD) -o $@ $(OPTS) .
