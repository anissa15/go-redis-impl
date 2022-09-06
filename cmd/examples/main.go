package main

import (
	"context"
	"encoding/json"
	"fmt"

	redis "github.com/anissa15/go-redis-impl"
)

var (
	redisURL = "redis://:@localhost:26379/0"
)

func main() {
	ctx := context.Background()

	r := redis.NewUrl(redisURL)

	// checking run_id field in the INFO command
	maps, err := r.Info(ctx, redis.SectionServer)
	if err != nil {
		panic(err)
	}
	fmt.Println("run_id =", maps["run_id"])

	simpleScriptName := redis.ScriptName("simple script")
	simpleScript := `
	--[[
		this simple script do:
		- increment index key to get new index
		- set data using index
		- push index to list of index
	]]

	local indexKey = KEYS[1]
	local dataKey = KEYS[2]
	local indexListKey = KEYS[3]

	local data = cjson.decode(ARGV[1])
	
	-- uncomment below to write data in redis log
	-- redis.log(redis.LOG_NOTICE, "elements:", #data)
	
	-- get new index
	local index = redis.call("INCR", indexKey)

	-- add data
	local n = table.getn(data)
	for i=1,n,2 do
		local k = data[i]
		local v = data[i+1]
		redis.call("HSET", dataKey .. index, k, v)
	end

	-- append index to list
	redis.call("RPUSH", indexListKey, index)

	return index
	`

	err = r.RegisterScript(ctx, simpleScriptName, simpleScript)
	if err != nil {
		panic(err)
	}
	fmt.Println("success register script:", simpleScriptName)
	fmt.Println("script hash:", r.GetScriptHash(simpleScriptName))

	customerIndexKey := "index:customer"
	customerDataKey := "customer#"
	customerIndexListKey := "customer:ids"
	for i := 0; i < 5; i++ {
		keys := []string{customerIndexKey, customerDataKey, customerIndexListKey}
		customer := []string{"name", fmt.Sprintf("cust %d", i), "age", fmt.Sprint(i)}
		value, err := json.Marshal(customer)
		if err != nil {
			panic(err)
		}
		res, err := r.EvalScript(ctx, simpleScriptName, keys, value)
		if err != nil {
			panic(err)
		}
		fmt.Println("keys:", keys, "value:", res)
	}

	productIndexKey := "index:product"
	productDataKey := "product#"
	productIndexListKey := "product:ids"
	for i := 0; i < 5; i++ {
		keys := []string{productIndexKey, productDataKey, productIndexListKey}
		product := []string{"name", fmt.Sprintf("prod %d", i), "count", fmt.Sprint(i)}
		value, err := json.Marshal(product)
		if err != nil {
			panic(err)
		}
		res, err := r.EvalScript(ctx, simpleScriptName, keys, value)
		if err != nil {
			panic(err)
		}
		fmt.Println("keys:", keys, "value:", res)
	}

	fmt.Println("success")
}
