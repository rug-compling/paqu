
include ../Makefile.cfg

all: ../../bin/alpino ../../bin/prepare

../../bin/alpino: *.go
	go build -o $@ $(OPTS) .

../../bin/prepare: prepare.l
	flex prepare.l
	gcc -s -Wall -o $@ prepare.c -lfl
	rm prepare.c
