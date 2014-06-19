
# TODO: gebruik van any2any verwijderen

include ../Makefile.cfg

all: ../../bin/alpino ../../bin/prepare testdoc-ascii.txt testdoc-latin.txt testdoc-utf16.txt testdoc-windows.txt

../../bin/alpino: *.go
	go build -o $@ $(OPTS) .

../../bin/prepare: prepare.l
	flex prepare.l
	gcc -s -Wall -o $@ prepare.c -lfl
	rm prepare.c

testdoc-ascii.txt: testdoc-utf8.txt
	any2any u8 us-ascii < $< > $@

testdoc-latin.txt: testdoc-utf8.txt
	any2any u8 l1 < $< > $@

testdoc-utf16.txt: testdoc-utf8.txt
	any2any u8 utf16 < $< > $@

testdoc-windows.txt: testdoc-utf8.txt
	any2any u8 windows-1252 < $< > $@

