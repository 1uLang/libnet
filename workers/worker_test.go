package workers

import (
	"runtime"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	worker := NewWorker("1")
	for i := 0; i < 10; i++ {
		func(index int) {
			worker.Run(func() {
				t.Log("hello", index)
			})
		}(i)
	}

	time.Sleep(1 * time.Second)
}

func TestWorker_Panic(t *testing.T) {
	worker := NewWorker("1")
	worker.Run(func() {
		t.Log("0")
		var a = []int{}
		t.Log(a[1])
	})
	worker.Run(func() {
		t.Log("1")
	})
	worker.Run(func() {
		t.Log("2")
	})

	time.Sleep(1 * time.Second)
	t.Log("OK")
}

func BenchmarkNewWorker(b *testing.B) {
	runtime.GOMAXPROCS(1)
	worker := NewWorker("1")
	for i := 0; i < b.N; i++ {
		worker.Run(func() {

		})
	}
}
