package utils

import (
	"cmp"
	"container/heap"
)

type heapSlice[T cmp.Ordered, S any] struct {
	score func(S) T
	data  []S
}

func (h heapSlice[T, S]) Len() int {
	return len(h.data)
}

func (h heapSlice[T, S]) Less(i, j int) bool {
	return h.score(h.data[i]) < h.score(h.data[j])
}

func (h heapSlice[T, S]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *heapSlice[T, S]) Push(x any) {
	h.data = append(h.data, x.(S))
}

func (h *heapSlice[T, S]) Pop() any {
	s := h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	return s
}

type Heap[T cmp.Ordered, S any] struct {
	data *heapSlice[T, S]
}

func NewHeap[S any, T cmp.Ordered](score func(S) T, elems ...S) Heap[T, S] {
	if score == nil {
		panic("you need to specify a score function")
	}

	hs := &heapSlice[T, S]{
		score: score,
		data:  elems,
	}

	h := Heap[T, S]{
		data: hs,
	}
	heap.Init(h.data)

	return h
}

func (h Heap[T, S]) Push(elem S) {
	heap.Push(h.data, elem)
}

func (h Heap[T, S]) Pop() S {
	return heap.Pop(h.data).(S)
}

func (h Heap[T, S]) Len() int {
	return len(h.data.data)
}
