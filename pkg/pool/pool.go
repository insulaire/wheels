package pool

import (
	"log"
	"sync"
	"sync/atomic"
)

//关闭对象
type closeFn func(interface{}) bool

//初始化对象
type initFn func() interface{}

//
type watchCloseFn func(interface{})

type expiredFn func(interface{}) bool

//对象池配置
type PoolConfig struct {
	CloseFn   closeFn
	InitFn    initFn
	WatchFn   watchCloseFn
	ExpiredFn expiredFn
	Min, Max  uint32
}

//对象
type element struct {
	value      interface{}
	prve, next *element
}

type Pool interface {
	Get() interface{}
	Push(v interface{})
}

//池
type pool struct {
	start, end     *element
	mu             sync.Mutex
	closeFn        closeFn
	initFn         initFn
	watchFn        watchCloseFn
	expiredFn      expiredFn
	min, max, size uint32
}

func NewPool(config PoolConfig) Pool {
	if config.InitFn == nil {
		log.Panicln("InitFn not is nil")
	}
	if config.ExpiredFn == nil {
		config.ExpiredFn = func(i interface{}) bool {
			return false
		}
	}
	if config.WatchFn == nil {
		config.WatchFn = func(i interface{}) {
		}
	}
	if config.CloseFn == nil {
		config.CloseFn = func(i interface{}) bool {
			return true
		}
	}
	return &pool{
		closeFn:   config.CloseFn,
		initFn:    config.InitFn,
		watchFn:   config.WatchFn,
		expiredFn: config.ExpiredFn,
		mu:        sync.Mutex{},
		min:       config.Min,
		max:       config.Max,
	}
}

func (p *pool) put(el *element) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.start == nil {
		p.start = el
	} else {
		p.start.prve = el
		el.next = p.start
		p.start = el
	}
	if p.end == nil {
		p.end = el
	}

	atomic.AddUint32(&p.size, 1)
}

func (p *pool) newElement() {
	p.put(&element{value: p.initFn()})
}

func (p *pool) pop() *element {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.end == nil {
		return nil
	}
	v := p.end
	p.end = v.prve

	atomic.AddUint32(&p.size, ^uint32(1-1))
	return v
}

func (p *pool) Get() interface{} {
	for p.min > p.size {
		p.newElement()
	}
	for {
		v := p.pop()
		if v != nil && p.expiredFn(v) && p.closeFn(v) {
			go p.watchFn(v)
			continue
		}
		if v == nil {
			return v
		}
		return v.value
	}
}

func (p *pool) Push(v interface{}) {
	if p.expiredFn(v) {
		return
	}
	p.put(&element{value: v})
}
