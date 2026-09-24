#!/bin/bash

. env.sh

P=wordrel
D=/home/peter/paqu-data/data

for db in $(sqlite3 paqu-local.db 'select id from info')
do
    echo ================================================================
    echo $db
    rm -f $D/$db/data.duckdb
    echo "LOAD MYSQL;
ATTACH 'host=localhost' AS db1 (TYPE MYSQL, READ_ONLY);

CREATE TABLE arch   AS SELECT * FROM db1.${P}_c_${db}_arch;
CREATE UNIQUE INDEX arch_id ON arch(id);

CREATE TABLE file   AS SELECT * FROM db1.${P}_c_${db}_file;
CREATE UNIQUE INDEX file_id ON file(id);

CREATE TABLE deprel AS SELECT * FROM db1.${P}_c_${db}_deprel;
CREATE UNIQUE INDEX deprel_idd ON deprel(idd);
CREATE INDEX deprel_word    ON deprel(word);
CREATE INDEX deprel_lemma   ON deprel(lemma);
CREATE INDEX deprel_root    ON deprel(root);
CREATE INDEX deprel_postag  ON deprel(postag);
CREATE INDEX deprel_rel     ON deprel(rel);
CREATE INDEX deprel_hword   ON deprel(hword);
CREATE INDEX deprel_hlemma  ON deprel(hlemma);
CREATE INDEX deprel_hroot   ON deprel(hroot);
CREATE INDEX deprel_hpostag ON deprel(hpostag);
CREATE INDEX deprel_file    ON deprel(file);
CREATE INDEX deprel_arch    ON deprel(arch);

CREATE TABLE meta   AS SELECT * FROM db1.${P}_c_${db}_meta;
CREATE INDEX meta_id ON meta(id);
CREATE INDEX meta_file ON meta(file);
CREATE INDEX meta_arch ON meta(arch);
CREATE INDEX meta_tval ON meta(tval);
CREATE INDEX meta_ival ON meta(ival);
CREATE INDEX meta_fval ON meta(fval);
CREATE INDEX meta_dval ON meta(dval);
CREATE INDEX meta_idx  ON meta(idx);

CREATE TABLE midx   AS SELECT * FROM db1.${P}_c_${db}_midx;
CREATE UNIQUE INDEX midx_id ON midx(id);
CREATE INDEX midx_name ON midx("name");

CREATE TABLE minf   AS SELECT * FROM db1.${P}_c_${db}_minf;
CREATE INDEX minf_id ON minf(id);

CREATE TABLE mval   AS SELECT * FROM db1.${P}_c_${db}_mval;
CREATE INDEX mval_id ON mval(id);
CREATE INDEX mval_idx ON mval(idx);

CREATE TABLE sent   AS SELECT * FROM db1.${P}_c_${db}_sent;
CREATE INDEX sent_file ON sent(file);
CREATE INDEX sent_arch ON sent(arch);
CREATE INDEX sent_lbl  ON sent(lbl);

CREATE TABLE word   AS SELECT * FROM db1.${P}_c_${db}_word;
CREATE UNIQUE INDEX word_word ON word(word);
    " | duckdb -echo $D/$db/data.duckdb
    # read -p continue... X
done

