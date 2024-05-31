package goroutine

import (
	"runtime"
	"sync"

	"app/container/safe/mpsc"
	"go.uber.org/atomic"
)

/*
异步有序任务队列
当任务队列空闲时自动释放goroutine
*/

type Task func()

const (
	idle int32 = iota
	running
)

const DefaultGpsched = 128

var _workerPool = sync.Pool{New: func() interface{} { return newWorker() }}

func GetWorker() *GoWorker  { return _workerPool.Get().(*GoWorker) }
func PutWorker(w *GoWorker) { _workerPool.Put(w) }
func newWorker() *GoWorker  { return &GoWorker{queue: mpsc.New[Task](), gosched: DefaultGpsched} }

type GoWorker struct {
	queue           *mpsc.Queue[Task]
	schedulerStatus atomic.Int32
	gosched         int
}

func (w *GoWorker) SetGosched(gosched int) { w.gosched = gosched }

func (w *GoWorker) Push(fn Task) {
	w.queue.Push(fn)
	if w.schedulerStatus.CAS(idle, running) {
		if err := _pool.Submit(w.schedule); err != nil {
		}
	}
}

func (w *GoWorker) schedule() {
process:
	Try(w.run, nil)
	w.schedulerStatus.Store(idle)
	if !w.queue.Empty() && w.schedulerStatus.CAS(idle, running) {
		goto process
	}
}

func (w *GoWorker) run() {
	count := 0
	for {
		if w.gosched > 0 && count > w.gosched {
			count = 0
			runtime.Gosched()
		}

		fn := w.queue.Pop()
		if fn == nil {
			return
		}
		fn()
	}
}
