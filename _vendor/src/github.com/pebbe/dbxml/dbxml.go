/*
A basic Go interface to Oracle Berkeley DB XML.

http://www.oracle.com/us/products/database/berkeley-db/xml/index.html
*/
package dbxml

//. Imports

/*
#cgo LDFLAGS: -ldbxml
#include <stdlib.h>
#include "c_dbxml.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

//. Types

// A database connection.
type Db struct {
	opened bool
	db     C.c_dbxml
	lock   sync.Mutex
}

// An iterator over xml documents in the database.
type Docs struct {
	started bool
	opened  bool
	docs    C.c_dbxml_docs
	lock    sync.Mutex
}

//. Variables

var (
	errclosed = errors.New("Database is closed")
)

//. Open & Close

// Open a database.
//
// Call db.Close() to ensure all write operations to the database are finished, before terminating the program.
func Open(filename string) (*Db, error) {
	db := &Db{}
	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	db.db = C.c_dbxml_open(cs)
	if C.c_dbxml_error(db.db) != 0 {
		err := errors.New(C.GoString(C.c_dbxml_errstring(db.db)))
		C.c_dbxml_free(db.db)
		return db, err
	}
	db.opened = true
	runtime.SetFinalizer(db, (*Db).Close)
	return db, nil
}

// Close the database.
//
// This flushes all write operations to the database.
func (db *Db) Close() {
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.opened {
		C.c_dbxml_free(db.db)
		db.opened = false
	}
}

//. Write

// Put an xml file from disc into the database.
func (db *Db) PutFile(filename string, replace bool) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return errclosed
	}

	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	repl := C.int(0)
	if replace {
		repl = 1
	}
	r := C.c_dbxml_put_file(db.db, cs, repl)
	defer C.c_dbxml_result_free(r)
	if C.c_dbxml_result_error(r) != 0 {
		return errors.New(C.GoString(C.c_dbxml_result_string(r)))
	}
	return nil
}

// Put an xml document from memory into the database.
func (db *Db) PutXml(name string, data string, replace bool) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return errclosed
	}

	csname := C.CString(name)
	defer C.free(unsafe.Pointer(csname))
	csdata := C.CString(data)
	defer C.free(unsafe.Pointer(csdata))
	repl := C.int(0)
	if replace {
		repl = 1
	}
	r := C.c_dbxml_put_xml(db.db, csname, csdata, repl)
	defer C.c_dbxml_result_free(r)
	if C.c_dbxml_result_error(r) != 0 {
		return errors.New(C.GoString(C.c_dbxml_result_string(r)))
	}
	return nil
}

// Merge a database from disc into this database.
func (db *Db) Merge(filename string, replace bool) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return errclosed
	}

	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	repl := C.int(0)
	if replace {
		repl = 1
	}
	r := C.c_dbxml_merge(db.db, cs, repl)
	defer C.c_dbxml_result_free(r)
	if C.c_dbxml_result_error(r) != 0 {
		return errors.New(C.GoString(C.c_dbxml_result_string(r)))
	}
	return nil
}

// Remove an xml document from the database.
func (db *Db) Remove(name string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return errclosed
	}

	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	r := C.c_dbxml_remove(db.db, cs)
	defer C.c_dbxml_result_free(r)
	if C.c_dbxml_result_error(r) != 0 {
		return errors.New(C.GoString(C.c_dbxml_result_string(r)))
	}
	return nil
}

//. Read

// Get an xml document by name from the database.
func (db *Db) Get(name string) (string, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return "", errclosed
	}

	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	r := C.c_dbxml_get(db.db, cs)
	defer C.c_dbxml_result_free(r)
	s := C.GoString(C.c_dbxml_result_string(r))
	if C.c_dbxml_result_error(r) != 0 {
		return "", errors.New(s)
	}
	return s, nil
}

// Get the number of xml documents in the database.
func (db *Db) Size() (uint64, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return 0, errclosed
	}
	return uint64(C.c_dbxml_size(db.db)), nil
}

// Get all xml documents from the database.
//
// Example:
//
//      docs, err := db.All()
//      if err != nil {
//          fmt.Println(err)
//      } else {
//          for docs.Next() {
//              fmt.Println(docs.Name(), docs.Content())
//          }
//      }
func (db *Db) All() (*Docs, error) {
	docs := &Docs{}
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return docs, errclosed
	}
	docs.docs = C.c_dbxml_get_all(db.db)
	runtime.SetFinalizer(docs, (*Docs).Close)
	docs.opened = true
	return docs, nil
}

// Get all xml documents that match the XPATH query from the database.
//
// Example:
//
//      docs, err := db.Query(xpath_query)
//      if err != nil {
//          fmt.Println(err)
//      } else {
//          for docs.Next() {
//              fmt.Println(docs.Name(), docs.Content())
//          }
//      }
func (db *Db) Query(query string) (*Docs, error) {
	docs := &Docs{}
	db.lock.Lock()
	defer db.lock.Unlock()

	if !db.opened {
		return docs, errclosed
	}
	cs := C.CString(query)
	defer C.free(unsafe.Pointer(cs))
	docs.docs = C.c_dbxml_get_query(db.db, cs)
	if C.c_dbxml_get_query_error(docs.docs) != 0 {
		return docs, errors.New(C.GoString(C.c_dbxml_get_query_errstring(docs.docs)))
	}
	runtime.SetFinalizer(docs, (*Docs).Close)
	docs.opened = true
	return docs, nil
}

// Iterate to the next xml document in the list, that was returned by db.All() or db.Query(query).
func (docs *Docs) Next() bool {
	docs.lock.Lock()
	defer docs.lock.Unlock()
	if !docs.opened {
		return false
	}
	if C.c_dbxml_docs_next(docs.docs) == 0 {
		docs.close()
		return false
	}
	docs.started = true
	return true
}

// Get name of current xml document after call to docs.Next().
func (docs *Docs) Name() string {
	return docs.getNameContent(1)
}

// Get content of current xml document after call to docs.Next().
func (docs *Docs) Content() string {
	return docs.getNameContent(2)
}

// Get matched subtree from current xml document after call to docs.Next().
func (docs *Docs) Match() string {
	return docs.getNameContent(3)
}

func (docs *Docs) getNameContent(what int) string {
	docs.lock.Lock()
	defer docs.lock.Unlock()
	if !(docs.opened && docs.started) {
		return ""
	}
	switch what {
	case 1:
		return C.GoString(C.c_dbxml_docs_name(docs.docs))
	case 2:
		return C.GoString(C.c_dbxml_docs_content(docs.docs))
	case 3:
		return C.GoString(C.c_dbxml_docs_match(docs.docs))
	}
	return ""
}

// Close iterator over xml documents in the database, that was returned by db.All() or db.Query(query).
//
// This is called automaticly if docs.Next() reaches false, but you can also call it inside a loop to exit it prematurely:
//
//      docs, _ := db.All()
//      for docs.Next() {
//          fmt.Println(docs.Name(), docs.Content())
//          if some_condition {
//              docs.Close()
//          }
//      }
func (docs *Docs) Close() {
	docs.lock.Lock()
	defer docs.lock.Unlock()
	docs.close()
}

func (docs *Docs) close() {
	if docs.opened {
		C.c_dbxml_docs_free(docs.docs)
		docs.opened = false
	}
}
