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

func TestPurge(t *testing.T) {
	Set(request, key2, "value2")
	removedCount := Purge(0)
	assertEqual(t, removedCount, 1)
	assertEqual(t, len(GetAll(request)), 0)

	Set(request, key1, "value1")
	removedCount = Purge(10)
	assertEqual(t, removedCount, 0)
}

func parallelReader(r *http.Request, key string, iterations int, wait, done chan struct{}) {
	<-wait
	for i := 0; i < iterations; i++ {
		Get(r, key)
	}
	done <- struct{}{}

}

func parallelWriter(r *http.Request, key, value string, iterations int, wait, done chan struct{}) {
	<-wait
	for i := 0; i < iterations; i++ {
		Set(r, key, value)
	}
	done <- struct{}{}

}

func benchmarkMutex(b *testing.B, numReaders, numWriters, iterations int) {

	b.StopTimer()
	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	done := make(chan struct{})
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		wait := make(chan struct{})

		for i := 0; i < numReaders; i++ {
			go parallelReader(r, "test", iterations, wait, done)
		}

		for i := 0; i < numWriters; i++ {
			go parallelWriter(r, "test", "123", iterations, wait, done)
		}

		close(wait)

		for i := 0; i < numReaders+numWriters; i++ {
			<-done
		}

	}

}

func BenchmarkMutexSameReadWrite1(b *testing.B) {
	benchmarkMutex(b, 1, 1, 32)
}
func BenchmarkMutexSameReadWrite2(b *testing.B) {
	benchmarkMutex(b, 2, 2, 32)
}
func BenchmarkMutexSameReadWrite4(b *testing.B) {
	benchmarkMutex(b, 4, 4, 32)
}
func BenchmarkMutex1(b *testing.B) {
	benchmarkMutex(b, 2, 8, 32)
}
func BenchmarkMutex2(b *testing.B) {
	benchmarkMutex(b, 16, 4, 64)
}
func BenchmarkMutex3(b *testing.B) {
	benchmarkMutex(b, 1, 2, 128)
}
func BenchmarkMutex4(b *testing.B) {
	benchmarkMutex(b, 128, 32, 256)
}
func BenchmarkMutex5(b *testing.B) {
	benchmarkMutex(b, 1024, 2048, 64)
}
func BenchmarkMutex6(b *testing.B) {
	benchmarkMutex(b, 2048, 1024, 512)
}
