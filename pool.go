package velty

import (
	"github.com/viant/velty/est"
	"sync"
	"sync/atomic"
)

type (
	Pool struct {
		statePool *sync.Pool
		lock      *sync.RWMutex
		counter   int64
		size      int64
	}
)

func (p *Pool) State() *est.State {
	atomic.AddInt64(&p.counter, 1)

	for {
		state := p.statePool.Get().(*est.State)
		if state.Take() {
			return state
		}
	}
}

func (p *Pool) Put(state *est.State) {
	if atomic.AddInt64(&p.counter, -1) > p.size-1 {
		return
	}

	state.Reset()
	p.statePool.Put(state)
}

func NewPool(size int, newState func() *est.State) *Pool {
	statePool := &sync.Pool{
		New: func() interface{} {
			return newState()
		},
	}

	return &Pool{
		statePool: statePool,
		counter:   int64(0),
		size:      int64(size),
		lock:      &sync.RWMutex{},
	}
}
