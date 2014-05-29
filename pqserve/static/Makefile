
include ../../Makefile.cfg

DEPS = \
	busy.gif \
	clarinnl.png \
	favicon.ico \
	info.html \
	jquery.js \
	paqu.css \
	paqu.png \
	rename.html \
	robots.txt \
	share.html \
	tooltip.css \
	tooltip.js \


../static_data.go: mkstatic $(DEPS)
	./mkstatic $(DEPS) > $@

mkstatic: mkstatic.go
	$(VARS) go build $(OPTS) $^
