package cache

import (
	"log"
	"testing"
)

func TestSetGet(t *testing.T) {
	ca := New(21, nil)
	ca.DelOldst()
	ca.Set("k1", "abcd")
	ca.Set("k2", "dddd")
	log.Println(ca.Len())
	v := ca.Get("k2")
	log.Println(v)

	ca.Del("k1")

	log.Println(ca.Len())

}
