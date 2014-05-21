
DEPS = clarinnl.png favicon.ico info.html jquery.js tooltip.css tooltip.js wordrel.css

../static_data.go: mkstatic $(DEPS)
	./mkstatic $(DEPS) > $@

mkstatic: mkstatic.go
	go build $^
