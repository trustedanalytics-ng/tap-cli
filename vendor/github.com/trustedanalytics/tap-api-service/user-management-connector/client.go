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

package user_management_connector

import (
	"net/http"

	"fmt"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
	"strings"
)

const (
	DEFAULT_ORG string = "00000000-0000-0000-0000-000000000000"
)

type InvitationRequest struct {
	Email string `json:"email"`
}

type InvitationResponse struct {
	State   string `json:"state"`
	Details string `json:"details"`
}

type User struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UaaUser struct {
	Guid     string `json:"guid"`
	Username string `json:"username"`
}

type UserManagementApi interface {
	InviteUser(email string) (*InvitationResponse, error)
	DeleteUser(email string) (int, error)
	GetInvitations() ([]string, error)
	GetUsers() ([]UaaUser, error)
}

type UserManagementApiConnectorFactory struct {
	Address string
	Client  *http.Client
}

type UserManagementApiConnector struct {
	Address       string
	Authorization string
	Client        *http.Client
}

func (c *UserManagementApiConnectorFactory) GetConfiguredUserManagementConnector(authorization string) *UserManagementApiConnector {
	return &UserManagementApiConnector{c.Address, authorization, c.Client}
}

func NewUserManagementConnectorFactory(address string) (*UserManagementApiConnectorFactory, error) {
	client, _, err := brokerHttp.GetHttpClient()
	if err != nil {
		return nil, err
	}
	return &UserManagementApiConnectorFactory{address, client}, nil
}

func NewUserManagementConnectorFactoryWithSSL(address, certPemFile, keyPemFile, caPemFile string) (*UserManagementApiConnectorFactory, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &UserManagementApiConnectorFactory{address, client}, nil
}

func (c *UserManagementApiConnector) getApiConnector(url string) brokerHttp.ApiConnector {
	auth := strings.Split(c.Authorization, " ")
	if len(auth) != 2 {
		auth = []string{"", ""}
	}
	return brokerHttp.ApiConnector{
		OAuth2: &brokerHttp.OAuth2{auth[0], auth[1]},
		Client: c.Client,
		Url:    url,
	}
}

func (u *UserManagementApiConnector) CurrentUser() (*User, error) {
	user := &User{}

	connector := u.getApiConnector(fmt.Sprintf("%s/rest/users/current", u.Address))

	_, err := brokerHttp.GetModel(connector, http.StatusOK, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserManagementApiConnector) InviteUser(email string) (*InvitationResponse, error) {
	invResp := &InvitationResponse{}

	url := fmt.Sprintf("%s/rest/invitations", u.Address)
	connector := u.getApiConnector(url)

	reqBody := InvitationRequest{
		Email: email,
	}
	_, err := brokerHttp.PostModel(connector, reqBody, http.StatusOK, invResp)
	if err != nil {
		return nil, err
	}

	if invResp.State == "NEW" {
		return invResp, nil
	}

	url += "/" + email + "/resend"
	connector = u.getApiConnector(url)

	_, err = brokerHttp.PostModel(connector, reqBody, http.StatusOK, "")
	if err != nil {
		return nil, err
	}

	return invResp, nil
}

func (u *UserManagementApiConnector) GetInvitations() ([]string, error) {
	connector := u.getApiConnector(fmt.Sprintf("%s/rest/invitations", u.Address))

	invitations := []string{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, &invitations)
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

func (u *UserManagementApiConnector) invitationExists(email string) bool {
	invitations, err := u.GetInvitations()
	if err != nil {
		return false
	}

	for _, inv := range invitations {
		if inv == email {
			return true
		}
	}

	return false
}

func (u *UserManagementApiConnector) GetUsers() ([]UaaUser, error) {
	connector := u.getApiConnector(fmt.Sprintf("%s/rest/orgs/%s/users", u.Address, DEFAULT_ORG))

	users := []UaaUser{}
	_, err := brokerHttp.GetModel(connector, http.StatusOK, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserManagementApiConnector) getUserUUID(email string) string {
	users, err := u.GetUsers()
	if err != nil {
		return ""
	}

	for _, u := range users {
		if u.Username == email {
			return u.Guid
		}
	}

	return ""
}

func (u *UserManagementApiConnector) DeleteUser(email string) (int, error) {
	connector := u.getApiConnector(fmt.Sprintf("%s/rest/invitations/%s", u.Address, email))

	if u.invitationExists(email) {
		status, err := brokerHttp.DeleteModel(connector, http.StatusOK)
		if err != nil {
			return status, err
		}
		return http.StatusNoContent, nil
	}

	if uid := u.getUserUUID(email); uid != "" {
		connector = u.getApiConnector(fmt.Sprintf("%s/rest/orgs/%s/users/%s", u.Address, DEFAULT_ORG, uid))
		status, err := brokerHttp.DeleteModel(connector, http.StatusOK)
		if err != nil {
			return status, err
		}
		return http.StatusNoContent, nil
	}

	return http.StatusNotFound, nil
}