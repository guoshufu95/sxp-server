package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
	"sxp-server/app/model"
	"sxp-server/common/initial"
	"sxp-server/common/logger"
	"sync"
	"time"
)

type DelayQueue struct {
	name          string
	redisCli      RedisCli
	cb            func(string) bool
	pendingKey    string // sorted set: message id -> delivery time
	readyKey      string // list
	unAckKey      string // sorted set: message id -> retry time
	retryKey      string // list
	retryCountKey string // hash: message id -> remain retry count
	garbageKey    string // set: message id
	useHashTag    bool
	ticker        *time.Ticker
	logger        *logger.ZapLog
	close         chan struct{}

	maxConsumeDuration time.Duration
	msgTTL             time.Duration
	defaultRetryCount  uint
	fetchInterval      time.Duration
	fetchLimit         uint
	concurrent         uint
}

// NilErr represents redis nil
var NilErr = errors.New("nil")

// RedisCli is abstraction for redis client, required commands only not all commands
type RedisCli interface {
	Eval(script string, keys []string, args []interface{}) (interface{}, error) // args should be string, integer or float
	Set(key string, value string, expiration time.Duration) error
	Set2(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(keys []string) error
	Del1(key string) error
	HSet(key string, field string, value string) error
	HMget(key, filed string) ([]interface{}, error)
	HDel(key string, fields []string) error
	SMembers(key string) ([]string, error)
	SRem(key string, members []string) error
	ZAdd(key string, values map[string]float64) error
	ZRem(key string, fields []string) error
}

type hashTagKeyOpt int

func UseHashTagKey() interface{} {
	return hashTagKeyOpt(1)
}

func NewQueue0(name string, cli RedisCli, callback func(string) bool, opts ...interface{}) *DelayQueue {
	if name == "" {
		panic("name is required")
	}
	if cli == nil {
		panic("cli is required")
	}
	if callback == nil {
		panic("callback is required")
	}
	useHashTag := false
	for _, opt := range opts {
		switch opt.(type) {
		case hashTagKeyOpt:
			useHashTag = true
		}
	}
	var keyPrefix string
	if useHashTag {
		keyPrefix = "{dp:" + name + "}"
	} else {
		keyPrefix = "dp:" + name
	}
	return &DelayQueue{
		name:               name,
		redisCli:           cli,
		cb:                 callback,
		pendingKey:         keyPrefix + ":pending",
		readyKey:           keyPrefix + ":ready",
		unAckKey:           keyPrefix + ":unack",
		retryKey:           keyPrefix + ":retry",
		retryCountKey:      keyPrefix + ":retry:cnt",
		garbageKey:         keyPrefix + ":garbage",
		useHashTag:         useHashTag,
		close:              make(chan struct{}, 1),
		maxConsumeDuration: 5 * time.Second,
		msgTTL:             time.Hour,
		logger:             logger.GetLogger(),
		defaultRetryCount:  3,
		fetchInterval:      time.Second,
		concurrent:         1,
	}
}

func (q *DelayQueue) WithLogger(logger *logger.ZapLog) *DelayQueue {
	q.logger = logger
	return q
}

func (q *DelayQueue) WithFetchInterval(d time.Duration) *DelayQueue {
	q.fetchInterval = d
	return q
}

func (q *DelayQueue) WithMaxConsumeDuration(d time.Duration) *DelayQueue {
	q.maxConsumeDuration = d
	return q
}

func (q *DelayQueue) WithFetchLimit(limit uint) *DelayQueue {
	q.fetchLimit = limit
	return q
}

// WithConcurrent sets the number of concurrent consumers
func (q *DelayQueue) WithConcurrent(c uint) *DelayQueue {
	if c == 0 {
		return q
	}
	q.concurrent = c
	return q
}

func (q *DelayQueue) WithDefaultRetryCount(count uint) *DelayQueue {
	q.defaultRetryCount = count
	return q
}

func (q *DelayQueue) genMsgKey(idStr string) string {
	if q.useHashTag {
		return "{dp:" + q.name + "}" + ":msg:" + idStr
	}
	return "dp:" + q.name + ":msg:" + idStr
}

type RetryCountOpt int

func WithRetryCount(count int) interface{} {
	return RetryCountOpt(count)
}

type msgTTLOpt time.Duration

func WithMsgTTL(d time.Duration) interface{} {
	return msgTTLOpt(d)
}

func (q *DelayQueue) SendScheduleMsg(payload model.TaskField, t time.Time, opts ...interface{}) error {
	// parse options
	retryCount := q.defaultRetryCount
	for _, opt := range opts {
		switch o := opt.(type) {
		case RetryCountOpt:
			retryCount = uint(o)
		case msgTTLOpt:
			q.msgTTL = time.Duration(o)
		}
	}
	// generate id
	idStr := uuid.Must(uuid.NewRandom()).String()
	v, _ := json.Marshal(payload)
	r := q.redisCli.Set2(q.genMsgKey(idStr)+"table", v, 0)
	fmt.Println(r)
	now := time.Now()
	// store msg
	value := payload.Value
	msgTTL := t.Sub(now) + q.msgTTL // delivery + q.msgTTL
	err := q.redisCli.Set(q.genMsgKey(idStr), value, msgTTL)
	if err != nil {
		return fmt.Errorf("store msg failed: %v", err)
	}
	// store retry count
	err = q.redisCli.HSet(q.retryCountKey, idStr, strconv.Itoa(int(retryCount)))
	if err != nil {
		return fmt.Errorf("store retry count failed: %v", err)
	}
	// put to pending
	err = q.redisCli.ZAdd(q.pendingKey, map[string]float64{idStr: float64(t.Unix())})
	if err != nil {
		return fmt.Errorf("push to pending failed: %v", err)
	}
	return nil
}

func (q *DelayQueue) SendDelayMsg(payload model.TaskField, duration time.Duration, opts ...interface{}) error {
	t := time.Now().Add(duration)
	return q.SendScheduleMsg(payload, t, opts...)
}

const pending2ReadyScript = `
local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get ready msg
if (#msgs == 0) then return end
local args2 = {'LPush', KEYS[2]} -- push into ready
for _,v in ipairs(msgs) do
	table.insert(args2, v) 
    if (#args2 == 4000) then
		redis.call(unpack(args2))
		args2 = {'LPush', KEYS[2]}
	end
end
if (#args2 > 2) then 
	redis.call(unpack(args2))
end
redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from pending
`

func (q *DelayQueue) pending2Ready() error {
	now := time.Now().Unix()
	keys := []string{q.pendingKey, q.readyKey}
	_, err := q.redisCli.Eval(pending2ReadyScript, keys, []interface{}{now})
	if err != nil && err != NilErr {
		return fmt.Errorf("pending2ReadyScript failed: %v", err)
	}
	return nil
}

const ready2UnackScript = `
local msg = redis.call('RPop', KEYS[1])
if (not msg) then return end
redis.call('ZAdd', KEYS[2], ARGV[1], msg)
return msg
`

func (q *DelayQueue) ready2Unack() (string, error) {
	retryTime := time.Now().Add(q.maxConsumeDuration).Unix()
	keys := []string{q.readyKey, q.unAckKey}
	ret, err := q.redisCli.Eval(ready2UnackScript, keys, []interface{}{retryTime})
	if err == NilErr {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
	}
	str, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("illegal result: %#v", ret)
	}
	return str, nil
}

func (q *DelayQueue) retry2Unack() (string, error) {
	retryTime := time.Now().Add(q.maxConsumeDuration).Unix()
	keys := []string{q.retryKey, q.unAckKey}
	ret, err := q.redisCli.Eval(ready2UnackScript, keys, []interface{}{retryTime, q.retryKey, q.unAckKey})
	if err == NilErr {
		return "", NilErr
	}
	if err != nil {
		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
	}
	str, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("illegal result: %#v", ret)
	}
	return str, nil
}

func (q *DelayQueue) callback(idStr string) error {
	payload, err := q.redisCli.Get(q.genMsgKey(idStr))
	if err == NilErr {
		return nil
	}
	if err != nil {
		// Is an IO error?
		return fmt.Errorf("get message payload failed: %v", err)
	}
	ack := q.cb(payload)
	if ack {
		err = q.ack(idStr)
	} else {
		err = q.nack(idStr)
	}
	return err
}

func (q *DelayQueue) batchCallback(ids []string) {
	if len(ids) == 1 || q.concurrent == 1 {
		for _, id := range ids {
			err := q.callback(id)
			if err != nil {
				q.logger.Errorf("consume msg %s failed: %v", id, err)
			}
		}
		return
	}
	ch := make(chan string, len(ids))
	for _, id := range ids {
		ch <- id
	}
	close(ch)
	wg := sync.WaitGroup{}
	concurrent := int(q.concurrent)
	if concurrent > len(ids) { // too many goroutines is no use
		concurrent = len(ids)
	}
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			for id := range ch {
				err := q.callback(id)
				if err != nil {
					q.logger.Errorf("consume msg %s failed: %v", id, err)
				}
			}
		}()
	}
	wg.Wait()
}

// ack
//
//	@Description: 任务执行成功的后续操作
//	@receiver q
//	@param idStr
//	@return error
func (q *DelayQueue) ack(idStr string) error {
	var (
		val   string
		err   error
		filed model.TaskField
	)
	app := initial.App
	val, err = q.redisCli.Get(q.genMsgKey(idStr) + "table")
	if err != nil {
		app.Logger.Errorf("更新task表status=1,定时任务: %s, 失败: %s", val, err.Error())
		return err
	}
	err = json.Unmarshal([]byte(val), &filed)
	if err != nil {
		app.Logger.Errorf("更新task表status=1,获取字段反序列化失败: %s", err.Error())
		return err
	}
	err = app.Db.Table("task").Debug().Where("task_name = ?", filed.TaskName).Update("status", 1).Error
	if err != nil {
		app.Logger.Errorf("更新task表staus=1 失败: %s", err.Error())
		return err
	}
	err = q.redisCli.Del1(q.genMsgKey(idStr) + "table")
	if err != nil {
		return err
	}
	err = q.redisCli.ZRem(q.unAckKey, []string{idStr})
	if err != nil {
		return fmt.Errorf("remove from unack failed: %v", err)
	}
	// msg key has ttl, ignore result of delete
	_ = q.redisCli.Del([]string{q.genMsgKey(idStr)})
	_ = q.redisCli.HDel(q.retryCountKey, []string{idStr})
	return nil
}

func (q *DelayQueue) nack(idStr string) error {
	// update retry time as now, unack2Retry will move it to retry immediately
	res, err := q.redisCli.HMget(q.retryCountKey, idStr)
	if res[0].(string) == "0" {

	}
	//var t time.Duration
	//if res[0].(string) == "3" {
	//	t = 5 * time.Second
	//}
	//if res[0].(string) == "2" {
	//	t = 10 * time.Second
	//}
	//if res[0].(string) == "1" {
	//	t = 15 * time.Second
	//}
	err = q.redisCli.ZAdd(q.unAckKey, map[string]float64{
		idStr: float64(time.Now().Add(3 * time.Second).Unix()),
	})
	if err != nil {
		return fmt.Errorf("negative ack failed: %v", err)
	}
	return nil
}

const unack2RetryScript = `
local unack2retry = function(msgs)
	local retryCounts = redis.call('HMGet', KEYS[2], unpack(msgs)) -- get retry count
	for i,v in ipairs(retryCounts) do
		local k = msgs[i]
		if v ~= false and v ~= nil and v ~= '' and tonumber(v) > 0 then
			redis.call("HIncrBy", KEYS[2], k, -1) -- reduce retry count
			redis.call("LPush", KEYS[3], k) -- add to retry
		else
			redis.call("HDel", KEYS[2], k) -- del retry count
			redis.call("SAdd", KEYS[4], k) -- add to garbage
		end
	end
end

local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get retry msg
if (#msgs == 0) then return end
if #msgs < 4000 then
	unack2retry(msgs)
else
	local buf = {}
	for _,v in ipairs(msgs) do
		table.insert(buf, v)
		if #buf == 4000 then
			unack2retry(buf)
			buf = {}
		end
	end
	if (#buf > 0) then
		unack2retry(buf)
	end
end
redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from unack
`

func (q *DelayQueue) unack2Retry() error {
	keys := []string{q.unAckKey, q.retryCountKey, q.retryKey, q.garbageKey}
	now := time.Now()
	_, err := q.redisCli.Eval(unack2RetryScript, keys, []interface{}{now.Unix()})
	if err != nil && err != NilErr {
		return fmt.Errorf("unack to retry script failed: %v", err)
	}
	return nil
}

// garbageCollect
//
//	@Description: 任务执行失败的后续操作
//	@receiver q
//	@return error
func (q *DelayQueue) garbageCollect() error {
	msgIds, err := q.redisCli.SMembers(q.garbageKey)
	if err != nil {
		return fmt.Errorf("smembers failed: %v", err)
	}
	if len(msgIds) == 0 {
		return nil
	}

	//更新数据库task表的status
	for _, k := range msgIds {
		var (
			val   string
			filed model.TaskField
		)
		app := initial.App
		val, err = q.redisCli.Get(q.genMsgKey(k) + "table")
		if err != nil {
			app.Logger.Errorf("更新task表,定时任务: %s, 失败: %s", val, err.Error())
			continue
		}
		err = json.Unmarshal([]byte(val), &filed)
		if err != nil {
			app.Logger.Errorf("更新task表,获取字段反序列化失败: %s", err.Error())
			continue
		}
		err = app.Db.Table("task").Debug().Where("task_name = ?", filed.TaskName).Update("status", 2).Error
		if err != nil {
			app.Logger.Errorf("更新task表失败: %s", err.Error())
			continue
		}
		q.redisCli.Del1(q.genMsgKey(k) + "table")
	}

	// allow concurrent clean
	msgKeys := make([]string, 0, len(msgIds))
	for _, idStr := range msgIds {
		msgKeys = append(msgKeys, q.genMsgKey(idStr))
	}
	err = q.redisCli.Del(msgKeys)
	if err != nil && err != NilErr {
		return fmt.Errorf("del msgs failed: %v", err)
	}
	err = q.redisCli.SRem(q.garbageKey, msgIds)
	if err != nil && err != NilErr {
		return fmt.Errorf("remove from garbage key failed: %v", err)
	}
	return nil
}

func (q *DelayQueue) consume() error {
	// pending to ready
	err := q.pending2Ready()
	if err != nil {
		return err
	}
	// consume
	ids := make([]string, 0, q.fetchLimit)
	for {
		idStr, err := q.ready2Unack()
		if err == NilErr { // consumed all
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, idStr)
		if q.fetchLimit > 0 && len(ids) >= int(q.fetchLimit) {
			break
		}
	}
	if len(ids) > 0 {
		q.batchCallback(ids)
	}
	// unack to retry
	err = q.unack2Retry()
	if err != nil {
		return err
	}
	err = q.garbageCollect()
	if err != nil {
		return err
	}
	// retry
	ids = make([]string, 0, q.fetchLimit)
	for {
		idStr, err := q.retry2Unack()
		if errors.Is(err, NilErr) { // consumed all
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, idStr)
		if q.fetchLimit > 0 && len(ids) >= int(q.fetchLimit) {
			break
		}
	}
	if len(ids) > 0 {
		q.batchCallback(ids)
	}
	return nil
}

func (q *DelayQueue) StartConsume() (done <-chan struct{}) {
	q.ticker = time.NewTicker(q.fetchInterval)
	done0 := make(chan struct{})
	go func() {
	tickerLoop:
		for {
			select {
			case <-q.ticker.C:
				err := q.consume()
				if err != nil {
					log.Printf("consume error: %v", err)
				}
			case <-q.close:
				fmt.Println("##########定时任务队列消费者退出###########")
				break tickerLoop
			}
		}
		close(done0)
	}()
	return done0
}

// StopConsume stops consumer goroutine
func (q *DelayQueue) StopConsume() {
	fmt.Println("定时任务停止")
	close(q.close)
	if q.ticker != nil {
		q.ticker.Stop()
	}
}
