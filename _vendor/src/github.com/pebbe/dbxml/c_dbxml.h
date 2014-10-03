#ifndef C_DBXML_H
#define C_DBXML_H

#ifdef __cplusplus
extern "C" {
#endif

    typedef struct c_dbxml_t *c_dbxml;

    typedef struct c_dbxml_result_t *c_dbxml_result;

    typedef struct c_dbxml_docs_t *c_dbxml_docs;

    typedef struct c_dbxml_query_t *c_dbxml_query;

    c_dbxml c_dbxml_open(char const *filename);
    void c_dbxml_free(c_dbxml db);

    int c_dbxml_error(c_dbxml db);
    char const * c_dbxml_errstring(c_dbxml db);

    /**** RESULTS ****/

    void c_dbxml_result_free(c_dbxml_result r);

    int c_dbxml_result_error(c_dbxml_result r);
    char const *c_dbxml_result_string(c_dbxml_result r);

    /**** WRITE ****/

    /* replace if replace != 0
     */
    c_dbxml_result c_dbxml_put_file(c_dbxml db, char const *filename, int replace);

    /* replace if replace != 0
     */
    c_dbxml_result c_dbxml_put_xml(c_dbxml db, char const *name, char const *data, int replace);

    /* replace if replace != 0
     */
    c_dbxml_result c_dbxml_merge(c_dbxml db, char const *dbxmlfile, int replace);

    c_dbxml_result c_dbxml_remove(c_dbxml db, char const *name);

    /**** READ ****/

    c_dbxml_result c_dbxml_get(c_dbxml db, char const * name);

    unsigned long long c_dbxml_size(c_dbxml db);

    c_dbxml_docs c_dbxml_get_all(c_dbxml db);
    c_dbxml_docs c_dbxml_get_query(c_dbxml db, char const *query);
    int c_dbxml_get_query_error(c_dbxml_docs docs);
    char const *c_dbxml_get_query_errstring(c_dbxml_docs docs);
    int c_dbxml_docs_next(c_dbxml_docs docs);
    char const * c_dbxml_docs_name(c_dbxml_docs docs);
    char const * c_dbxml_docs_content(c_dbxml_docs docs);
    char const * c_dbxml_docs_match(c_dbxml_docs docs);
    void c_dbxml_docs_free(c_dbxml_docs docs);
    c_dbxml_query c_dbxml_prepare_query(c_dbxml db, char const *query);
    c_dbxml_docs c_dbxml_run_query(c_dbxml_query query);
    void c_dbxml_cancel_query(c_dbxml_query query);
    void c_dbxml_query_free(c_dbxml_query query);
    int c_dbxml_get_prepared_error(c_dbxml_query query);
    char const *c_dbxml_get_prepared_errstring(c_dbxml_query query);

    /**** CHECK ****/

    c_dbxml_result c_dbxml_check(char const *query);

#ifdef __cplusplus
}
#endif

#endif /* C_DBXML_H */
