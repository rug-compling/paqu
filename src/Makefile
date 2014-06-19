
info:
	@echo
	@echo Kopieer Makefile.cfg.example naar Makefile.cfg
	@echo Pas Makefile.cfg aan aan de lokale situatie
	@echo Run daarna: make all
	@echo

all:
	mkdir -p ../bin
	make -C alpino
	make -C dbinit
	make -C pqbuild
	make -C pqserve
	make -C pqstatus
	make -C rmuser
	make -C setquota
