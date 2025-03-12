package loadbalance

type Loadbalance interface {
	Add(keys ...string)
	Get(key string) string
	Remove(key string)
}
