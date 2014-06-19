
include ../../Makefile.cfg

DEPS = \
	busy.gif \
	clarinnl.png \
	favicon.ico \
	info.html \
	jquery.js \
	paqu.css \
	paqu.png \
	robots.txt \
	tooltip.css \
	tooltip.js \


../static_data.go: Makefile mkstatic $(DEPS)
	./mkstatic $(DEPS) > $@

mkstatic: mkstatic.go
	go build $(OPTS) $^
