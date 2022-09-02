package main

import (
	"context"
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
	local indexKey = KEYS[1]
	
	-- get new index
	local index = redis.call("INCR", indexKey)

	return index
	`

	err = r.RegisterScript(ctx, simpleScriptName, simpleScript)
	if err != nil {
		panic(err)
	}
	fmt.Println("success register script:", simpleScriptName)
	fmt.Println("script hash:", r.GetScriptHash(simpleScriptName))

	customerIndexKey := "index:customer"
	for i := 0; i < 5; i++ {
		// only 1 key for simpleScript
		keys := []string{customerIndexKey}
		res, err := r.EvalScript(ctx, simpleScriptName, keys)
		if err != nil {
			panic(err)
		}
		fmt.Println("keys:", keys, "value:", res)
	}

	productIndexKey := "index:product"
	for i := 0; i < 5; i++ {
		// only 1 key for simpleScript
		keys := []string{productIndexKey}
		res, err := r.EvalScript(ctx, simpleScriptName, keys)
		if err != nil {
			panic(err)
		}
		fmt.Println("keys:", keys, "value:", res)
	}

	fmt.Println("success")
}
