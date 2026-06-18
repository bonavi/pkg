package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfigEnv struct {
	Host     string `env:"REDIS_HOST"`
	User     string `env:"REDIS_USER"`
	Password string `env:"REDIS_PASSWORD"`
}

func NewClientRedis(cfg RedisConfigEnv, database int) (*redis.Client, error) {

	// Создаем клиент
	client := redis.NewClient(&redis.Options{
		Network:                    "",
		Addr:                       cfg.Host,
		ClientName:                 "",
		Dialer:                     nil,
		OnConnect:                  nil,
		Protocol:                   0,
		Username:                   "",
		Password:                   cfg.Password,
		CredentialsProvider:        nil,
		CredentialsProviderContext: nil,
		DB:                         database,
		MaxRetries:                 3,
		MinRetryBackoff:            0,
		MaxRetryBackoff:            0,
		DialTimeout:                5 * time.Second,
		ReadTimeout:                5 * time.Second,
		WriteTimeout:               5 * time.Second,
		ContextTimeoutEnabled:      false,
		PoolFIFO:                   false,
		PoolSize:                   100,
		PoolTimeout:                10 * time.Second,
		MinIdleConns:               10,
		MaxIdleConns:               0,
		MaxActiveConns:             0,
		ConnMaxIdleTime:            0,
		ConnMaxLifetime:            0,
		TLSConfig:                  nil,
		Limiter:                    nil,
		DisableIndentity:           false,
		IdentitySuffix:             "",
		UnstableResp3:              false,
	})

	// Проверяем соединение
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	// Возвращаем клиент
	return client, nil
}
