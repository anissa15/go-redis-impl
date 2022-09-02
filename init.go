package redisimpl

import (
	"context"
	"fmt"

	goredis "github.com/go-redis/redis/v8"
)

type Redis interface {
	Info
	LuaScript
}

type redis struct {
	client    *goredis.Client
	runID     string
	luascript luascript
}

func New(client *goredis.Client) Redis {
	if client == nil {
		panic(fmt.Errorf("redis client is nil"))
	}
	r := &redis{client: client}
	r.prepare()
	return r
}

func NewUrl(redisURL string) Redis {
	opt, err := goredis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Errorf("failed parse redisURL: %v", err))
	}
	client := goredis.NewClient(opt)
	r := &redis{client: client}
	r.prepare()
	return r
}

func (r *redis) prepare() {
	server, err := r.Info(context.Background(), SectionServer)
	if err != nil {
		panic(fmt.Errorf("failed get server info: %v", err))
	}
	r.runID = fmt.Sprint(server["run_id"])
	r.prepareLuaScript()
}
