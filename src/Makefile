
include Makefile.cfg

info:
	@echo
	@echo Kopieer Makefile.cfg.example naar Makefile.cfg
	@echo Pas Makefile.cfg aan aan de lokale situatie
	@echo
	@echo Kopieer internal/dir/default.go.example naar internal/dir/default.go
	@echo Pas internal/dir/default.go aan aan de lokale situatie
	@echo
	@echo Run daarna: make all
	@echo

all:
	mkdir -p ../bin
	make -C pqserve
	make -C pqalpino
	make -C pqbugtest
	make -C pqbuild
	make -C pqclean
	make -C pqconfig
	make -C pqconllu
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
	make -C pqudupgrade
	make -C pqudep

# pqudep moet als laatst, voor het geval dat libxqilla
# niet aanwezig is en dit faalt.
