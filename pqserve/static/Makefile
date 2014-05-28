
include ../../Makefile.cfg

DEPS = \
	busy.gif \
	clarinnl.png \
	favicon.ico \
	info.html \
	jquery.js \
	rename.html \
	robots.txt \
	share.html \
	tooltip.css \
	tooltip.js \
	wordrel.css \



../static_data.go: mkstatic $(DEPS)
	./mkstatic $(DEPS) > $@

mkstatic: mkstatic.go
	$(VARS) go build $(OPTS) $^
