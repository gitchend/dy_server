package mpsc

import (
	"sync/atomic"
	"unsafe"
)

type (
	Queue[T any] struct {
		head   unsafe.Pointer
		tail   unsafe.Pointer
		length int64
	}
	Node[T any] struct {
		value T
		next  unsafe.Pointer
	}
)

func New[T any]() *Queue[T] {
	n := unsafe.Pointer(new(Node[T]))
	return &Queue[T]{
		head: n,
		tail: n,
	}
}

func (q *Queue[T]) Push(v T) {
	n := &Node[T]{value: v}
	for {
		tail := load[T](&q.tail)
		next := load[T](&tail.next)
		if tail == load[T](&q.tail) {
			atomic.AddInt64(&q.length, 1)
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n)
					return
				}
			} else {
				cas(&q.tail, tail, next)
			}
		}
	}
}

func (q *Queue[T]) Pop() (v T) {
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)
		if head == load[T](&q.head) {
			if head == tail {
				if next == nil {
					return
				}
				cas(&q.tail, tail, next)
			} else {
				v = next.value
				atomic.AddInt64(&q.length, -1)
				if cas(&q.head, head, next) {
					return
				}
			}
		}
	}
}
func (q *Queue[T]) Length() int64 {
	return q.length
}

func (q *Queue[T]) Empty() bool {
	return q.length == 0
}

func load[T any](p *unsafe.Pointer) *Node[T] {
	return (*Node[T])(atomic.LoadPointer(p))
}

func cas[T any](p *unsafe.Pointer, old, new *Node[T]) bool {
	return atomic.CompareAndSwapPointer(p,
		unsafe.Pointer(old), unsafe.Pointer(new))
}
