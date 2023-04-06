all: index.rdf

index.rdf: index.ttl
	rapper -q -i turtle -o rdfxml-abbrev $< > $@

index.ttl: index.jsonld
	./jsonld2ttl.py $< > $@

index.jsonld: ../codemeta.json
	cp $< $@
