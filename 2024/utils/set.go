package utils

type Set[T comparable] map[T]struct{}

func SetFromSlice[T comparable, S ~[]T](sl S) Set[T] {
	s := make(Set[T], len(sl))

	for _, t := range sl {
		s[t] = struct{}{}
	}

	return s
}

func (s Set[T]) Put(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Has(t T) bool {
	_, ok := s[t]
	return ok
}
