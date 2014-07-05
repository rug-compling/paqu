
include ../Makefile.cfg

../../bin/pqtexter: pqtexter.l
	flex pqtexter.l
	gcc -s -Wall -o $@ pqtexter.c -lfl
	rm pqtexter.c
