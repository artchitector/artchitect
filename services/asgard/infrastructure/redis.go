package infrastructure

import "github.com/go-redis/redis/v8"

func InitRedis(config *Config) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword,
			DB:       0,
		},
	)
}
