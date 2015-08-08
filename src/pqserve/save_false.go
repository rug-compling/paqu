// +build nodbxml

package main

func saveOpenDact(q *Context, prefix string, arch int) (interface{}, string) {
	return nil, ""
}

func saveGetDact(q *Context, dact interface{}, filename string) []byte {
	return []byte{}
}

func saveCloseDact(dact interface{}) {
}
