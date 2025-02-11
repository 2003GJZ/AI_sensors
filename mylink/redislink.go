package mylink

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var REDIS_JUST_ONCE bool

type RedisLink struct {
	Ctx    context.Context
	Client *redis.Client
}

func NewRedisLink(databasenum int) (*RedisLink, error) {

	// 创建 Redis 客户端
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       databasenum, // 默认数据库为 0
		//// 设置连接池参数
		//PoolSize:     100000,         // 连接池最大连接数
		//MinIdleConns: 2000,         // 最小空闲连接数
		//DialTimeout:  5 * time.Second, // 连接超时时间
		//ReadTimeout:  3 * time.Second, // 读取超时时间
		//WriteTimeout: 3 * time.Second, // 写入超时时间
		//PoolTimeout:  0,               // 等待可用连接的最大时间，0表示无限等待
		PoolSize:     1000, // 非常大的连接池大小
		MinIdleConns: 100,  // 最小空闲连接数
		MaxIdleConns: 1000, // 最大空闲连接数
		PoolTimeout:  0,    // 等待可用连接的最大时间
	})
	// 测试链接
	if !REDIS_JUST_ONCE {
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			return nil, err
		}
		log.Println("Redis 链接成功")
		REDIS_JUST_ONCE = true
	}

	//创建链接不测试链接
	redislink := &RedisLink{
		Ctx:    ctx,
		Client: rdb,
	}

	return redislink, nil
}

func GetredisLink() (*RedisLink, error) {
	return NewRedisLink(0)
}
