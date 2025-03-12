package cache

type Cache[K comparable, V any] interface {
	Add(key K, value V) bool

	Remove(key K) (k K, v V, ok bool)

	Get(key K) (value V, ok bool)
	Contains(key K) bool
	Keys() []K
	Values() []V

	Len() int
	Cap() int
}
