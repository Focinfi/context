package context

import (
	"net/http"
	"sync"
	"time"
)

var (
	mux   sync.RWMutex
	data  = make(map[*http.Request]map[interface{}]interface{})
	dataT = make(map[*http.Request]int64)
)

func Set(req *http.Request, key, value interface{}) {
	mux.Lock()
	if data[req] == nil {
		data[req] = make(map[interface{}]interface{})
	}
	data[req][key] = value
	dataT[req] = time.Now().Unix()
	mux.Unlock()
}

func Get(req *http.Request, key interface{}) interface{} {
	mux.RLock()
	defer mux.RUnlock()
	return data[req][key]
}

func GetOk(req *http.Request, key interface{}) (val interface{}, ok bool) {
	if val = Get(req, key); val != nil {
		ok = true
	}
	return
}

func GetAll(req *http.Request) map[interface{}]interface{} {
	mux.RLock()
	defer mux.RUnlock()
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
		mux.Lock()
		delete(data[req], key)
		mux.Unlock()
	}
}

func Clear(req *http.Request) {
	if _, ok := GetAllOk(req); ok {
		mux.Lock()
		delete(data, req)
		delete(dataT, req)
		mux.Unlock()
	}
}

func ClearHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer Clear(req)
		handler.ServeHTTP(rw, req)
	})
}

func Purge(maxAge int) (count int) {
	if maxAge <= 0 {
		count = len(data)
		mux.Lock()
		data = make(map[*http.Request]map[interface{}]interface{})
		dataT = make(map[*http.Request]int64)
		mux.Unlock()
	} else {
		for req := range data {
			if time.Now().Unix() >= dataT[req]+int64(maxAge) {
				Clear(req)
				count++
			}
		}
	}
	return
}
