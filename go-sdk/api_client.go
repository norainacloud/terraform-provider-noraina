package go_sdk

import (
	"net/http"
)

const (
	norainaDomain    = "https://nacp01.noraina.net/"
	loginRoute       = "api/login"
	instanceRoute    = "api/instance"
	certificateRoute = "api/certificate"
)

type NorainaApiClient struct {
	Client *http.Client
	Token  string
}

func NewNorainaApiClient(email string, password string) (*NorainaApiClient, error) {
	c := &NorainaApiClient{
		Client: &http.Client{},
	}

	token, err := c.GetAuthToken(email, password)
	if err != nil {
		return nil, err
	}
	c.Token = token

	return c, nil
}
