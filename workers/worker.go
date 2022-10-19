package workers

import (
	"time"
)

type Worker struct {
	id       string
	ch       chan func()
	isClosed bool
}

func NewWorker(id string) *Worker {
	worker := &Worker{
		id: id,
		ch: make(chan func(), 128),
	}
	worker.setup()
	return worker
}

func (this *Worker) Run(task func()) {
	defer func() {
		recover()
	}()

	// 此处可能造成阻塞，但保证了任务是同步执行的
	this.ch <- task
}

func (this *Worker) Id() string {
	return this.id
}

func (this *Worker) Close() {
	if this.isClosed {
		return
	}
	this.isClosed = true

	// 延时N秒再关闭
	timeout := time.NewTimer(5 * time.Second)
	go func() {
		defer func() {
			recover()
		}()

		select {
		case <-timeout.C:
			close(this.ch)
		}
	}()
}

func (this *Worker) setup() {
	go func() {
		for task := range this.ch {
			if task == nil {
				break
			}
			task()
		}
	}()
}
