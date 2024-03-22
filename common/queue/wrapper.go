package queue

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewQueue(name string, cli *redis.Client, callback func(string) bool, opts ...interface{}) *DelayQueue {
	rc := &redisV9Wrapper{
		inner: cli,
	}
	return NewQueue0(name, rc, callback, opts...)
}

func wrapErr(err error) error {
	if errors.Is(err, redis.Nil) {
		return NilErr
	}
	return err
}

type redisV9Wrapper struct {
	inner *redis.Client
}

func (r *redisV9Wrapper) Eval(script string, keys []string, args []interface{}) (interface{}, error) {
	ctx := context.Background()
	ret, err := r.inner.Eval(ctx, script, keys, args...).Result()
	return ret, wrapErr(err)
}

func (r *redisV9Wrapper) Set(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return wrapErr(r.inner.Set(ctx, key, value, expiration).Err())
}

func (r *redisV9Wrapper) Set2(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return wrapErr(r.inner.Set(ctx, key, value, expiration).Err())
}

func (r *redisV9Wrapper) Get(key string) (string, error) {
	ctx := context.Background()
	ret, err := r.inner.Get(ctx, key).Result()
	return ret, wrapErr(err)
}

func (r *redisV9Wrapper) Del(keys []string) error {
	ctx := context.Background()
	return wrapErr(r.inner.Del(ctx, keys...).Err())
}

func (r *redisV9Wrapper) Del1(key string) error {
	ctx := context.Background()
	return wrapErr(r.inner.Del(ctx, key).Err())
}

func (r *redisV9Wrapper) HMget(key, filed string) ([]interface{}, error) {
	ctx := context.Background()
	res, err := r.inner.HMGet(ctx, key, filed).Result()
	return res, wrapErr(err)
}

func (r *redisV9Wrapper) HSet(key string, field string, value string) error {
	ctx := context.Background()
	return wrapErr(r.inner.HSet(ctx, key, field, value).Err())
}

func (r *redisV9Wrapper) HDel(key string, fields []string) error {
	ctx := context.Background()
	return wrapErr(r.inner.HDel(ctx, key, fields...).Err())
}

func (r *redisV9Wrapper) SMembers(key string) ([]string, error) {
	ctx := context.Background()
	ret, err := r.inner.SMembers(ctx, key).Result()
	return ret, wrapErr(err)
}

func (r *redisV9Wrapper) SRem(key string, members []string) error {
	ctx := context.Background()
	members2 := make([]interface{}, len(members))
	for i, v := range members {
		members2[i] = v
	}
	return wrapErr(r.inner.SRem(ctx, key, members2...).Err())
}

func (r *redisV9Wrapper) ZAdd(key string, values map[string]float64) error {
	ctx := context.Background()
	var zs []redis.Z
	for member, score := range values {
		zs = append(zs, redis.Z{
			Score:  score,
			Member: member,
		})
	}
	return wrapErr(r.inner.ZAdd(ctx, key, zs...).Err())
}

func (r *redisV9Wrapper) ZRem(key string, members []string) error {
	ctx := context.Background()
	members2 := make([]interface{}, len(members))
	for i, v := range members {
		members2[i] = v
	}
	return wrapErr(r.inner.ZRem(ctx, key, members2...).Err())
}

type redisClusterWrapper struct {
	inner *redis.ClusterClient
}

func (r *redisClusterWrapper) Eval(script string, keys []string, args []interface{}) (interface{}, error) {
	ctx := context.Background()
	ret, err := r.inner.Eval(ctx, script, keys, args...).Result()
	return ret, wrapErr(err)
}

func (r *redisClusterWrapper) Set(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return wrapErr(r.inner.Set(ctx, key, value, expiration).Err())
}

func (r *redisClusterWrapper) Get(key string) (string, error) {
	ctx := context.Background()
	ret, err := r.inner.Get(ctx, key).Result()
	return ret, wrapErr(err)
}

func (r *redisClusterWrapper) Del(keys []string) error {
	ctx := context.Background()
	return wrapErr(r.inner.Del(ctx, keys...).Err())
}

func (r *redisClusterWrapper) HSet(key string, field string, value string) error {
	ctx := context.Background()
	return wrapErr(r.inner.HSet(ctx, key, field, value).Err())
}

func (r *redisClusterWrapper) HDel(key string, fields []string) error {
	ctx := context.Background()
	return wrapErr(r.inner.HDel(ctx, key, fields...).Err())
}

func (r *redisClusterWrapper) SMembers(key string) ([]string, error) {
	ctx := context.Background()
	ret, err := r.inner.SMembers(ctx, key).Result()
	return ret, wrapErr(err)
}

func (r *redisClusterWrapper) SRem(key string, members []string) error {
	ctx := context.Background()
	members2 := make([]interface{}, len(members))
	for i, v := range members {
		members2[i] = v
	}
	return wrapErr(r.inner.SRem(ctx, key, members2...).Err())
}

func (r *redisClusterWrapper) ZAdd(key string, values map[string]float64) error {
	ctx := context.Background()
	var zs []redis.Z
	for member, score := range values {
		zs = append(zs, redis.Z{
			Score:  score,
			Member: member,
		})
	}
	return wrapErr(r.inner.ZAdd(ctx, key, zs...).Err())
}

func (r *redisClusterWrapper) ZRem(key string, members []string) error {
	ctx := context.Background()
	members2 := make([]interface{}, len(members))
	for i, v := range members {
		members2[i] = v
	}
	return wrapErr(r.inner.ZRem(ctx, key, members2...).Err())
}
