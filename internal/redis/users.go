package redis

import (
	"context"
	"fff-hedgedoc/internal/models"

	"github.com/redis/go-redis/v9"
)

const USERS_KEY_PREFIX = "csh:users"
const USER_LIST_KEY_PREFIX = USERS_KEY_PREFIX + ":list"

const USERNAME_TO_ID_KEY = USERS_KEY_PREFIX + ":usertoid:"

const MAX_USERS_PER_REQ = 50

func AddUser(ctx context.Context, uuid string, username string) error {
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		if err := pipe.Set(ctx, USERNAME_TO_ID_KEY+username, uuid, -1).Err(); err != nil {
			return err
		}

		if err := pipe.ZAdd(ctx, USER_LIST_KEY_PREFIX, redis.Z{Score: 0, Member: username}).Err(); err != nil {
			return err
		}

		return nil
	})
	return err
}

func GetUsersSection(ctx context.Context, prefix string, cursor *string) ([]models.UserInfo, error) {
	lo := "[" + prefix
	if cursor != nil {
		lo = "(" + *cursor
	}
	hi := "[" + prefix + "\xff"

	usernames, err := client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   USER_LIST_KEY_PREFIX,
		Start: lo,
		Stop:  hi,
		Count: MAX_USERS_PER_REQ,
		ByLex: true,
	}).Result()
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(usernames))
	for i, name := range usernames {
		keys[i] = USERNAME_TO_ID_KEY + name
	}

	uuids, err := client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	users := make([]models.UserInfo, len(usernames))
	for index, user := range usernames {
		if uuids[index] != nil {
			users[index] = models.UserInfo{
				UUID:     uuids[index].(string),
				Username: user,
			}
		}

	}

	return users, nil
}
