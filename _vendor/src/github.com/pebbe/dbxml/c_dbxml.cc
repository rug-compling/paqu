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
	std::string errstring;
    };

    struct c_dbxml_result_t {
	std::string result;
	bool error;
    };

    struct c_dbxml_docs_t {
	DbXml::XmlDocument doc;
	DbXml::XmlResults it;
	DbXml::XmlQueryContext context;
	bool more;
	std::string name;
	std::string content;
	bool error;
	std::string errstring;
    };

    c_dbxml c_dbxml_open(char const *filename)
    {
	c_dbxml db;

	db = new c_dbxml_t;

	for (int i = 0; i < 2; i++) {
	    try {
		db->context = db->manager.createUpdateContext();
		if (i == 1) {
		    db->config.setReadOnly(true);
		}
		db->container = db->manager.existsContainer(filename) ?
		    db->manager.openContainer(filename, db->config) :
		    db->manager.createContainer(filename);
		if (!db->container.addAlias(ALIAS)) {
		    db->errstring = "Unable to add alias \"" ALIAS "\"";
		    db->error = true;
		}
		db->error = false;
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
            db->container.putDocument(filename, is, db->context);
	    r->error = false;
        } catch (DbXml::XmlException &xe) {
	    r->error = true;
	    r->result = xe.what();
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
	    r->error = true;
	    r->result = xe.what();
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
		r->error = true;
		r->result = xe.what();
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
	    r->error = true;
	    r->result = xe.what();
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

    c_dbxml_docs c_dbxml_get_query(c_dbxml db, char const *query)
    {
	c_dbxml_docs docs;
	docs = new c_dbxml_docs_t;
	docs->more = true;
	try {

	    docs->context = db->manager.createQueryContext(DbXml::XmlQueryContext::LiveValues, DbXml::XmlQueryContext::Lazy);
	    docs->context.setDefaultCollection(ALIAS);
	    docs->it = db->manager.query(std::string("collection('" ALIAS "')") + query,
					 docs->context,
					 DbXml::DBXML_LAZY_DOCS | DbXml::DBXML_WELL_FORMED_ONLY
					 );
	    docs->error = false;
	} catch (DbXml::XmlException const &xe) {
	    docs->more = false;
	    docs->error = true;
	    docs->errstring = xe.what();
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

    int c_dbxml_docs_next(c_dbxml_docs docs)
    {
	if (docs->more) {
	    docs->more = docs->it.next(docs->doc);
	    docs->name.clear();
	    docs->content.clear();
	}
	return docs->more ? 1 : 0;
    }

    char const * c_dbxml_docs_name(c_dbxml_docs docs)
    {
	if (docs->more && ! docs->name.size())
	    docs->name = docs->doc.getName();

	return docs->name.c_str();
    }

    char const * c_dbxml_docs_content(c_dbxml_docs docs)
    {
	if (docs->more && ! docs->content.size())
	    docs->doc.getContent(docs->content);

	return docs->content.c_str();
    }

    void c_dbxml_docs_free(c_dbxml_docs docs)
    {
	delete docs;
    }

}
