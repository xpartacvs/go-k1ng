package core

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// A type to group allowed HTTP method
type ConsumeMethod string

// A type to group main modules
type Module string

// The main struct of the K1NG's core service
type Core struct {
	url  *url.URL
	key  string
	pass string
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

// Create a new core service from the given arguments
func New(hosturl, apikey, password string) (*Core, error) {
	u, err := url.Parse(hosturl)
	if err != nil {
		return nil, err
	}

	u.RawQuery = ""
	u.RawFragment = ""
	u.Fragment = ""

	return &Core{
		url:  u,
		key:  apikey,
		pass: password,
	}, nil
}

func (c Core) baseUrl() string {
	return strings.TrimRight(c.url.String(), "/?#")
}

// Consume API according to the given arguments
func (c Core) ConsumeAPI(method ConsumeMethod, endpoint string, data url.Values) (*http.Response, error) {
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
