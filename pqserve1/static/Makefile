
include ../../Makefile.cfg

DEPS = \
	clarinnl.png \
	favicon.ico \
	info.html \
	jquery.js \
	paqu.css \
	paqu.png \
	tooltip.css \
	tooltip.js \


../static_data.go: mkstatic $(DEPS)
	./mkstatic $(DEPS) > $@

mkstatic: mkstatic.go
	$(VARS) go build $(OPTS) $^
