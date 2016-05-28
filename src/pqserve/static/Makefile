
include ../../Makefile.cfg

DEPS = \
	busy.gif \
	clariah.png \
	clarinnl.png \
	favicon.ico \
	foliahelp1.html \
	foliahelp2.html \
	foliahelp3.html \
	foliahelp4.html \
	foliahelp5.html \
	info.html \
	jquery.js \
	jquery.textcomplete.js \
	paqu.css \
	paqu.png \
	relhead.png \
	relnone.png \
	relother.png \
	robots.txt \
	tooltip.css \
	tooltip.js \


../static_generated.go: Makefile mkstatic.go $(DEPS)
	go run mkstatic.go $(DEPS) > $@
