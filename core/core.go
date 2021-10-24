package core

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type ConsumeMethod string
type Module string

type core struct {
	url  *url.URL
	key  string
	pass string
}

type Core interface {
	ConsumeAPI(method ConsumeMethod, endpoint string, data url.Values) (*http.Response, error)
}

const (
	// MySQL Time format
	TimeFormatMySQLDate     string = "2006-01-02"
	TimeFormatMySQLTime     string = "15:04:05"
	TimeFormatMySQLDateTime string = "2006-01-02 15:04:05"

	// Allowed Method
	MethodGet  ConsumeMethod = http.MethodGet
	MethodPost ConsumeMethod = http.MethodPost

	// Available Module
	ModuleSms      Module = "SMS"
	ModuleEmail    Module = "EMAIL"
	ModuleWhatsapp Module = "WA"
)

func New(hosturl, apikey, password string) (Core, error) {
	u, err := url.Parse(hosturl)
	if err != nil {
		return nil, err
	}

	u.RawQuery = ""
	u.RawFragment = ""
	u.Fragment = ""

	return &core{
		url:  u,
		key:  apikey,
		pass: password,
	}, nil
}

func (c core) baseUrl() string {
	return strings.TrimRight(c.url.String(), "/?#")
}

func (c core) ConsumeAPI(method ConsumeMethod, endpoint string, data url.Values) (*http.Response, error) {
	urlPath := c.baseUrl() + "/" + endpoint
	data.Add("api_key", c.key)
	data.Add("api_pass", c.pass)

	switch method {
	case http.MethodGet:
		return http.Get(urlPath + "?" + data.Encode())
	case http.MethodPost:
		return http.PostForm(urlPath, data)
	}
	return nil, errors.New("method not allowed")
}
