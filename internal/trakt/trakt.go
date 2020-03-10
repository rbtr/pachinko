/*
Copyright © 2020 The Pachinko Authors

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package trakt

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/rbtr/go-trakt"
	gotrakt "github.com/rbtr/go-trakt"
)

const (
	DefaultAuthfile = "/etc/pachinko/trakt"
	ClientID        = "76a0c1e8d3331021f6e312115e27fe4c29f4ef23ef89a0a69143a62d136ab994"
	// nolint: gosec
	ClientSecret = "fe8d1f0921413028f92428d2922e13a728e27d2f35b26e315cf3dde31228568d"
)

type Auth struct {
	AccessToken  string        `json:"access-token,omitempty"`
	ClientID     string        `json:"client-id,omitempty"`
	ClientSecret string        `json:"client-secret,omitempty"`
	CreatedAt    time.Time     `json:"created-at,omitempty"`
	ExpiresAfter time.Duration `json:"expires-after,omitempty"`
	RefreshToken string        `json:"refresh-token,omitempty"`
}

func (auth *Auth) IsExpired() bool {
	return time.Now().After(auth.CreatedAt.Add(auth.ExpiresAfter))
}

func (auth *Auth) ShouldRefresh(threshold time.Duration) bool {
	return time.Now().After(auth.CreatedAt.Add(auth.ExpiresAfter).Add(-threshold))
}

func ReadAuthFile(path string) (*Auth, error) {
	auth := &Auth{}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return auth, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, auth); err != nil {
		return nil, err
	}
	return auth, nil
}

func WriteAuthFile(path string, auth *Auth) error {
	b, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0600)
}

type Trakt struct {
	*gotrakt.Client
	auth *Auth
}

func NewTrakt(auth *Auth) (*Trakt, error) {
	if auth.ClientID == "" {
		auth.ClientID = ClientID
	}
	if auth.ClientSecret == "" {
		auth.ClientSecret = ClientSecret
	}
	client, err := gotrakt.NewClient(nil, auth.ClientID, auth.ClientSecret)
	if auth.AccessToken != "" {
		client.SetAuthorization(auth.AccessToken)
	}
	return &Trakt{client, auth}, err
}

// Authorize authorizes the client using 2-legged oauth.
// Authorized credentials are stored in the client and also returned.
func (t *Trakt) Authorize(ctx context.Context) (*Auth, error) {
	res, err := t.DeviceCode(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf(
		"Authenticating in Trakt!\nPlease open in your browser:\t%s\n\t and enter the code:\t\t%s\n",
		res.VerificationURL,
		res.UserCode,
	)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(res.ExpiresIn)*time.Second)
	defer cancel()

	ticker := time.NewTicker(time.Duration(res.Interval) * time.Second)
	go func(ctx context.Context, cancelFunc func()) {
		<-ctx.Done()
		cancelFunc()
	}(ctx, ticker.Stop)

	var result *trakt.AuthResult
	code := res.DeviceCode
	for range ticker.C {
		result, err = t.DeviceToken(ctx, code)
		if err == nil {
			break
		}
	}
	if result == nil {
		return nil, err
	}
	fmt.Printf("Success! Your Authorization token is:\n\t> %s <\n", result.AccessToken)
	t.auth.AccessToken = result.AccessToken
	t.auth.RefreshToken = result.RefreshToken
	t.auth.ExpiresAfter = time.Duration(result.ExpiresIn) * time.Second
	t.auth.CreatedAt = time.Unix(int64(result.CreatedAt), 0)
	return t.auth, nil
}

// Refresh reauthorizes the client using the refresh token.
// Authorized credentials are stored in the client and also returned.
func (t *Trakt) Refresh(ctx context.Context) (*Auth, error) {
	res, err := t.RefreshToken(ctx, t.auth.RefreshToken)
	if err != nil {
		return nil, err
	}
	t.SetAuthorization(res.AccessToken)
	t.auth.AccessToken = res.AccessToken
	t.auth.RefreshToken = res.RefreshToken
	t.auth.ExpiresAfter = time.Duration(res.ExpiresIn) * time.Second
	t.auth.CreatedAt = time.Unix(int64(res.CreatedAt), 0)
	return t.auth, nil
}
