
info:
	@echo
	@echo Kopieer Makefile.cfg.example naar Makefile.cfg
	@echo Pas Makefile.cfg aan aan de lokale situatie
	@echo Run daarna: make all
	@echo

all:
	mkdir -p ../bin
	make -C pqalpino
	make -C pqbugtest
	make -C pqbuild
	make -C pqclean
	make -C pqconfig
	make -C pqinit
	make -C pqrmcorpus
	make -C pqrmuser
	make -C pqserve
	make -C pqsetquota
	make -C pqstatus
	make -C pqtexter
