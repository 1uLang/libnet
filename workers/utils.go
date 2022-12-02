package workers

import (
	"sync"
)

var workers = []*Worker{}
var count = 0
var index = 0
var locker = sync.Mutex{}

//func init() {
//	count = runtime.NumCPU() * 4096
//	if count == 0 {
//		count = 8
//	}
//	for i := 0; i < count; i++ {
//		workers = append(workers, NewWorker(strconv.Itoa(i)))
//	}
//}
//
//func Get() *Worker {
//	locker.Lock()
//	defer locker.Unlock()
//
//	index++
//	if index > count-1 {
//		index = 0
//	}
//
//	return workers[index]
//}
