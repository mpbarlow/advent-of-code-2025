package support

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(item T) bool {
	if !s.Has(item) {
		s[item] = struct{}{}

		return true
	}

	return false
}

func (s Set[T]) Has(item T) bool {
	_, ok := s[item]

	return ok
}

func NewSet[T comparable](initial ...T) Set[T] {
	s := make(Set[T], len(initial))

	for _, v := range initial {
		s.Add(v)
	}

	return s
}
