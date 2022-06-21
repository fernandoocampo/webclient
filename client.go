package webclient

import (
	"net/http"
)

type client struct {
	httpClient *http.Client
}

func New() *client {
	httpClient := http.Client{}
	newclient := client{
		httpClient: &httpClient,
	}
	return &newclient
}

func (c *client) HTTPClient() *http.Client {
	return c.httpClient
}
