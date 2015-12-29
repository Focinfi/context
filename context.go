package context

import (
	"net/http"
)

var data = make(map[*http.Request]map[interface{}]interface{})

func Set(req *http.Request, key, value interface{}) {
	if data[req] == nil {
		data[req] = make(map[interface{}]interface{})
	}
	data[req][key] = value
}

func Get(req *http.Request, key interface{}) interface{} {
	return data[req][key]
}

func GetOk(req *http.Request, key interface{}) (val interface{}, ok bool) {
	if val = Get(req, key); val != nil {
		ok = true
	}
	return
}

func GetAll(req *http.Request) map[interface{}]interface{} {
	vals := make(map[interface{}]interface{})
	for key := range data[req] {
		vals[key] = data[req][key]
	}
	return vals
}

func GetAllOk(req *http.Request) (allVal map[interface{}]interface{}, ok bool) {
	if allVal = GetAll(req); len(allVal) != 0 {
		ok = true
	}
	return
}

func Delete(req *http.Request, key interface{}) {
	if _, ok := GetOk(req, key); ok {
		delete(data[req], key)
	}
}

func Clear(req *http.Request) {
	if _, ok := GetAllOk(req); ok {
		delete(data, req)
	}
}

func ClearHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer Clear(req)
		handler.ServeHTTP(rw, req)
	})
}
