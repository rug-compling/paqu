
include Makefile.cfg

info:
	@echo
	@echo Kopieer Makefile.cfg.example naar Makefile.cfg
	@echo Pas Makefile.cfg aan aan de lokale situatie
	@echo Run daarna: make all
	@echo

all:
ifneq "$(findstring paqu/_vendor, $(GOPATH))" ""
	go install github.com/BurntSushi/toml
	go install github.com/beevik/etree
	go install github.com/dchest/authcookie
	go install github.com/go-sql-driver/mysql
	go install github.com/pebbe/compactcorpus
	CGO_LDFLAGS="-Wl,-rpath=$(LIBRARY_PATH)" go install github.com/pebbe/dbxml
	go install github.com/pebbe/util
	go install github.com/robertkrimen/otto
endif
	mkdir -p ../bin
	make -C pqserve
	make -C pqalpino
	make -C pqbugtest
	make -C pqbuild
	make -C pqclean
	make -C pqconfig
	make -C pqdactx
	make -C pqfolia
	make -C pqinit
	make -C pqrmcorpus
	make -C pqrmuser
	make -C pqsetquota
	make -C pqstatus
	make -C pqtexter
	make -C pqupgrade
	make -C pqxok
	make -C pqudep
	make -C pqudupgrade

# pqserve moet als eerste gecompileerd worden omdat
# sommige andere programma's ervan afhankelijk zijn.

# udep moet als laatst, voor het geval dat libxqilla
# niet aanwezig is en dit faalt.
