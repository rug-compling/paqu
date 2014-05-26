all:
	make -C alpino
	make -C pqbuild
	make -C pqbuild1
	make -C pqserve
	make -C pqserve1
	make -C pqstatus
	make -C rmuser
