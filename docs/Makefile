all: index.rdf

index.rdf: index.ttl
	rapper -q -i turtle -o rdfxml-abbrev $< > $@.tmp
	mv $@.tmp $@

index.ttl: index.jsonld
	./jsonld2ttl.py $< > $@.tmp
	mv $@.tmp $@

index.jsonld: ../codemeta.json
	cp $< $@
