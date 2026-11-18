package oidc

import (
	"context"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/models"
	"fff-hedgedoc/internal/redis"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func (client *OIDCClient) collect(ctx context.Context, groups []models.KeycloakGroup) error {
	for _, g := range groups {
		redis.SetKeyValueStore(context.Background(), redis.GROUPS_KEY_PREFIX+g.ID, g.Name, -1)

		subs := g.SubGroups
		if len(subs) == 0 && g.SubGroupCount > 0 {
			if err := client.get(ctx, "/groups/"+g.ID+"/children?max=1000", &subs); err != nil {
				return err
			}
		}

		if err := client.collect(ctx, subs); err != nil {
			return err
		}

		for _, group := range subs {
			if err := redis.AddGroup(ctx, group.ID, group.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

func (client *OIDCClient) FindOIDCGroupID(ctx context.Context, name string) string {
	if groupCache[name] != "" {
		return groupCache[name]
	}

	ret := make([]map[string]any, 0)
	if err := client.get(ctx, "/groups?exact=true&search="+name, &ret); err != nil {
		logging.Logger.WithFields(logrus.Fields{"method": "FindOIDCGroupID"}).Error(err)
		return ""
	}

	//Example:
	//[{"id":"47dd1a94-853c-426d-b181-6d0714074892","name":"eboard","path":"/eboard","subGroups":[{"id":"66b9578a-2b58-46a6-8040-59388e57e830","name":"eboard-opcomm","path":"/eboard/eboard-opcomm","subGroups":[]}]}]
	//it returns as an array for SOME reason, so we cut to the group we want
	group := ret[0]
	// and now we have SUBgroups, so we do this fucked parse
	subGroups := group["subGroups"].([]any)
	fmt.Println(subGroups)
	subGroup := subGroups[0].(map[string]interface{})
	fmt.Println(subGroup)
	gid := subGroup["id"].(string)
	groupCache[name] = gid
	return gid
}

func (client *OIDCClient) GetOIDCGroup(ctx context.Context, groupID string) []models.KeycloakUser {
	ret := make([]models.KeycloakUser, 0)
	if err := client.get(ctx, "/groups/"+groupID+"/members", &ret); err != nil {
		logging.Logger.WithFields(logrus.Fields{"method": "GetOIDCGroup"}).Error(err)
		return nil
	}
	return ret
}

func (client *OIDCClient) LoadMembersFromGroup(ctx context.Context, groupID string, groupName string) ([]models.KeycloakUser, error) {
	const pageSize = 100
	var members []models.KeycloakUser

	for first := 0; ; first += pageSize {
		q := url.Values{}
		q.Set("first", strconv.Itoa(first)) //ts exists??
		q.Set("max", strconv.Itoa(pageSize))
		q.Set("briefRepresentation", "true")

		uri := fmt.Sprintf("/groups/%s/members?%s", url.PathEscape(groupID), q.Encode())

		var page []models.KeycloakUser
		if err := client.get(ctx, uri, &page); err != nil {
			return nil, err
		}
		members = append(members, page...)

		if len(page) < pageSize {
			break
		}
	}

	return members, nil
}

func (client *OIDCClient) LoadAllGroups(ctx context.Context) error {
	logging.Logger.Warn("Starting Reload of All Groups!")

	var top []models.KeycloakGroup

	if err := client.get(ctx, "/groups?briefRepresentation=true&max=1000", &top); err != nil {
		return err
	}

	if err := client.collect(ctx, top); err != nil {
		logging.Logger.Warnf("Failed to load User groups! %s", err)
		return err
	}

	for _, group := range top {
		members, err := client.LoadMembersFromGroup(ctx, group.ID, group.Name)
		if err != nil {
			logging.Logger.Warnf("Failure loading members from group %s, with error %s", group.ID, err)
			return err
		}

		if err := redis.AddMembersOfGroup(ctx, group.ID, members); err != nil {
			return err
		}

		if err := redis.AddGroup(ctx, group.ID, group.Name); err != nil {
			return err
		}
	}

	logging.Logger.Warn("Finished Reload of All Groups Successfully!")

	return nil
}

func ReloadGroupsWorker(ctx context.Context, client *OIDCClient) {
	ticker := time.NewTicker(3 * time.Hour) //add 30 minutes to offset from the user

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := client.LoadAllGroups(context.Background()); err != nil {
					logging.Logger.Warnf("FAILED TO RELOAD THE GROUPS FROM OIDC!: %s", err)
				}
				ticker.Reset(6 * time.Hour)
			case <-client.quit:
				ticker.Stop()
				return
			}
		}
	}()
}
