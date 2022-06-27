package webclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type client struct {
	serverAddress string
	httpClient    *http.Client
}

func New(serverAddress string) *client {
	httpClient := http.Client{}
	newclient := client{
		serverAddress: serverAddress,
		httpClient:    &httpClient,
	}
	return &newclient
}

// Get makes a HTTP request with Get method using the given path
func (c *client) Get(ctx context.Context, path string) (*Response, error) {
	req, err := c.newGetRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) HTTPClient() *http.Client {
	return c.httpClient
}

// buildURL join the server address with the given relative path.
// relative path must start with /
func (c *client) buildURL(path string) string {
	return c.serverAddress + path
}

func (c *client) do(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get resource: %w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}
	result := Response{
		StatusCode: resp.StatusCode,
		Data:       data,
	}
	return &result, nil
}

func (c *client) newGetRequest(ctx context.Context, path string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	return req, nil
}
