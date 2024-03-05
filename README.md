# go-future

  

# 1.Method List
    1. func NewFuture[T any]() *Future[T] 
    2. func (f *Future[T]) BlockingGet() (*T, error)
    3. func (f *Future[T]) BlockingUntil(duration time.Duration) (*T, error)
    4. func (f *Future[T]) Complete(value *T) bool
    5. func (f *Future[T]) CompleteExceptionally(e error) bool
    6. func (f *Future[T]) WhenComplete(fc func(*T, error)) *Future[T]
    7. func (f *Future[T]) Or(other *Future[T]) *Future[T]
    8. func (f *Future[T]) Then(fc func(*T) *any) *Future[any]
    9. func (f *Future[T]) Delay(duration time.Duration) *Future[T]
    10. func Then[T any, V any](f *Future[T], fc func(*T) *V) *Future[V]
    11. func And[T1 any, T2 any](f1 *Future[T1], f2 *Future[T2]) *Future[Tuple2[T1, T2]] 
    


# 2.How To Use

## 2.1.NewFuture And Complete

```go
var f = futures.NewFuture[string]()
f.WhenComplete(func(s *string,err error){
    if err != nil {
        log.Printf("task fail !")
    } else {
        log.Printf("task success, result: %s!", * s)
    }
})
gofunc() {
    time.Sleep(time.Second * 3)
    var r = rand.Intn(2)
    if r/2 == 0 {
        var str = strconv.Itoa(r)
        f.Complete(&str)
    } else {
        f.CompleteExceptionally(errors.New(""))
    }
}()
time.Sleep(time.Second * 4)

```
## 2.2.FromFunc

```go
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

```
## 2.3.And Future

```go
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
```
