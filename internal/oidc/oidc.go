package oidc

import (
	"bytes"
	"context"
	"encoding/json"
	"fff-hedgedoc/internal/logging"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/sirupsen/logrus"
)

var httpClient = &http.Client{Timeout: 45 * time.Second}

type OIDCClient struct {
	oidcClientId     string
	oidcClientSecret string

	accessToken  string
	providerBase string
	quit         chan struct{}
}

type OIDCUser struct {
	Uuid     string `json:"id"`
	Username string `json:"username"`
	Gatekeep bool   `json:"result"`
	SlackUID string `json:"slackuid"`
}

var groupCache map[string]string

func (client *OIDCClient) SetupOIDCClient(oidcClientId, oidcClientSecret string) {
	client.oidcClientId = oidcClientId
	client.oidcClientSecret = oidcClientSecret
	parse, err := url.Parse(csh_auth.ProviderURI)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"method": "setupOidcClient"}).Error(err)
		return
	}
	groupCache = make(map[string]string)
	client.providerBase = parse.Scheme + "://" + parse.Host
	exp := client.GetAccessToken()
	ticker := time.NewTicker(time.Duration(exp) * time.Second)
	// this will async get the token
	go func() {
		for {
			select {
			case <-ticker.C:
				exp = client.GetAccessToken()
				ticker.Reset(time.Duration(exp) * time.Second)
			case <-client.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (client *OIDCClient) GetAccessToken() int {
	htclient := http.DefaultClient
	//request body
	authData := url.Values{}
	authData.Set("client_id", client.oidcClientId)
	authData.Set("client_secret", client.oidcClientSecret)
	authData.Set("grant_type", "client_credentials")
	resp, err := htclient.PostForm(csh_auth.ProviderURI+"/protocol/openid-connect/token", authData)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"method": "getAccessToken"}).Error(err)
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logging.Logger.WithFields(logrus.Fields{"method": "getAccessToken", "statusCode": resp.StatusCode}).Error(resp.Status)
		return 0
	}
	respData := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"method": "getAccessToken"}).Error(err)
		return 0
	}
	if respData["error"] != nil {
		logging.Logger.WithFields(logrus.Fields{"method": "getAccessToken"}).Error(respData)
		return 0
	}
	client.accessToken = respData["access_token"].(string)
	return int(respData["expires_in"].(float64))
}

func (c *OIDCClient) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.providerBase+"/auth/admin/realms/csh"+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("status %d: %s", resp.StatusCode, bytes.TrimSpace(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
