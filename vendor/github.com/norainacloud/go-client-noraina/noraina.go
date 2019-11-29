package noraina

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://nacp01.noraina.net/"
	mediaType      = "application/json"
)

type Client struct {
	// HTTP client used to communicate with the Noraina API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// Authentication token
	Token string
}

// NewClient returns a new Noraina API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}

	return c
}

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mediaType)

	if len(c.Token) != 0 {
		req.Header.Add("x-access-token", c.Token)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	res := map[string]string{}
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &res)
		if err != nil {
			return err
		}
	}

	return &ErrorResponse{
		StatusCode: r.StatusCode,
		Message:    res,
	}
}

type ErrorResponse struct {
	StatusCode int
	Message    map[string]string
}

func (r *ErrorResponse) Error() string {
	err := ""
	for key, value := range r.Message {
		err = fmt.Sprintf("%s: %s", key, value)
	}
	return err
}
