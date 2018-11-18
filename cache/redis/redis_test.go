package redis

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mofancloud/xmicro/cache"
)

func TestRedisCache(t *testing.T) {
	bm, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379"}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("goods", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("goods") {
		t.Error("check err")
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist("goods") {
		t.Error("check err")
	}
	if err = bm.Put("goods", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if v, _ := redis.Int(bm.Get("goods"), err); v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("goods"); err != nil {
		t.Error("Incr Error", err)
	}

	if v, _ := redis.Int(bm.Get("goods"), err); v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("goods"); err != nil {
		t.Error("Decr Error", err)
	}

	if v, _ := redis.Int(bm.Get("goods"), err); v != 1 {
		t.Error("get err")
	}
	bm.Delete("goods")
	if bm.IsExist("goods") {
		t.Error("delete err")
	}

	//test string
	if err = bm.Put("goods", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("goods") {
		t.Error("check err")
	}

	if v, _ := redis.String(bm.Get("goods"), err); v != "author" {
		t.Error("get err")
	}

	//test GetMulti
	if err = bm.Put("goods1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("goods1") {
		t.Error("check err")
	}

	vv := bm.GetMulti([]string{"goods", "goods1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if v, _ := redis.String(vv[0], nil); v != "author" {
		t.Error("GetMulti ERROR")
	}
	if v, _ := redis.String(vv[1], nil); v != "author1" {
		t.Error("GetMulti ERROR")
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}
