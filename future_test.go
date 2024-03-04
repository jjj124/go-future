package go_future

import (
	"log"
	"testing"
	"time"
)

func TestCase1(t *testing.T) {
	var foo = "foo"
	var bar = "bar"
	var f1 = Just(&foo)
	var f2 = Just(&bar)
	var f3 = f1.WhenComplete(func(s *string, err error) {
		log.Printf("recv %s sleep 3 seconds...", *s)
		time.Sleep(time.Second * 3)
	})
	var f4 = Then(f2, func(t *string) *string {
		var ret = *t + "-"
		return &ret
	})
	Zip(f4, f3).WhenComplete(func(t *Tuple2[string, string], err error) {
		var t1 = *t.GetT1()
		var t2 = *t.GetT2()
		log.Printf("f4=%s,f3=%s", t1, t2)
	})
	time.Sleep(time.Second * 4)

}
