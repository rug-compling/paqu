#include "c_dbxml.h"
#include <dbxml/DbXml.hpp>
#include <string>

#define ALIAS "c_dbxml"

extern "C" {

    struct c_dbxml_t {
	DbXml::XmlManager manager;
	DbXml::XmlUpdateContext context;
	DbXml::XmlContainer container;
	DbXml::XmlContainerConfig config;
	bool error;
	std::string filename;
	std::string errstring;
    };

    struct c_dbxml_result_t {
	std::string result;
	bool error;
    };

    struct c_dbxml_docs_t {
	DbXml::XmlDocument doc;
	DbXml::XmlValue value;
	DbXml::XmlResults it;
	DbXml::XmlQueryContext context;
	bool validDoc;
	bool more;
	std::string name;
	std::string content;
	std::string match;
	std::string result;
	bool error;
	std::string errstring;
    };

    struct c_dbxml_query_t {
	DbXml::XmlQueryContext context;
	DbXml::XmlQueryExpression expression;
	bool error;
	std::string errstring;
    };

    c_dbxml c_dbxml_open(char const *filename, int readwrite, int read)
    {
	c_dbxml db;

	db = new c_dbxml_t;
	db->filename = filename;

	for (int i = 0; i < 2; i++) {
	    /* if both: first attempt is read+write */
	    if (i == 0 && readwrite == 0) {
		continue;
	    }
	    if (i == 1 && read == 0) {
		continue;
	    }
	    try {
		db->context = db->manager.createUpdateContext();
		if (i == 0) {
		    db->config.setAllowCreate(true);
		    db->config.setMode(0666);
		} else {
		    db->config.setAllowCreate(false);
		    db->config.setReadOnly(true);
		}
		db->container = db->manager.openContainer(filename, db->config);
		db->error = false;
		if (!db->container.addAlias(ALIAS)) {
		    db->errstring = "Unable to add alias \"" ALIAS "\"";
		    db->error = true;
		}
	    } catch (DbXml::XmlException &xe) {
		db->errstring = xe.what();
		db->error = true;
	    }
	    if (db->error == false) {
		break;
	    }
	}

	return db;
    }

    void c_dbxml_free(c_dbxml db)
    {
	delete db;
    }

    int c_dbxml_error(c_dbxml db)
    {
	return db->error ? 1 : 0;
    }

    char const *c_dbxml_errstring(c_dbxml db)
    {
	return db->errstring.c_str();
    }

    void c_dbxml_result_free(c_dbxml_result r)
    {
	delete r;
    }

    int c_dbxml_result_error(c_dbxml_result r)
    {
	return r->error ? 1 : 0;
    }

    char const *c_dbxml_result_string(c_dbxml_result r)
    {
	return r->result.c_str();
    }

    c_dbxml_result c_dbxml_put_file(c_dbxml db, char const * filename, int replace)
    {
	c_dbxml_result r;
	r = new c_dbxml_result_t;

	if (replace) {
	    try {
		db->container.deleteDocument(filename, db->context);
	    } catch (DbXml::XmlException &xe) {
		;
	    }
	}
        try {
            DbXml::XmlInputStream *is = db->manager.createLocalFileInputStream(filename);
            db->container.putDocument(filename, is, db->context, DbXml::DBXML_WELL_FORMED_ONLY);
	    r->error = false;
        } catch (DbXml::XmlException &xe) {
	    r->result = xe.what();
	    r->error = true;
        }

	return r;
    }

    // replace if replace != 0
    c_dbxml_result c_dbxml_put_xml(c_dbxml db, char const *name, char const *data, int replace)
    {
	c_dbxml_result r;
	r = new c_dbxml_result_t;

	if (replace) {
	    try {
		db->container.deleteDocument(name, db->context);
	    } catch (DbXml::XmlException &xe) {
		;
	    }
	}

        try {
            db->container.putDocument(name, data, db->context);
	    r->error = false;
        } catch (DbXml::XmlException &xe) {
	    r->result = xe.what();
	    r->error = true;
        }
	return r;
    }

    // replace if replace != 0
    c_dbxml_result c_dbxml_merge(c_dbxml db, char const * dbxmlfile, int replace) {
	c_dbxml_result r;
	r = new c_dbxml_result_t;

	DbXml::XmlContainer input = db->manager.openContainer(dbxmlfile);
	DbXml::XmlDocument doc;
	DbXml::XmlResults it = input.getAllDocuments(DbXml::DBXML_LAZY_DOCS);
	while (it.next(doc)) {
	    if (replace) {
		try {
		    db->container.deleteDocument(doc.getName(), db->context);
		} catch (DbXml::XmlException &xe) {
		    ;
		}
	    }
	    try {
		db->container.putDocument(doc, db->context);
		r->error = false;
	    } catch (DbXml::XmlException &xe) {
		r->result = xe.what();
		r->error = true;
		return r;
	    }
	}
	return r;
    }

    c_dbxml_result c_dbxml_remove(c_dbxml db, char const * filename)
    {
	c_dbxml_result r;
	r = new c_dbxml_result_t;

	try {
	    db->container.deleteDocument(filename, db->context);
	    r->error = false;
        } catch (DbXml::XmlException &xe) {
	    r->result = xe.what();
	    r->error = true;
	}
	return r;
    }

    c_dbxml_result c_dbxml_get(c_dbxml db, char const * name)
    {
	c_dbxml_result r;
	r = new c_dbxml_result_t;
	try {
	    DbXml::XmlDocument doc = db->container.getDocument(name);
	    doc.getContent(r->result);
	    r->error = false;
	} catch (DbXml::XmlException &xe) {
	    r->result = xe.what();
	    r->error = true;
	}
	return r;
    }

    unsigned long long c_dbxml_size(c_dbxml db)
    {
	return (unsigned long long) db->container.getNumDocuments();
    }

    c_dbxml_docs c_dbxml_get_all(c_dbxml db)
    {
	c_dbxml_docs docs;
	docs = new c_dbxml_docs_t;
	docs->it = db->container.getAllDocuments(DbXml::DBXML_LAZY_DOCS);
	docs->more = true;
	docs->error = false;
	return docs;
    }

    c_dbxml_query c_dbxml_prepare_query(c_dbxml db, char const *query, int useImplicitCollection, char const **namespaces)
    {
	int i;
	c_dbxml_query q;
	q = new c_dbxml_query_t;
	try {
	    q->context = db->manager.createQueryContext(DbXml::XmlQueryContext::LiveValues, DbXml::XmlQueryContext::Lazy);
	    q->context.setDefaultCollection(ALIAS);
	    for (i = 0; namespaces[i]; i += 2) {
		q->context.setNamespace(namespaces[i], namespaces[i+1]);
	    }
	    q->expression = db->manager.prepare(useImplicitCollection ? std::string("collection('" ALIAS "')") + query : query, q->context);
	    q->error = false;
	    if (q->expression.isUpdateExpression()) {
		q->errstring = "Update Expressions are not allowed";
		q->error = true;
	    }
	} catch (DbXml::XmlException const &xe) {
	    q->errstring = xe.what();
	    q->error = true;
	}
	return q;
    }

    c_dbxml_docs c_dbxml_run_query(c_dbxml_query query)
    {
	c_dbxml_docs docs;
	docs = new c_dbxml_docs_t;
	docs->more = true;
	docs->context = query->context;
	try {
	    docs->it = query->expression.execute(docs->context,
						 DbXml::DBXML_LAZY_DOCS | DbXml::DBXML_WELL_FORMED_ONLY
						 );
	    docs->error = false;
	} catch (DbXml::XmlException const &xe) {
	    docs->more = false;
	    docs->errstring = xe.what();
	    docs->error = true;
	}

	return docs;
    }

    int c_dbxml_get_query_error(c_dbxml_docs docs)
    {
	return docs->error ? 1 : 0;
    }

    char const *c_dbxml_get_query_errstring(c_dbxml_docs docs)
    {
	return docs->errstring.c_str();
    }

    int c_dbxml_get_prepared_error(c_dbxml_query query)
    {
	return query->error ? 1 : 0;
    }

    char const *c_dbxml_get_prepared_errstring(c_dbxml_query query)
    {
	return query->errstring.c_str();
    }

    int c_dbxml_docs_next(c_dbxml_docs docs)
    {
	if (docs->more) {

	    // goal: advance to next result and get both doc and value
	    // two times peek() and then next() to prevent a segment violation

	    try {
		docs->more = docs->it.peek(docs->value);
	    } catch (DbXml::XmlException &xe) {
		// while there are more results, this should always succeed, as the result is always an XmlValue
		docs->errstring = xe.what();
		docs->error = true;
		docs->more = false;
	    }

	    docs->validDoc = false;
	    if (docs->more) {
		try {
		    // this would cause a segment violation if first peek() failed
		    docs->it.peek(docs->doc);
		    docs->validDoc = true;
		} catch (...) {
		    // result is not always of type document, but we don't know this until we try.
		}

		// repeat first peek() as next() to advance
		try {
		    docs->more = docs->it.next(docs->value);
		} catch (DbXml::XmlException &xe) {
		    docs->errstring = xe.what();
		    docs->error = true;
		    docs->more = false;
		}
	    }

	    docs->name.clear();
	    docs->content.clear();
	    docs->match.clear();
	    docs->result.clear();
	}
	return docs->more ? 1 : 0;
    }

    char const * c_dbxml_docs_name(c_dbxml_docs docs)
    {
	if (docs->more && ! docs->name.size()) {
	    if (docs->validDoc)
		docs->name = docs->doc.getName();
	    else
		docs->name = "";
	}

	return docs->name.c_str();
    }

    char const * c_dbxml_docs_content(c_dbxml_docs docs)
    {
	if (docs->more && ! docs->content.size()) {
	    if (docs->validDoc) {
		docs->doc.getContent(docs->content);
	    } else {
		docs->content = "";
	    }
	}

	return docs->content.c_str();
    }

    char const * c_dbxml_docs_match(c_dbxml_docs docs)
    {
	if (docs->more && ! docs->match.size() && docs->value.isNode()) {
	    docs->match = docs->value.asString();
	}

	return docs->match.c_str();
    }

    char const * c_dbxml_docs_value(c_dbxml_docs docs)
    {
	if (docs->more && ! docs->result.size()) {
	    docs->result = docs->value.asString();
	}

	return docs->result.c_str();
    }

    void c_dbxml_docs_free(c_dbxml_docs docs)
    {
	delete docs;
    }

    void c_dbxml_query_free(c_dbxml_query query)
    {
	query->context.clearNamespaces(); // is this necessary?
	delete query;
    }

    void c_dbxml_cancel_query(c_dbxml_query query)
    {
	query->context.interruptQuery();
    }

    c_dbxml_result c_dbxml_check(char const *query, char const **namespaces)
    {
	c_dbxml_result r;
	r = new c_dbxml_result_t;
	int i;
	try {
	    DbXml::XmlManager manager;
	    DbXml::XmlQueryContext context;
	    DbXml::XmlQueryExpression expr;
	    context = manager.createQueryContext(DbXml::XmlQueryContext::LiveValues, DbXml::XmlQueryContext::Lazy);
	    for (i = 0; namespaces[i]; i += 2) {
		context.setNamespace(namespaces[i], namespaces[i+1]);
	    }
	    expr = manager.prepare(std::string(query), context);
	    r->error = false;
	    if (expr.isUpdateExpression()) {
		r->result = "Update Expressions are not allowed";
		r->error = true;
	    }
	    context.clearNamespaces(); // is this necessary?
	} catch (DbXml::XmlException const &xe) {
	    r->result = xe.what();
	    r->error = true;
	}
	return r;
    }

    void c_dbxml_version(int *major, int *minor, int *patch)
    {
	DbXml::dbxml_version(major, minor, patch);
    }

}
