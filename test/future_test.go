package test

import (
	"github.com/jjj124/go-future"
	"log"
	"testing"
	"time"
)

func TestCase1(t *testing.T) {
	var foo = "foo"
	var bar = "bar"
	var f1 = futures.Just(&foo)
	var f2 = futures.Just(&bar)
	var f3 = f1.WhenComplete(func(s *string, err error) {
		log.Printf("recv %s sleep 3 seconds...", *s)
		time.Sleep(time.Second * 3)
	})
	var f4 = futures.Then(f2, func(t *string) *string {
		var ret = *t + "-"
		return &ret
	})
	futures.And(f4, f3).WhenComplete(func(t *futures.Tuple2[string, string], err error) {
		var t1 = *t.GetT1()
		var t2 = *t.GetT2()
		log.Printf("f4=%s,f3=%s", t1, t2)
	})
	time.Sleep(time.Second * 4)

}
func TestFromFun(t *testing.T) {

	var fun = func() *string {
		log.Printf("sleep 3 seconds ...")
		time.Sleep(time.Second * 3)
		var foo = "foo"
		return &foo
	}
	var f1 = futures.FromFunc[string](fun)
	var val, _ = f1.WhenComplete(func(s *string, err error) {
		log.Printf(*s)
	}).BlockingGet()
	log.Printf(*val)

	time.Sleep(time.Second * 4)

}

func TestThen(t *testing.T) {

}
