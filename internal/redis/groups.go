package redis

import (
	"context"
	"fff-hedgedoc/internal/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

const GROUPS_KEY_PREFIX = "csh:groups"
const GROUP_LIST_KEY_PREFIX = GROUPS_KEY_PREFIX + ":list"

const GROUPNAME_TO_ID_KEY = GROUPS_KEY_PREFIX + ":grouptoid:"

const MAX_GROUPS_PER_REQ = 50

func AddGroup(ctx context.Context, uuid string, groupName string) error {
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		if err := pipe.Set(ctx, GROUPNAME_TO_ID_KEY+groupName, uuid, -1).Err(); err != nil {
			return err
		}

		if err := pipe.ZAdd(ctx, GROUP_LIST_KEY_PREFIX, redis.Z{Score: 0, Member: groupName}).Err(); err != nil {
			return err
		}

		return nil
	})
	return err
}

func GetGroupsSection(ctx context.Context, prefix string, cursor *string) ([]models.GroupInfo, error) {
	lo := "[" + prefix
	if cursor != nil && strings.TrimSpace(*cursor) != "" {
		if *cursor >= prefix {
			lo = "(" + *cursor
		}
	}
	hi := "[" + prefix + "\xff"

	groupnames, err := client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   GROUP_LIST_KEY_PREFIX,
		Start: lo,
		Stop:  hi,
		Count: MAX_GROUPS_PER_REQ,
		ByLex: true,
	}).Result()
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(groupnames))
	for i, name := range groupnames {
		keys[i] = GROUPNAME_TO_ID_KEY + name
	}

	uuids, err := client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	groups := make([]models.GroupInfo, len(groupnames))
	for index, group := range groupnames {
		if uuids[index] != nil {
			groups[index] = models.GroupInfo{
				UUID:      uuids[index].(string),
				GroupName: group,
			}
		}

	}

	return groups, nil
}
