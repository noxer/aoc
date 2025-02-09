package utils

type Queue[T any] struct {
	p, s int
	data []T
}

func (q *Queue[T]) ensureSpace(size int) {
	if len(q.data)-q.s < size {

	}
}

func (q *Queue[T]) Push(el T) {

}
