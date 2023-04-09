package set

type HashSet[T comparable] struct {
	m map[T]struct{}
}

func NewHashSet[T comparable]() Set[T] {
	return &HashSet[T]{m: make(map[T]struct{})}
}

func (s *HashSet[T]) Insert(x T) bool {
	_, ok := s.m[x]
	cover := false
	if ok {
		cover = true
	}
	s.m[x] = struct{}{}
	return cover
}

func (s *HashSet[T]) Del(x T) bool {
	_, ok := s.m[x]
	if !ok {
		return false
	}
	delete(s.m, x)
	return true
}

func (s *HashSet[T]) Exist(x T) bool {
	_, ok := s.m[x]
	return ok
}

func (s *HashSet[T]) Len() int {
	return len(s.m)
}

func (s *HashSet[T]) Range(f func(int, T) bool) {
	var i = 0
	for key := range s.m {
		if !f(i, key) {
			return
		}
		i++
	}
}
