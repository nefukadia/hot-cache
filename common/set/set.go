package set

type Set[T comparable] interface {
	Insert(x T) bool
	Del(x T) bool
	Exist(x T) bool
	Len() int
	Range(f func(int, T) bool)
}
