INSTALL mysql;
INSTALL sqlite;

LOAD MYSQL;
ATTACH 'host=localhost' AS db1 (TYPE MYSQL, READ_ONLY);

LOAD SQLITE;
ATTACH 'paqu.db' AS db2 (TYPE SQLITE);

INSERT INTO db2.corpora SELECT * FROM db1.wordrel_corpora;
INSERT INTO db2.ignore  SELECT * FROM db1.wordrel_ignore;
INSERT INTO db2.info    SELECT * FROM db1.wordrel_info;
INSERT INTO db2.macros  SELECT * FROM db1.wordrel_macros;
INSERT INTO db2.users   SELECT * FROM db1.wordrel_users;
INSERT INTO db2.version SELECT * FROM db1.wordrel_version;

