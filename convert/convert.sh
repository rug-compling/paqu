#!/bin/bash

set -e

. env.sh

if [ ! -f $PAQU/data/paqu.db ]
then

    rm -f $PAQU/data/paqu.db.tmp

    cat <<'EOT' | sqlite3 -bail -echo $PAQU/data/paqu.db.tmp
CREATE TABLE corpora (
  user    TEXT NOT NULL,
  prefix  TEXT NOT NULL,
  enabled INTEGER(1) NOT NULL DEFAULT 1
);
CREATE INDEX corpora_prefix ON corpora(prefix);
CREATE INDEX corpora_user ON corpora(user);
CREATE INDEX corpora_enabled ON corpora(enabled);

CREATE TABLE ignore (
  user   TEXT,
  prefix TEXT
);
CREATE INDEX ignore_user ON ignore(user);
CREATE INDEX ignore_prefix ON ignore(prefix);

CREATE TABLE info (
  id          TEXT NOT NULL,
  description TEXT NOT NULL,
  owner       TEXT NOT NULL DEFAULT 'none',
  status      TEXT NOT NULL DEFAULT 'QUEUING',
  msg         TEXT NOT NULL,
  nline       INTEGER NOT NULL DEFAULT 0,
  nword       INTEGER NOT NULL DEFAULT 0,
  params      TEXT NOT NULL,
  shared      TEXT NOT NULL DEFAULT 'PRIVATE',
  created     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  active      TEXT NOT NULL DEFAULT '1000-01-01 00:00:00',
  protected   INT(1) NOT NULL DEFAULT 0,
  hasmeta     INT(1) NOT NULL DEFAULT 0,
  hasud       INT(1) NOT NULL DEFAULT 0,
  hasis       INT(1) NOT NULL DEFAULT 0,
  info        TEXT NOT NULL DEFAULT '',
  infop       TEXT NOT NULL DEFAULT ''
);
CREATE UNIQUE INDEX info_id ON info(id);
CREATE INDEX info_owner ON info(owner);
CREATE INDEX info_status ON info(status);

CREATE TABLE macros (
  user   TEXT NOT NULL,
  macros TEXT NOT NULL
);
CREATE UNIQUE INDEX macros_user ON macros(user);

CREATE TABLE users (
  mail   TEXT NOT NULL,
  sec    TEXT NOT NULL,
  pw     TEXT NOT NULL,
  active TEXT NOT NULL DEFAULT '1000-01-01 00:00:00',
  quotum INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX users_mail ON users(mail);
CREATE INDEX users_sec ON users(sec);

CREATE TABLE version (
  id      INTEGER NOT NULL,
  version INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX version_id ON version(id);
EOT

    cat <<EOT | duckdb -bail -echo
INSTALL mysql;
INSTALL sqlite;

LOAD MYSQL;
ATTACH '' AS db1 (TYPE MYSQL, READ_ONLY);

LOAD SQLITE;
ATTACH '$PAQU/data/paqu.db.tmp' AS db2 (TYPE SQLITE);

INSERT INTO db2.corpora SELECT * FROM db1.${PREFIX}_corpora;
INSERT INTO db2.ignore  SELECT * FROM db1.${PREFIX}_ignore;
INSERT INTO db2.info    SELECT * FROM db1.${PREFIX}_info;
INSERT INTO db2.macros  SELECT * FROM db1.${PREFIX}_macros;
INSERT INTO db2.users   SELECT * FROM db1.${PREFIX}_users;
INSERT INTO db2.version SELECT * FROM db1.${PREFIX}_version;
EOT

    mv $PAQU/data/paqu.db.tmp $PAQU/data/paqu.db

fi

for db in $(sqlite3 $PAQU/data/paqu.db 'select id from info')
do
    echo ================================================================
    echo $db
    if [ -f $PAQU/data/$db/data.duckdb ]
    then
        continue
    fi
    rm -f $PAQU/data/$db/data.duckdb.tmp
    cat <<EOT | duckdb -bail -echo $PAQU/data/$db/data.duckdb.tmp
LOAD MYSQL;
ATTACH 'host=localhost' AS db1 (TYPE MYSQL, READ_ONLY);

CREATE TABLE arch   AS SELECT * FROM db1.${PREFIX}_c_${db}_arch;
CREATE UNIQUE INDEX arch_id ON arch(id);

CREATE TABLE file   AS SELECT * FROM db1.${PREFIX}_c_${db}_file;
CREATE UNIQUE INDEX file_id ON file(id);

CREATE TABLE deprel AS SELECT * FROM db1.${PREFIX}_c_${db}_deprel;
ALTER TABLE deprel ALTER word   TYPE VARCHAR COLLATE NOCASE.NOACCENT;
ALTER TABLE deprel ALTER lemma  TYPE VARCHAR COLLATE NOCASE.NOACCENT;
ALTER TABLE deprel ALTER root   TYPE VARCHAR COLLATE NOCASE.NOACCENT;
ALTER TABLE deprel ALTER hword  TYPE VARCHAR COLLATE NOCASE.NOACCENT;
ALTER TABLE deprel ALTER hlemma TYPE VARCHAR COLLATE NOCASE.NOACCENT;
ALTER TABLE deprel ALTER hroot  TYPE VARCHAR COLLATE NOCASE.NOACCENT;
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

CREATE TABLE sent   AS SELECT * FROM db1.${PREFIX}_c_${db}_sent;
CREATE INDEX sent_file ON sent(file);
CREATE INDEX sent_arch ON sent(arch);
CREATE INDEX sent_lbl  ON sent(lbl);

CREATE TABLE word   AS SELECT * FROM db1.${PREFIX}_c_${db}_word;
ALTER TABLE word ALTER word TYPE VARCHAR COLLATE NOCASE.NOACCENT;
CREATE UNIQUE INDEX word_word ON word(word);

EOT

    if [ $(echo 'select hasmeta from info where id = '\'$db\' | sqlite3 $PAQU/data/paqu.db) = 1 ]
    then
        cat <<EOT | duckdb -bail -echo $PAQU/data/$db/data.duckdb.tmp
LOAD MYSQL;
ATTACH 'host=localhost' AS db1 (TYPE MYSQL, READ_ONLY);

CREATE TABLE meta   AS SELECT * FROM db1.${PREFIX}_c_${db}_meta;
CREATE INDEX meta_id ON meta(id);
CREATE INDEX meta_file ON meta(file);
CREATE INDEX meta_arch ON meta(arch);
CREATE INDEX meta_tval ON meta(tval);
CREATE INDEX meta_ival ON meta(ival);
CREATE INDEX meta_fval ON meta(fval);
CREATE INDEX meta_dval ON meta(dval);
CREATE INDEX meta_idx  ON meta(idx);

CREATE TABLE midx   AS SELECT * FROM db1.${PREFIX}_c_${db}_midx;
CREATE UNIQUE INDEX midx_id ON midx(id);
CREATE INDEX midx_name ON midx("name");

CREATE TABLE minf   AS SELECT * FROM db1.${PREFIX}_c_${db}_minf;
CREATE INDEX minf_id ON minf(id);

CREATE TABLE mval   AS SELECT * FROM db1.${PREFIX}_c_${db}_mval;
CREATE INDEX mval_id ON mval(id);
CREATE INDEX mval_idx ON mval(idx);
EOT
    fi

    mv $PAQU/data/$db/data.duckdb.tmp $PAQU/data/$db/data.duckdb

done

