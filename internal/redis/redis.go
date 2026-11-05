package redis

import (
	"context"
	"errors"
	"fff-hedgedoc/internal/logging"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var client *redis.Client
var RedisInitialized bool = false

const GROUPS_KEY_PREFIX = "csh:groups:"

func InitRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}

	newClient := redis.NewClient(opt)
	pong, err := newClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	client = newClient
	RedisInitialized = true
	logging.Logger.WithFields(logrus.Fields{"module": "redis", "method": "InitRedis"}).Info(fmt.Sprintf("Pinged Redis!: %s", pong))
	return nil
}

func GetAllFromKeyPrefix(ctx context.Context, prefix string, maxSize int) (map[string]string, error) {
	if !RedisInitialized {
		return nil, errors.New("redis was not initialized")
	}

	var (
		cursor      uint64
		pattern     = prefix + "*"
		result      = make(map[string]string)
		fullResults = 0
	)

	if maxSize == -1 {
		fullResults = 99999999
	} else {
		fullResults = int(maxSize)
	}
	for {
		keys, next, err := client.Scan(ctx, cursor, pattern, 500).Result()
		if err != nil {
			return nil, err
		}

		if len(keys) > 0 {
			vals, err := client.MGet(ctx, keys...).Result()
			if err != nil {
				return nil, err
			}

			for i, v := range vals {
				if s, ok := v.(string); ok {
					result[keys[i]] = s
				}
			}
		}

		cursor = next
		if cursor == 0 || len(result) >= fullResults {
			break
		}
	}
	return result, nil
}

func ResetKeysWithPrefix(ctx context.Context, prefix string) error {
	if !RedisInitialized {
		return errors.New("redis was not initialized")
	}

	var (
		cursor  uint64
		pattern = prefix + "*"
	)

	for {
		keys, next, err := client.Scan(ctx, cursor, pattern, 500).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			_, err := client.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

func SetKeyValueStore(ctx context.Context, key, value string, ttl time.Duration) error {
	if !RedisInitialized {
		return errors.New("redis was not initialized")
	}
	return client.Set(ctx, key, value, ttl).Err()
}
