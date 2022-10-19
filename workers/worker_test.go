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
				time.Sleep(300 * time.Millisecond)
			})
		}(i)
	}

	time.Sleep(1 * time.Second)
}

func BenchmarkNewWorker(b *testing.B) {
	runtime.GOMAXPROCS(1)
	worker := NewWorker("1")
	for i := 0; i < b.N; i++ {
		worker.Run(func() {

		})
	}
}
