package redisimpl

import (
	"context"
	"fmt"
)

// based on
// https://redis.io/docs/manual/programmability/eval-intro/

type LuaScript interface {
	RegisterScript(ctx context.Context, scriptName ScriptName, script string) error
	EvalScript(ctx context.Context, scriptName ScriptName, keys []string, args ...interface{}) (interface{}, error)
	GetScriptHash(scriptName ScriptName) ScriptHash
}

type ScriptName string
type ScriptHash string

type luascript struct {
	scripts map[ScriptName]ScriptHash
}

func (r *redis) prepareLuaScript() {
	r.luascript.scripts = make(map[ScriptName]ScriptHash)
}

func (r *redis) RegisterScript(ctx context.Context, scriptName ScriptName, script string) error {
	hash, err := r.client.ScriptLoad(ctx, script).Result()
	if err != nil {
		return fmt.Errorf("failed to script load: %v", err)
	}
	r.luascript.scripts[scriptName] = ScriptHash(hash)
	return nil
}

func (r *redis) GetScriptHash(scriptName ScriptName) ScriptHash {
	return r.luascript.scripts[scriptName]
}

func (r *redis) EvalScript(ctx context.Context, scriptName ScriptName, keys []string, args ...interface{}) (interface{}, error) {
	hash := r.luascript.scripts[scriptName]
	return r.client.EvalSha(ctx, string(hash), keys, args...).Result()
}
