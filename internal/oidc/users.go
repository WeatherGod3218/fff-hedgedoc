package oidc

import (
	"context"
	"fff-hedgedoc/internal/logging"
	"fff-hedgedoc/internal/models"
	"fff-hedgedoc/internal/redis"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (client *OIDCClient) GetActiveUsers(ctx context.Context) []models.KeycloakUser {
	return client.GetOIDCGroup(ctx, "a97a191e-5668-43f5-bc0c-6eefc2b958a7")
}

func (client *OIDCClient) GetEBoardMembers(ctx context.Context) []models.KeycloakUser {
	return client.GetOIDCGroup(ctx, "47dd1a94-853c-426d-b181-6d0714074892")
}

func (c *OIDCClient) LoadAllUsers(ctx context.Context) error {
	logging.Logger.Info("Starting Reload of All Users!")
	const pageSize = 250

	for first := 0; ; first += pageSize {
		var page []models.KeycloakUser
		path := fmt.Sprintf("/users?briefRepresentation=true&first=%d&max=%d", first, pageSize)
		if err := c.get(ctx, path, &page); err != nil {
			logging.Logger.WithFields(logrus.Fields{"method": "LoadAllUsers"}).Error(err)
			return nil
		}

		for _, user := range page {
			if err := redis.AddUser(ctx, user.ID, user.Username); err != nil {
				return err
			}
		}

		if len(page) < pageSize {
			break
		}
	}
	logging.Logger.Info("Finished Reload of All Users Successfully!")

	return nil
}

func ReloadUsersWorker(ctx context.Context, client *OIDCClient) {
	ticker := time.NewTicker(time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := client.LoadAllUsers(context.Background()); err != nil {
					logging.Logger.Warnf("FAILED TO RELOAD THE USERS FROM OIDC!: %s", err)
				}
				ticker.Reset(time.Hour)
			case <-client.quit:
				ticker.Stop()
				return
			}
		}
	}()
}
