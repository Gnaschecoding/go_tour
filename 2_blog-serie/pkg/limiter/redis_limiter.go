package limiter

import (
	"Golang_Programming_Journey/2_blog-serie/global"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// RedisLimiter 基于 Redis 的令牌桶限流器
type RedisLimiter struct {
	client *redis.Client
	prefix string
}

// NewRedisLimiter 创建一个新的 Redis 限流器实例

func NewRedisLimiter(client *redis.Client, prefix string) *RedisLimiter {
	return &RedisLimiter{client, prefix}
}

// Key 生成限流器的键
// 每个 IP 对每个 URL 都会限流。
func (l *RedisLimiter) Key(c *gin.Context) string {
	ip := c.ClientIP()
	uri := c.Request.RequestURI
	index := strings.Index(uri, "?")
	if index != -1 {
		uri = uri[:index] // 去掉查询参数
	}

	// 按 IP + URL 生成限流 Key,对IP进行限制的话 可能会导致Redis内存激增，所以需要进行定期过期Redis中的数据
	return l.prefix + ip + ":" + uri
}

// GetBucket 获取令牌桶信息,返回值 有 *ratelimit.Bucket完全是为了兼容 这个接口LimiterIface
func (l *RedisLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	// 这里不直接返回 ratelimit.Bucket，而是使用 Redis 操作来模拟令牌桶
	ctx := context.Background()
	_, err := l.client.Get(ctx, key).Int64()
	if err != nil {
		return nil, false
	}

	return nil, true
}

// AddBuckets 添加限流规则

func (l *RedisLimiter) AddBuckets(rules ...LimiterBucketRule) LimiterIface {
	for _, rule := range rules {
		// 初始化 Redis 中的令牌桶
		key := rule.Key
		err := l.client.Set(context.Background(), key, rule.Capacity, global.LimiterSetting.Expiration*time.Second).Err()
		if err != nil {
			return nil
		}

		// 设置令牌填充间隔 rule.Quantum
		go l.fillTokens(key, rule.FillInterval, rule.Quantum, rule.Capacity)
	}
	return l
}

// fillTokens 定时填充令牌
func (l *RedisLimiter) fillTokens(key string, interval time.Duration, quantum, capacity int64) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	luaScript := `
			local tokens = redis.call('GET', KEYS[1])   -- 获取当前令牌数
			if not tokens then                          -- 如果不存在，初始化为 0
				tokens = 0
			else
				tokens = tonumber(tokens)               -- 转换为数字
			end
			
			tokens = tokens + tonumber(ARGV[1])         -- 增加令牌数量
			
			if tokens > tonumber(ARGV[2]) then          -- 超过容量则限制
				tokens = tonumber(ARGV[2])
			end
			
			redis.call('SET', KEYS[1], tokens, 'EX', tonumber(ARGV[3]))          -- 写回 Redis
			return tokens
    `

	for range ticker.C {
		ctx := context.Background()

		// 使用 Lua 脚本确保原子性增加令牌且不会超出容量
		_, err := l.client.Eval(ctx, luaScript, []string{key}, quantum, capacity, global.LimiterSetting.Expiration*time.Second).Result()
		if err != nil {
			// 记录日志，防止 Redis 故障时出错
			global.Logger.Errorf(ctx, "Failed to execute luaScript script err:%v", err)
		}
	}
}

// TakeToken 尝试获取一个令牌
func (l *RedisLimiter) TakeToken(key string) bool {
	ctx := context.Background()
	// 使用 Lua 脚本原子性地检查和减少令牌数量
	script := `
       		 	-- 从 Redis 中获取指定键（KEYS[1]）对应的值，并将其转换为数字类型
				-- KEYS 是 Redis 执行 Lua 脚本时传入的键名数组，KEYS[1] 表示第一个键
		local tokens = tonumber(redis.call('get', KEYS[1]))
		redis.log(redis.LOG_NOTICE, "tokens value: "..tostring(tokens))
				-- 检查获取到的令牌数量是否存在且大于 0
				-- 如果 tokens 不为 nil 且大于 0，说明令牌桶中有可用令牌
		if tokens and tokens > 0 then
    			-- 若有可用令牌，调用 Redis 的 DECR 命令将该键对应的值减 1
    			-- 这一步相当于消耗掉一个令牌
			redis.call('decr', KEYS[1])
			redis.call('expire', KEYS[1], tonumber(ARGV[1]))			
				-- 返回 1 表示成功获取到一个令牌
			return 1
		end
				-- 如果没有可用令牌，返回 0 表示获取令牌失败
		return 0
    `
	result, err := l.client.Eval(ctx, script, []string{key}, global.LimiterSetting.Expiration*time.Second).Int64()

	if err != nil {
		return false
	}

	return result == 1
}
