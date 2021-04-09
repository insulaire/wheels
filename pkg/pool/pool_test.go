package pool

import (
	"testing"
)

func Test(t *testing.T) {
	p := NewPool(PoolConfig{
		Min: 5,
		Max: 10,
		InitFn: func() interface{} {
			return "abc"
		},
	})
	for i := 0; i < 100; i++ {
		if i%3 == 0 {
			p.Push("bcd")
		}
		p.Get()
	}

}
