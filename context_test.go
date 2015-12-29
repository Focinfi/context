package context

import (
	"net/http"
	"testing"
)

var request, _ = http.NewRequest("get", "localhost:8080/", nil)
var emptyReq, _ = http.NewRequest("get", "localhost:8080/", nil)

const (
	key1 = iota
	key2
)

func assertEqual(t *testing.T, result, expected interface{}) {
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSet(t *testing.T) {
	Set(request, key1, "value")
	assertEqual(t, data[request][key1], "value")
}

func TestGet(t *testing.T) {
	val := Get(emptyReq, key1)
	val = "v"
	val = Get(emptyReq, key1)
	assertEqual(t, val, nil)
	val = Get(request, key1)
	assertEqual(t, val, "value")
}

func TestGetok(t *testing.T) {
	val, ok := GetOk(request, key1)
	assertEqual(t, val, "value")
	assertEqual(t, ok, true)
}

func TestGetAll(t *testing.T) {
	allVal := GetAll(emptyReq)
	assertEqual(t, len(allVal), 0)

	allVal = GetAll(request)
	allVal[key2] = "value2"
	allVal = GetAll(request)
	assertEqual(t, len(allVal), 1)
}

func TestGetAllOk(t *testing.T) {
	allVal, ok := GetAllOk(request)
	assertEqual(t, len(allVal), 1)
	assertEqual(t, ok, true)
}

func TestDelete(t *testing.T) {
	Delete(request, key1)
	allVal := GetAll(request)
	assertEqual(t, len(allVal), 0)
}

func TestClear(t *testing.T) {
	Clear(request)
	_, ok := GetAllOk(request)
	assertEqual(t, ok, false)
}

func TestClearHandler(t *testing.T) {
	clearHandler := ClearHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))
	_, ok := clearHandler.(http.Handler)
	assertEqual(t, ok, true)
}
