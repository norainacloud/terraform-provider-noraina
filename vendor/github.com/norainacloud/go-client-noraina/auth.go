package noraina

import (
	"context"
	"net/http"
)

const loginRoute = "api/login"

type AuthRequest struct {
	Email    string `json:"mail"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Status string `json:"status"`
	Data   AuthResponseData
}

type AuthResponseData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func (c *Client) GetAuthToken(ctx context.Context, authData *AuthRequest) error {
	req, err := c.NewRequest(ctx, http.MethodPost, loginRoute, authData)
	if err != nil {
		return err
	}

	data := new(AuthResponseData)

	err = c.Do(ctx, req, data)
	if err != nil {
		return err
	}

	c.Token = data.Token
	return nil
}
