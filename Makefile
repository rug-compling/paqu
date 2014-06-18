
all:
	make -C alpino
	make -C dbinit
	make -C pqbuild
	make -C pqserve
	make -C pqstatus
	make -C rmuser
	make -C setquota
