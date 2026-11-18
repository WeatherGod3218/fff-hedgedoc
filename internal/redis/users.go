package redis

import (
	"context"
	"errors"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

const USERS_KEY_PREFIX = "csh:users"
const USER_LIST_KEY_PREFIX = USERS_KEY_PREFIX + ":list"

const USERNAME_TO_ID_KEY = USERS_KEY_PREFIX + ":usertoid:"

const MAX_USERS_PER_REQ = 50

func AddUsers(ctx context.Context, users []models.KeycloakUser) error {
	if !RedisInitialized {
		return errors.New("redis was not initialized")
	}

	if len(users) == 0 {
		return nil
	}

	pipe := client.Pipeline()
	members := make([]redis.Z, 0, len(users))

	for _, user := range users {
		pipe.Set(ctx, USERNAME_TO_ID_KEY+user.Username, user.ID, 0)
		members = append(members, redis.Z{Score: 0, Member: strings.ToLower(user.Username)})
	}

	pipe.ZAdd(ctx, USER_LIST_KEY_PREFIX, members...)

	_, err := pipe.Exec(ctx)
	return err
}

func AddUser(ctx context.Context, uuid string, username string) error {
	if !RedisInitialized {
		return errors.New("redis was not initialized")
	}

	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		if err := pipe.Set(ctx, USERNAME_TO_ID_KEY+username, uuid, -1).Err(); err != nil {
			return err
		}

		if err := pipe.ZAdd(ctx, USER_LIST_KEY_PREFIX, redis.Z{Score: 0, Member: strings.ToLower(username)}).Err(); err != nil {
			return err
		}

		return nil
	})
	return err
}

func GetUsersSection(ctx context.Context, prefix string, cursor *string) ([]models.UserInfo, error) {
	if !RedisInitialized {
		return nil, errors.New("redis was not initialized")
	}

	lo := "[" + prefix
	if cursor != nil && strings.TrimSpace(*cursor) != "" {
		if cursor != nil && strings.HasPrefix(*cursor, prefix) {
			lo = "(" + *cursor
		}
	}
	hi := "[" + prefix + "\xff"

	usernames, err := client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:    USER_LIST_KEY_PREFIX,
		Start:  lo,
		Stop:   hi,
		Count:  MAX_USERS_PER_REQ,
		Offset: 0,
		ByLex:  true,
	}).Result()
	if err != nil {
		return nil, err
	}

	if len(usernames) == 0 {
		return []models.UserInfo{}, nil
	}

	keys := make([]string, len(usernames))
	for i, name := range usernames {
		keys[i] = USERNAME_TO_ID_KEY + name
		logging.Logger.Infof("%s", name)
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
