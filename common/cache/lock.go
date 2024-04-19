package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync/atomic"
	"time"
)

const RedisLockKeyPrefix = "REDIS_LOCK_PREFIX_"

var ErrLockAcquiredByOthers = errors.New("lock is acquired by others")

// RedisLock
// @Description: redisLock
type RedisLock struct {
	key           string
	token         string
	Client        *redis.Client
	expireSeconds int64 //key过期时间
	isBlock       bool  //是否阻塞
	blockTimeOut  int64 //轮询超时时间
	MaxIdle       int64 //最大消费时间
	watch         bool  //是否启用看门狗
	// 看门狗标识
	runningDog int32
	// 停止看门狗
	stopDog context.CancelFunc
}

type Options func(r *RedisLock)

// NewLock
//
//	@Description: 返回一个自己定义的lock结构体，可调用其实现的方法
//	@param key
//	@param token
//	@param client
//	@param opts
//	@return *RedisLock
func NewLock(key, token string, client *redis.Client, opts ...Options) *RedisLock {
	r := RedisLock{
		key:           key,
		token:         token, //唯一token
		Client:        client,
		expireSeconds: 10,   //默认锁30秒过期
		isBlock:       true, //是否阻塞
		blockTimeOut:  5,    //默认阻塞时间为5秒
		watch:         true, //启动看门狗
	}
	for _, opt := range opts {
		opt(&r)
	}
	return &r
}

func WithExpire(i int64) Options {
	return func(r *RedisLock) {
		r.expireSeconds = i
	}
}

func WithIsBlock(block bool) Options {
	return func(r *RedisLock) {
		r.isBlock = block
	}
}

func WithWatch(watch bool) Options {
	return func(r *RedisLock) {
		r.watch = watch
	}
}

func WithBlockTimeOut(i int64) Options {
	return func(r *RedisLock) {
		r.blockTimeOut = i
	}
}

func WithMaxIdle(i int64) Options {
	return func(r *RedisLock) {
		r.MaxIdle = i
	}
}

func IsRetryableErr(err error) bool {
	return errors.Is(err, ErrLockAcquiredByOthers)
}

// Lock
//
//	@Description: 上锁
//	@receiver r
//	@param ctx
//	@return err
func (r *RedisLock) Lock(ctx context.Context) (err error) {
	defer func() {
		// 如果加锁成功，启动一个看门狗程序
		r.watchDog(ctx)
	}()
	fmt.Println(fmt.Sprintf("%s加锁", r.token))
	err = r.TryLock(ctx)
	if err == nil {
		return
	}
	if !r.isBlock {
		return err
	}
	// 判断错误类型
	if !IsRetryableErr(err) {
		return err
	}
	// 持续轮询取锁
	err = r.blockingLock(ctx)
	return
}

// Unlock
//
//	@Description: 解锁
//	@receiver r
//	@param ctx
//	@return err
func (r *RedisLock) Unlock(ctx context.Context) (err error) {
	keys := []string{r.getLockKey()}
	args := []interface{}{r.token}
	reply, err := r.Client.Eval(ctx, DeleteDistributionLock, keys, args).Result()
	if err != nil {
		return err
	}
	if ret, _ := reply.(int64); ret != 1 {
		return errors.New("删除锁失败")
	}
	r.stopDog()
	fmt.Println(fmt.Sprintf("%s:解锁成功", r.token))
	return nil
}

// TryLock
//
//	@Description: 加锁
//	@receiver r
//	@param ctx
//	@return err
func (r *RedisLock) TryLock(ctx context.Context) (err error) {
	flag, err := r.Client.SetNX(ctx, r.getLockKey(), r.token, time.Duration(r.expireSeconds)*time.Second).Result()
	if err != nil {
		return err
	}
	if !flag {
		err = ErrLockAcquiredByOthers
		return
	}
	fmt.Println(fmt.Sprintf("%s加锁成功", r.token))
	return
}

// blockingLock
//
//	@Description: 轮询获取锁
//	@receiver r
//	@param ctx
//	@return err
func (r *RedisLock) blockingLock(ctx context.Context) (err error) {
	fmt.Println(fmt.Sprintf("%s:轮询中", r.token))
	// 每隔 50 ms 尝试取锁一次
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("lock failed, ctx timeout, err: %w", ctx.Err())
			return
		default:
		}
		// 继续上锁
		err = r.TryLock(ctx)
		if !IsRetryableErr(err) {
			return
		}
		// 获取锁成功直接返回
		if err == nil {
			return
		}
	}
	return
}

// watchDog
//
//	@Description: 看门狗启动
//	@receiver r
//	@param ctx
func (r *RedisLock) watchDog(ctx context.Context) {
	// 是否启动看门狗
	if !r.watch {
		return
	}
	if !atomic.CompareAndSwapInt32(&r.runningDog, 0, 1) {
		return
	}
	//ctx, r.stopDog = context.WithCancel(ctx)
	go func() {
		defer func() {
			atomic.StoreInt32(&r.runningDog, 0)
		}()
		ctx, r.stopDog = context.WithCancel(ctx)
		r.runWatchDog(ctx)
	}()
}

// runWatchDog
//
//	@Description: 看门狗续期
//	@receiver r
//	@param ctx
func (r *RedisLock) runWatchDog(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			fmt.Println(fmt.Sprintf("%s:看门狗停止续期", r.token))
			return
		default:
		}
		fmt.Println(fmt.Sprintf("%s:看门狗续期", r.token))
		// 看门狗持续为分布式锁进行续期
		_ = r.DelayExpire(ctx, 3)
	}
}

// DelayExpire
//
//	@Description: 通过lua脚本实现判断和续期原子性操作
//	@receiver r
//	@param ctx
//	@param expireSeconds
//	@return error
func (r *RedisLock) DelayExpire(ctx context.Context, expireSeconds int64) error {
	keys := []string{r.getLockKey()}
	args := []interface{}{r.token, expireSeconds}
	reply, err := r.Client.Eval(ctx, ExpireDistributionLock, keys, args).Result()
	if err != nil {
		return err
	}

	if ret, _ := reply.(int64); ret != 1 {
		return errors.New("锁续期失败")
	}

	return nil
}

func (r *RedisLock) getLockKey() string {
	return RedisLockKeyPrefix + r.key
}
