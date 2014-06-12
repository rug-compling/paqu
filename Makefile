
all:
	make -C alpino
	make -C pqbuild
	make -C pqserve
	make -C pqstatus
	make -C rmuser
