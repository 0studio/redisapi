package redisapi

import (
	"strconv"
)

type ScoreInterface interface {
	GetScore() interface{}
	GetMember() interface{}
}
type OrderSetRedis interface {
	Zadd(key string, score int, value interface{}) error
	ZaddBatch(key string, list []ScoreInterface) error
	Zcard(key string) (int, error)
	ZRrange(key string, begin int, end int) ([]ScoreStruct, error)
	ZRevRrange(key string, begin int, end int) ([]ScoreStruct, error)
	ZRrank(key string, value interface{}) (int, error)
	ZRevRrank(key string, value interface{}) (int, error)
	Zrem(key string, value interface{}) error
	ZRemRangeByRank(key string, begin int, end int) error
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
	Ping() bool
	Exists(key string) bool

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
type ScoreStruct struct {
	Member interface{}
	Score  interface{}
}

func (ss ScoreStruct) GetMemberAsString() string {
	return string(ss.Member.([]uint8))
}
func (ss ScoreStruct) GetMemberAsUint64() (member uint64) {
	var i int
	i, _ = strconv.Atoi(string(ss.Member.([]uint8)))
	return uint64(i)
}

func (ss ScoreStruct) GetScoreAsInt() (score int) {
	score, _ = strconv.Atoi(string(ss.Score.([]uint8)))
	return
}
