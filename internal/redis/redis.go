package redis

import (
	"context"
	"fff-hedgedoc/internal/logging"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var client *redis.Client

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
	logging.Logger.WithFields(logrus.Fields{"module": "redis", "method": "InitRedis"}).Info(fmt.Sprintf("Pinged Redis!: %s", pong))
	return nil
}
