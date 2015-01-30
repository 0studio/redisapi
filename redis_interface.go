package redisapi

type OrderSetRedis interface {
	Exists(key string) bool

	Zadd(key string, score int, value interface{}) error

	Zrem(key string, value interface{}) error

	Zcard(key string) (int, error)

	ZRrange(key string, begin int, end int) ([]interface{}, error)

	ZRrank(key string, value interface{}) (int, error)
}

type HashRedis interface {
	// Hdel(table, key string, value interface{}) error

	// Hset(table, key string, value interface{})
}

type QueueRedis interface {
	Lpush(key string, value interface{}) error
	Rpush(key string, value interface{}) error
	Lrange(key string, start, end int) ([]interface{}, error)
	Rpop(key string) (interface{}, error)
	Lset(key string, index int, value interface{}) error
	Ltrim(key string, start, end int) error
	Brpop(key string, timeoutSecs int) (interface{}, error)
	Lrem(key string, value interface{}, remType int) error
}

type Redis interface {
	QueueRedis
	OrderSetRedis
	HashRedis

	Set(key string, value []byte) error

	Get(key string) ([]byte, error)

	Delete(key string) error

	Incr(key string, step uint64) (int64, error)

	Decr(key string, step uint64) (int64, error)

	MultiGet(keys []interface{}) ([]interface{}, error)

	MultiSet(kvMap map[string][]byte) error

	ClearAll() error

	Pub(key string, value interface{}) error
	Sub(keys ...string) ([]string, error)
	UnSub(keys ...string) error
}
