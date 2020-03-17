package workpool

import (
	"errors"
	"runtime"
	"sync"
)

var (
	ErrClosed       = errors.New("closed")
	ErrZeroPoolSize = errors.New("zero pool size")
)

type WorkPool struct {
	size   int
	queuec chan Task
	donec  chan struct{}
	wg     sync.WaitGroup
}

type Task func()

func New(size ...int) (*WorkPool, error) {
	wp := WorkPool{
		size:   runtime.GOMAXPROCS(0),
		queuec: make(chan Task),
		donec:  make(chan struct{}),
	}
	if len(size) > 0 {
		wp.size = size[0]
	}
	if wp.size == 0 {
		return nil, ErrZeroPoolSize
	}
	for i := 0; i < wp.size; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
	return &wp, nil
}

func (wp *WorkPool) Add(t Task) error {
	select {
	case <-wp.donec:
		return ErrClosed
	case wp.queuec <- t:
		return nil
	}
}

func (wp *WorkPool) Close() {
	select {
	case <-wp.donec:
	default:
		close(wp.donec)
		wp.wg.Wait()
	}
}

func (wp *WorkPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.donec:
			return
		case t := <-wp.queuec:
			t()
		}
	}
}
