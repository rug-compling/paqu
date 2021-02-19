
include ../Makefile.cfg

SRC = attributes_generated.go pqdactx_true.go pqdactx_false.go

../../bin/pqdactx: $(SRC)
	$(GO) build $(BUILD) -o $@ $(TAGS) $(OPTS) .

attributes_generated.go: ../pqserve/attributes_generated.go
	echo '// +build !nodbxml' > attributes_generated.go
	echo >> attributes_generated.go
	cat ../pqserve/attributes_generated.go >> attributes_generated.go
