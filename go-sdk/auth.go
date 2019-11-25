package go_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (c *NorainaApiClient) GetAuthToken(email string, password string) (string, error) {
	authRequest := AuthRequest{
		Email:    email,
		Password: password,
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(authRequest)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(norainaDomain+loginRoute, "application/json", b)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	authResponse := &AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(authResponse)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("[ERROR] Auth Error with Status code %v, message %v", resp.StatusCode, authResponse.Data.Message)
	}

	if authResponse.Status != "success" {
		return "", fmt.Errorf("[ERROR] Auth returns HTTP Status 200 but no token could be retrieved, status is %v, messsage %v", authResponse.Status, authResponse.Data.Message)
	}

	return authResponse.Data.Token, nil
}
