package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"strconv"
	"strings"

	"regexp"

	"github.com/mofancloud/xmicro/cache"

	"github.com/garyburd/redigo/redis"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "redis"
)

var tReg *regexp.Regexp

func init() {
	tReg, _ = regexp.Compile(` m=.*`)
}

type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	maxIdle  int
}

func NewRedisCache() cache.Cache {
	return &Cache{key: DefaultKey}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

func (rc *Cache) Set(key string, val interface{}) (reply interface{}, err error) {
	reply, err = rc.do("SET", key, val)
	return
}

func (rc *Cache) Incrby(key string, value int64) (reply interface{}, err error) {
	reply, err = rc.do("INCRBY", key, value)
	return
}

func (rc *Cache) SetWithExpired(key string, val interface{}, expired int) (reply interface{}, err error) {
	reply, err = rc.do("SET", key, val, "EX", expired)
	return
}

func (rc *Cache) LPush(key string, val interface{}) (reply interface{}, err error) {
	reply, err = rc.do("LPUSH", key, val)
	return
}

func (rc *Cache) LPop(key string) (reply interface{}, err error) {
	reply, err = rc.do("LPOP", key)
	return
}

func (rc *Cache) RPush(key string, val interface{}) (reply interface{}, err error) {
	reply, err = rc.do("RPUSH", key, val)
	return
}

func (rc *Cache) RPop(key string) (reply interface{}, err error) {
	reply, err = rc.do("RPOP", key)
	return
}

func (rc *Cache) LTRIM(key string, start int64, stop int64) (err error) {
	_, err = rc.do("LTRIM", key, start, stop)
	return err
}

func (rc *Cache) LRange(key string, start int64, stop int64) (reply interface{}, err error) {
	reply, err = rc.do("LRANGE", key, start, stop)
	return
}

func (rc *Cache) SADD(key string, vals ...interface{}) (reply interface{}, err error) {
	params := make([]interface{}, len(vals)+1)
	params[0] = key
	for i, v := range vals {
		params[i+1] = v
	}

	reply, err = rc.do("SADD", params...)
	return
}

func (rc *Cache) SISMMBER(key string, val interface{}) (reply interface{}, err error) {
	reply, err = rc.do("SISMEMBER", key, val)
	return
}

func (rc *Cache) SMEMBERS(key string) (reply interface{}, err error) {
	reply, err = rc.do("SMEMBERS", key)
	return
}

func (rc *Cache) SREM(key string, members ...interface{}) (reply interface{}, err error) {
	params := make([]interface{}, len(members)+1)
	params[0] = key
	for i, v := range members {
		params[i+1] = v
	}
	reply, err = rc.do("SREM", params...)
	return
}

func (rc *Cache) ZADD(key string, score, member interface{}) (reply interface{}, err error) {
	reply, err = rc.do("ZADD", key, score, member)
	return
}

func (rc *Cache) ZPOPMIN(key string, count int32) (reply interface{}, err error) {

	if count == 0 {
		reply, err = rc.do("ZPOPMIN", key)
	} else {
		reply, err = rc.do("ZPOPMIN", key, count)
	}
	return
}

func (rc *Cache) ZREM(key string, members ...interface{}) (reply interface{}, err error) {
	params := []interface{}{}
	params = append(params, key)
	for _, member := range members {
		params = append(params, member)
	}
	reply, err = rc.do("ZREM", params...)
	return
}

func (rc *Cache) ZCARD(key string) (reply interface{}, err error) {
	reply, err = rc.do("ZCARD", key)
	return
}

func (rc *Cache) ZCOUNT(key string, minScore, maxScore interface{}) (reply interface{}, err error) {
	reply, err = rc.do("ZCOUNT", key, minScore, maxScore)
	return
}

func (rc *Cache) ZSCORE(key string, member interface{}) (reply interface{}, err error) {
	reply, err = rc.do("ZSCORE", key, member)
	return
}

func (rc *Cache) ZRANGE(key string, startRank, endRank int64, bWithScorse bool) (reply interface{}, err error) {
	if bWithScorse {
		reply, err = rc.do("ZRANGE", key, startRank, endRank, "WITHSCORES")
	} else {
		reply, err = rc.do("ZRANGE", key, startRank, endRank)
	}
	return
}

func (rc *Cache) ZREVRANGE(key string, startRank, endRank int64, bWithScorse bool) (reply interface{}, err error) {
	if bWithScorse {
		reply, err = rc.do("ZREVRANGE", key, startRank, endRank, "WITHSCORES")
	} else {
		reply, err = rc.do("ZREVRANGE", key, startRank, endRank)
	}
	return
}

func (rc *Cache) ZRANGEBYSCORE(key string, startScore, endScore interface{}) (reply interface{}, err error) {
	reply, err = rc.do("ZRANGEBYSCORE", key, startScore, endScore)
	return
}

func (rc *Cache) ZREVRANGEBYSCORE(key string, startScore, endScore interface{}) (reply interface{}, err error) {
	reply, err = rc.do("ZREVRANGEBYSCORE", key, startScore, endScore)
	return
}

func (rc *Cache) ZREMRANGEBYSCORE(key string, startScore, endScore interface{}) (reply interface{}, err error) {
	reply, err = rc.do("ZREMRANGEBYSCORE", key, startScore, endScore)
	return
}

func (rc *Cache) ZINCRBY(key string, increment, member interface{}) (reply interface{}, err error) {
	reply, err = rc.do("ZINCRBY", key, increment, member)
	return
}

func (rc *Cache) SINTER(key ...interface{}) (reply interface{}, err error) {
	reply, err = rc.do("SINTER", key...)
	return
}

func (rc *Cache) HGET(key string, filed string) (reply interface{}, err error) {
	reply, err = rc.do("HGET", key, filed)
	return
}

func (rc *Cache) HMSET(value []interface{}) (reply interface{}, err error) {
	reply, err = rc.do("HMSET", value...)
	return
}

func (rc *Cache) HGETALL(key string) (reply interface{}, err error) {
	reply, err = rc.do("HGETALL", key)
	return
}

func (rc *Cache) HSET(key string, filed string, value int64) (reply interface{}, err error) {
	reply, err = rc.do("HSET", key, filed, value)
	return
}

func (rc *Cache) HDEL(key string, fields []string) (reply interface{}, err error) {
	reply, err = rc.do("HDEL", key, fields)
	return
}

func (rc *Cache) HINCRBY(key string, filed string, increment int64) (reply interface{}, err error) {
	reply, err = rc.do("HINCRBY", key, filed, increment)
	return
}

func (rc *Cache) EXPIRE(key string, seconds int32) (reply interface{}, err error) {
	reply, err = rc.do("EXPIRE", key, seconds)
	return
}

func (rc *Cache) Push(key string, value string) error {
	_, err := rc.do("RPush", key, value)
	return err
}

func (rc *Cache) LLen(key string, lmtcount int64) error {
	listcnt, err := rc.do("LLen", key)
	if err != nil {
		return err
	}

	cnt := listcnt.(int64)
	if cnt > lmtcount {
		for i := lmtcount; i <= cnt; i++ {
			rc.do("LPop", key)
		}
	}

	return nil
}

// Get cache from redis.
func (rc *Cache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// GetMulti get cache from redis.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Put put cache to redis.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete delete cache in redis.
func (rc *Cache) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *Cache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (rc *Cache) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

func (rc *Cache) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

func (rc *Cache) Publich(channel string, value string) (err error) {
	_, err = rc.do("PUBLISH", channel, value)
	return err
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *Cache) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	// Format redis://<password>@<host>:<port>
	cf["conn"] = strings.Replace(cf["conn"], "redis://", "", 1)
	if i := strings.Index(cf["conn"], "@"); i > -1 {
		cf["password"] = cf["conn"][0:i]
		cf["conn"] = cf["conn"][i+1:]
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("redis", NewRedisCache)
}
