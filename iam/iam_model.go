// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"time"
)

type ValidateRequestBody struct {
	Auth struct {
		Request struct {
			Headers map[string]string `json:"headers"`
			Params  map[string]string `json:"params"`
			Method  string            `json:"method"`
			URI     string            `json:"uri"`
		} `json:"request"`
		Authorization string `json:"authorization"`
		SecurityToken string `json:"security_token,omitempty"`
	} `json:"auth"`
}

type ServiceTokenRequestBody struct {
	Auth struct {
		Scope struct {
			Domain struct {
				ID string `json:"id"`
			} `json:"domain"`
		} `json:"scope"`
		Identity struct {
			Password struct {
				User struct {
					Domain struct {
						Name string `json:"name"`
					} `json:"domain"`
					Password string `json:"password"`
					Name     string `json:"name"`
				} `json:"user"`
			} `json:"password"`
			Methods []string `json:"methods"`
		} `json:"identity"`
	} `json:"auth"`
}

type Token struct {
	TokenID string `json:"-"`
	// id will be return in GetStsToken response
	ID      string `json:"id,omitempty"`
	Catalog []struct {
		Endpoints []struct {
			ID        string `json:"id"`
			Interface string `json:"interface"`
			Region    string `json:"region"`
			URL       string `json:"url"`
		} `json:"endpoints"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"catalog"`
	ExpiresAt string `json:"expires_at"`
	Extras    struct {
	} `json:"extras"`
	IssuedAt string        `json:"issued_at"`
	Methods  []interface{} `json:"methods"`
	Project  struct {
		Domain struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"domain"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Roles []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"roles"`
	User struct {
		Domain struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"domain"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

type IamAuthResponseBody struct {
	IamToken Token `json:"token"`
}

type AccessKey struct {
	Access       string `json:"access"`
	CredentialID string `json:"credential_id"`
	ProjectID    string `json:"project_id"`
	Secret       string `json:"secret"`
	UserID       string `json:"user_id"`
}

type SecretKeyResponseBody struct {
	Key AccessKey `json:"accesskey"`
}

type AssumeRoleResponse struct {
	AccessKey       string    `json:"accessKeyId"`
	SecretAccessKey string    `json:"secretAccessKey"`
	SessionToken    string    `json:"sessionToken"`
	CreateTime      time.Time `json:"createTime"`
	Expiration      time.Time `json:"expiration"`
	UserId          string    `json:"userId"`
	RoleId          string    `json:"roleId"`
	// token will be empty if withToken is false
	OpenStackToken string `json:"-"`
	Token          Token  `json:"token,omitempty"`
}

type RequestArguments struct {
	Authorization string
	Method        string
	URI           string
	SecurityToken string
	SignHeaders   map[string]string
	QueryParams   map[string]string
	RequestID     string // use request id for log tracking
}

type SignArguments struct {
	UserAccessKey       string
	UserSecretKey       string
	Method              string
	Path                string
	Headers             map[string]string
	QueryParams         map[string]string
	TimeStamps          time.Time
	ExpirationInSeconds int
	SignHeaders         []string
}
