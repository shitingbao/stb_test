package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func luaLoad() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0, // 默认DB 0
	})

	defer rdb.Close()
	ctx := context.Background()
	// 普通执行
	// res, err := rdb.Eval(ctx, "return {KEYS[1],ARGV[1]}", []string{"name"}, "stb").Result()
	// 多行脚本执行

	res, err := rdb.Eval(ctx, `
	local key1 = KEYS[1]
	local key2 = KEYS[2]
	local val1 = ARGV[1]
	local val2 = ARGV[2]
	
	local res = redis.call("get",key1)
	if (res == 18) 
	then
		redis.call("set",key1,val1)
	else
		redis.call("set",key1,10)
	end
	
	return redis.call("get",key2)
	`, []string{"name", "age"}, "stb", 18).Result()
	if err != nil {
		return
	}

	log.Println("res=:", res)
}

// 预加载脚本
func loadHas(ctx context.Context, rdb *redis.Client) {
	rdb.EvalSha(ctx, "", []string{}).Result()
	rdb.ScriptLoad(ctx, "")
}

var (
	// 对应的可使用种草任务数量减1
	CompanyTaskDeductionLua = `
	local test = tonumber(redis.call('GET', KEYS[1]))
	if test and test > 0 then
		redis.call('DECR', KEYS[1])
		return 1
	else
		return 0
	end`

	CompanyTaskDeductionLuaHas = ""
)

// address: "118.31.12.177:6379"
// db: 1
// pass: "Awjgrayredis!@#$%"

func RedisLoad() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		// Addr:     "47.99.104.79:6379",
		// Addr:     "127.0.0.1:6379",
		// Addr:     "10.176.107.70:6379",
		// Password: "12345678",

		Addr: "118.31.12.177:6379",
		// Addr:     "r-bp1us08xv15vg063ex.tairpena.rds.aliyuncs.com:6379",
		Password: "Awjgrayredis!@#$%",

		// Addr:     "47.97.176.134:6379",
		// Password: "Awjtestredis!@#$%",
		DB: 1, // 默认DB 0
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Println("redis clietn:", err)
		return nil, err
	}

	luaLoad()
	loadHas(ctx, rdb)
	redisLoad()
	return rdb, nil
}

func LoadLua(ctx context.Context, rdb *redis.Client) error {
	exist, err := rdb.ScriptExists(ctx, CompanyTaskDeductionLuaHas).Result()

	if err != nil {
		return err
	}
	log.Println("exists:", exist)
	if len(exist) > 0 && exist[0] {
		return errors.New("exist")
	}

	res, err := rdb.ScriptLoad(ctx, CompanyTaskDeductionLua).Result()

	if err != nil {
		return err
	}
	CompanyTaskDeductionLuaHas = res

	log.Println("CompanyTaskDeductionLuaHas:", CompanyTaskDeductionLuaHas)
	return nil
}

func EvalSha(ctx context.Context, rdb *redis.Client) error {
	res, err := rdb.EvalSha(ctx, "3d838308c0c4395241fbfaa178a2002cfac8d978", []string{"task1"}).Result()

	if err != nil {
		return err
	}

	b, err := json.Marshal(res)

	if err != nil {
		return err
	}

	log.Println("lua:", err, string(b))
	log.Println("CompanyTaskDeductionLuaHas:", CompanyTaskDeductionLuaHas)
	return nil
}

func redisLoad() error {
	ctx := context.Background()
	rdb, err := RedisLoad()

	if err != nil {
		return err
	}
	defer rdb.Close()

	for i := 1; i < 10; i++ {
		keyLock := fmt.Sprintf("context:select:task:lock:%d", i)
		key := fmt.Sprintf("context:select:task:%d", i)
		// res, err := rdb.XInfoGroups(ctx, "stb").Result()
		lockRes, err := rdb.Get(ctx, keyLock).Result()
		if err != nil {
			return err
		}
		res, err := rdb.Get(ctx, key).Result()
		// rdb.Set(ctx, key, 2, -1).Result()
		// rdb.Expire(ctx, keyLock, time.Second)
		// res, err := rdb.XGroupCreateMkStream(ctx, "stream_stb", "stream_stb_group", "0").Result()

		// args := &redis.XAddArgs{
		// 	Stream: "stream_stb",
		// }
		// res, err := rdb.XAdd(ctx, args).Result()
		log.Println("i:", i, "-lock:", lockRes, "-res:", res, err)
	}

	// if err := LoadLua(ctx, rdb); err != nil {
	// 	log.Println(err)
	// }

	// if err := LoadLua(ctx, rdb); err != nil {
	// 	log.Println(err)
	// }

	// if err := EvalSha(ctx, rdb); err != nil {
	// 	log.Println(err)
	// }

	return nil
}
