-- usage: sqlite3 paqu.db < mkSqlite1.sql

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

