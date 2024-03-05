package futures

import (
	"errors"
	"sync"
	"time"
)

type TimeoutErr struct {
}

func (t *TimeoutErr) Error() string {
	return "execute timeout !"
}

type Future[T any] struct {
	val    *T
	err    error
	ch     []chan int
	lock   sync.Locker
	result int
}

func NewFuture[T any]() *Future[T] {
	var x = &Future[T]{
		nil,
		nil,
		make([]chan int, 0),
		&sync.Mutex{},
		0,
	}
	return x
}
func (f *Future[T]) newWaiter() chan int {
	f.lock.Lock()
	defer f.lock.Unlock()
	var ch = make(chan int, 1)
	f.ch = append(f.ch, ch)
	if f.result != 0 {
		ch <- f.result
	}
	return ch
}
func (f *Future[T]) BlockingGet() (*T, error) {
	var ch = f.newWaiter()
	var r = <-ch
	if r > 0 {
		return f.val, nil
	} else {
		return nil, f.err
	}
}
func (f *Future[T]) BlockingUntil(duration time.Duration) (*T, error) {
	var ch = f.newWaiter()
	select {
	case r := <-ch:
		if r > 0 {
			return f.val, nil
		} else {
			return nil, f.err
		}
	case <-time.After(duration):
		return nil, &TimeoutErr{}
	}
}
func (f *Future[T]) Complete(value *T) bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.result != 0 {
		return false
	}

	f.val = value
	f.result = 1
	for _, chn := range f.ch {
		chn <- 1
	}
	return true
}
func (f *Future[T]) CompleteExceptionally(e error) bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.result != 0 {
		return false
	}
	f.err = e
	f.result = -1
	for _, c := range f.ch {
		c <- -1
	}
	return true
}
func (f *Future[T]) WhenComplete(fc func(*T, error)) *Future[T] {
	var ch = f.newWaiter()
	var ret = NewFuture[T]()
	go func() {
		r := <-ch
		func() {
			defer func() {
				var e = recover()
				if e != nil {
					ret.CompleteExceptionally(errors.New("execute map function occur error!"))
				}
			}()
			fc(f.val, f.err)
		}()
		if r > 0 {
			ret.Complete(f.val)
		} else {
			ret.CompleteExceptionally(f.err)
		}
	}()
	return ret
}
func (f *Future[T]) Or(other *Future[T]) *Future[T] {
	var ret = NewFuture[T]()
	f.WhenComplete(func(t *T, err error) {
		if err == nil {
			ret.Complete(t)
			return
		}
		other.WhenComplete(func(t *T, err error) {
			if err != nil {
				ret.CompleteExceptionally(err)
				return
			}
			ret.Complete(t)
		})
	})
	return ret
}
func (f *Future[T]) Then(fc func(*T) *any) *Future[any] {
	return Then(f, fc)
}
func (f *Future[T]) Delay(duration time.Duration) *Future[T] {
	var ret = NewFuture[T]()
	f.WhenComplete(func(t *T, err error) {
		time.Sleep(duration)
		if err != nil {
			ret.CompleteExceptionally(err)
			return
		}
		ret.Complete(t)
	})
	return ret
}
func Then[T any, V any](f *Future[T], fc func(*T) *V) *Future[V] {
	var ret = NewFuture[V]()
	f.WhenComplete(func(t *T, err error) {
		if err != nil {
			ret.CompleteExceptionally(err)
			return
		}
		func() {
			defer func() {
				var e = recover()
				if e != nil {
					ret.CompleteExceptionally(errors.New("execute func occur err !"))
				}
			}()
			var v2 = fc(t)
			ret.Complete(v2)
		}()
	})
	return ret
}
func And[T1 any, T2 any](f1 *Future[T1], f2 *Future[T2]) *Future[Tuple2[T1, T2]] {
	var ret = NewFuture[Tuple2[T1, T2]]()
	var waitGroup = sync.WaitGroup{}
	waitGroup.Add(2)
	f1.WhenComplete(func(t *T1, err error) {
		if err != nil {
			ret.CompleteExceptionally(err)
			waitGroup.Done()
			return
		}
		waitGroup.Done()
	})
	f2.WhenComplete(func(t *T2, err error) {
		if err != nil {
			ret.CompleteExceptionally(err)
			waitGroup.Done()
			return
		}
		waitGroup.Done()
	})
	go func() {
		waitGroup.Wait()
		var v1, err1 = f1.BlockingGet()
		if err1 != nil {
			ret.CompleteExceptionally(err1)
			return
		}
		var v2, err2 = f2.BlockingGet()
		if err2 != nil {
			ret.CompleteExceptionally(err2)
			return
		}
		var tuple2 = NewTuple2(v1, v2)
		ret.Complete(tuple2)
	}()
	return ret
}
func Just[T any](t *T) *Future[T] {
	var f = NewFuture[T]()
	f.Complete(t)
	return f
}
func Error[T any](err error) *Future[T] {
	var f = NewFuture[T]()
	f.CompleteExceptionally(err)
	return f
}

func FromFunc[T any](f func() *T) *Future[T] {
	var ret = NewFuture[T]()
	go func() {
		defer func() {
			var e = recover()
			if e != nil {
				ret.CompleteExceptionally(errors.New("execute function occur err!"))
			}
		}()
		var v = f()
		ret.Complete(v)
	}()
	return ret
}
