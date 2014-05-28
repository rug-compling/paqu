
ifeq ($(dbxml), true)
  DBXML = ok
else ifeq ($(dbxml), false)
  DBXML = ok
else
  DBXML = fout
endif

all: check
	make -C alpino
	make -C pqbuild
	make -C pqbuild1
	make -C pqserve
	make -C pqserve1
	make -C pqstatus
	make -C rmuser

check:
	@if [ $(DBXML) != ok ]; then echo "\nGebruik:\n\n  dbxml=true make\n  dbxml=false make\n"; exit 1; fi
