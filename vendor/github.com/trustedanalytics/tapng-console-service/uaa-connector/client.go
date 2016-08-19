/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package uaa_connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Jti          string `json:"jti"`
}

type TapJWTToken struct {
	Jti       string   `json:"jti"`
	Sub       string   `json:"sub"`
	Scope     []string `json:"scope"`
	ClientId  string   `json:"client_id"`
	Cid       string   `json:"cid"`
	Azp       string   `json:"azp"`
	GrantType string   `json:"grant_type"`
	UserId    string   `json:"user_id"`
	Username  string   `json:"user_name"`
	Email     string   `json:"email"`
	RevSig    string   `json:"rev_sig"`
	Iat       int64    `json:"iat"`
	Exp       int64    `json:"exp"`
	Iss       string   `json:"iss"`
	Zid       string   `json:"zid"`
	Aud       []string `json:"aud"`
}

type UaaApi interface {
	Login(username, password string) (*LoginResponse, error)
	ValidateOauth2Token(token string) (*TapJWTToken, error)
}

type UaaConnector struct {
	ClientId     string
	ClientSecret string
	Client       *http.Client
}

func NewUaaBasicAuth(clientId, clientSecret string) (*UaaConnector, error) {
	client, _, err := brokerHttp.GetHttpClient()
	if err != nil {
		return nil, err
	}
	return &UaaConnector{clientId, clientSecret, client}, nil
}

func (u *UaaConnector) Login(username, password string) (*LoginResponse, error) {
	loginResp := LoginResponse{}

	url := os.Getenv("SSO_TOKEN_URI")
	reqBody := fmt.Sprintf("grant_type=password&response_type=token&client_id=%s&client_secret=%s&username=%s&password=%s",
		u.ClientId, u.ClientSecret, username, password)

	auth := brokerHttp.BasicAuth{u.ClientId, u.ClientSecret}
	status, resp, err := brokerHttp.RestUrlEncodedPOST(url, reqBody, brokerHttp.GetBasicAuthHeader(&auth), u.Client)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.New("Bad response status: " + strconv.Itoa(status))
	}

	err = json.Unmarshal(resp, &loginResp)
	if err != nil {
		return nil, err
	}

	return &loginResp, nil
}

func (u *UaaConnector) ValidateOauth2Token(token string) (*TapJWTToken, error) {
	jwtToken := TapJWTToken{}

	url := os.Getenv("SSO_CHECK_TOKEN_URI")
	reqBody := fmt.Sprintf("token=%s", token)

	auth := brokerHttp.BasicAuth{u.ClientId, u.ClientSecret}
	status, resp, err := brokerHttp.RestUrlEncodedPOST(url, reqBody, brokerHttp.GetBasicAuthHeader(&auth), u.Client)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.New("Bad response status: " + strconv.Itoa(status))
	}

	err = json.Unmarshal(resp, &jwtToken)
	if err != nil {
		return nil, err
	}

	return &jwtToken, nil
}
