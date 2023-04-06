all: index.rdf

index.rdf: index.ttl
	rapper -q -i turtle -o rdfxml-abbrev $< > $@

index.ttl: index.json
	./jsonld2ttl.py $< > $@

index.json: ../codemeta.json
	cp $< $@
