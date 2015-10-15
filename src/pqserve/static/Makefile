
include ../../Makefile.cfg

DEPS = \
	busy.gif \
	clariah.png \
	clarinnl.png \
	favicon.ico \
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


../static_data.go: Makefile mkstatic $(DEPS)
	./mkstatic $(DEPS) > $@

mkstatic: mkstatic.go
	go build $(OPTS) $^
