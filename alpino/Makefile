
% : %.go
	go build $<

all: alpino prepare testdoc-ascii.txt testdoc-latin.txt testdoc-utf16.txt testdoc-windows.txt

testdoc-ascii.txt: testdoc-utf8.txt
	any2any u8 us-ascii < $< > $@

testdoc-latin.txt: testdoc-utf8.txt
	any2any u8 l1 < $< > $@

testdoc-utf16.txt: testdoc-utf8.txt
	any2any u8 utf16 < $< > $@

testdoc-windows.txt: testdoc-utf8.txt
	any2any u8 windows-1252 < $< > $@

prepare: prepare.l
	flex prepare.l
	gcc -s -Wall -o prepare prepare.c -lfl
	rm prepare.c
