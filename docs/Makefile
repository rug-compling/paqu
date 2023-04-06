all: index.rdf

index.rdf: index.ttl
	rapper -q -i turtle -o rdfxml-abbrev $< > $@

index.ttl: ../codemeta.json
	./jsonld2ttl.py $< > $@
