package pool

import (
	"sync"
	"sync/atomic"
)

type closeFn func(*element)

type initFn func() interface{}

type watchFn func(*element)

type PoolConfig struct {
	CloseFn  closeFn
	InitFn   initFn
	WatchFn  watchFn
	Min, Max uint32
}

//池
type pool struct {
	start, end     *element
	mu             sync.Mutex
	closeFn        closeFn
	initFn         initFn
	watchFn        watchFn
	min, max, size uint32
}

func newPool(config PoolConfig) *pool {
	return &pool{
		closeFn: config.CloseFn,
		initFn:  config.InitFn,
		watchFn: config.WatchFn,
		mu:      sync.Mutex{},
		min:     config.Min,
		max:     config.Max,
	}
}

func (p *pool) add(el *element) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if el == nil {
		el = &element{value: p.initFn()}
	}
	if p.start == nil {
		p.start = el
	} else {
		el.next = p.start
		p.start = el
	}
	if p.end == nil {
		p.end = el
	}

	atomic.AddUint32(&p.size, 1)
}

//对象
type element struct {
	value      interface{}
	prve, next *element
}
