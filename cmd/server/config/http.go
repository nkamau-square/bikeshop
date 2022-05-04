package config

import "net/http"

type Configuration struct {
	squareVersion string
	token         string
	location      string
}

var Config = &Configuration{
	squareVersion: "2022-04-20",
	token:         "EAAAECW5OpsNF5N8XCNC3Pf2fFdJkHlIcjH75XXHSF8r00YbrdECo8SgV5bIaB02",
	location:      "LCDDGWKCAQRZY",
}

func (c *Configuration) NewRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Square-Version": []string{c.squareVersion},
		"Authorization":  []string{"Bearer " + c.token},
		"Content-Type":   []string{"application/json"},
	}
	return req, err
}

func (c *Configuration) GetLocation() string {
	return c.location
}

func (c *Configuration) GetAPIVersion() string {
	return c.squareVersion
}
