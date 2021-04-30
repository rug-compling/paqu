
include ../../../Makefile.cfg

DEPS = \
	busy.gif \
	clariah.png \
	clarinnl.png \
	expanded.png \
	favicon.ico \
	foliahelp0.html \
	foliahelp1.html \
	foliahelp2.html \
	foliahelp3.html \
	foliahelp4.html \
	foliahelp5.html \
	info.html \
	jquery.js \
	jquery.textcomplete.js \
	macros.txt \
	paqu.css \
	paqu.png \
	relhead.png \
	relnone.png \
	relother.png \
	robots.txt \
	tooltip.css \
	tooltip.js \


../file.go: Makefile mkstatic.go $(DEPS)
	$(GO) run mkstatic.go $(DEPS) > $@
