package futures

type Tuple2[T1 any, T2 any] struct {
	t1 *T1
	t2 *T2
}

func (t *Tuple2[T1, T2]) GetT1() *T1 {
	return t.t1
}
func (t *Tuple2[T1, T2]) GetT2() *T2 {
	return t.t2
}

func NewTuple2[T1 any, T2 any](t1 *T1, t2 *T2) *Tuple2[T1, T2] {
	return &Tuple2[T1, T2]{
		t1: t1,
		t2: t2,
	}
}
