package cache

type Cache[K comparable, V any] interface {
	New() Cache[K, V]

	Add(key K, value V) bool

	Remove(key K) bool

	Get(key K) (value V, ok bool)
	Contains(key K) bool
	Keys() []K
	Values() []V

	Len() int
	Cap() int
}
