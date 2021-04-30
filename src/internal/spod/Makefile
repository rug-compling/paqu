
include ../../Makefile.cfg

spod2xpath.go: spoddata.go ../file/data/macros.txt spodmake/spodmake.go
	$(GO) run spodmake/spodmake.go ../file/data/macros.txt > spod2xpath.tmp && mv spod2xpath.tmp spod2xpath.go || touch spod2xpath.go
	gofmt -w spod2xpath.go
