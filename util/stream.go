package util

type Stream[T any] interface {
	Length() int
	Skip(i uint) Stream[T]
	Limit(i uint) Stream[T]
	Filter(func(v T) bool) Stream[T]
	ForEach(func(v T) T)
	Do() Stream[T]
	Get() []T
	Set(vs []T) Stream[T]
}

type StreamArray[T any] struct {
	node  []T
	funcs []func(s *StreamArray[T]) *StreamArray[T]
}

func NewStreamArray[T any](s []T) *StreamArray[T] {
	return &StreamArray[T]{
		node:  s,
		funcs: make([]func(s *StreamArray[T]) *StreamArray[T], 0),
	}
}

func (s *StreamArray[T]) Length() int {
	return len(s.node)
}

func (s *StreamArray[T]) Skip(i uint) Stream[T] {
	s.funcs = append(s.funcs, func(s *StreamArray[T]) *StreamArray[T] {
		if s.node == nil || i > uint(s.Length()) {
			s.node = nil
			return s
		}
		s.node = s.node[i:]
		return s
	})
	return s
}

func (s *StreamArray[T]) Limit(i uint) Stream[T] {
	s.funcs = append(s.funcs, func(s *StreamArray[T]) *StreamArray[T] {
		if s.node == nil {
			return s
		}
		if i > uint(s.Length()) {
			i = uint(s.Length())
		}
		s.node = s.node[:i]
		return s
	})
	return s
}

func (s *StreamArray[T]) Filter(f func(v T) bool) Stream[T] {
	s.funcs = append(s.funcs, func(s *StreamArray[T]) *StreamArray[T] {
		newNode := make([]T, 0)
		for i := range s.node {
			if f(s.node[i]) {
				newNode = append(newNode, s.node[i])
			}
		}
		s.node = newNode
		return s
	})
	return s
}

func (s *StreamArray[T]) ForEach(f func(v T) T) {
	s.funcs = append(s.funcs, func(s *StreamArray[T]) *StreamArray[T] {
		for i := range s.node {
			s.node[i] = f(s.node[i])
		}
		return s
	})
	s.Do()
}

func (s *StreamArray[T]) Do() Stream[T] {
	for _, v := range s.funcs {
		s = v(s)
	}
	s.funcs = make([]func(s *StreamArray[T]) *StreamArray[T], 0)
	return s
}

func (s *StreamArray[T]) Get() []T {
	return s.node
}

func (s *StreamArray[T]) Set(vs []T) Stream[T] {
	s.node = vs
	s.funcs = make([]func(s *StreamArray[T]) *StreamArray[T], 0)
	return s
}
