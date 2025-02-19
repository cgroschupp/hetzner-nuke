package hetzner

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type Project struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type Token struct {
	Token struct {
		Project Project `json:"project"`
	} `json:"token"`
}

type Client struct {
	*hcloud.Client
	Token TokenClient
}

type TokenClient struct {
	client *hcloud.Client
}

func (t *TokenClient) Current(ctx context.Context) (*Token, *hcloud.Response, error) {
	token := &Token{}
	req, err := t.client.NewRequest(ctx, "GET", "/_tokens/current", nil)
	if err != nil {
		return nil, nil, err
	}
	response, err := t.client.Do(req, token)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to make request: %w", err)
	}
	return token, response, nil
}

func NewClient(options ...hcloud.ClientOption) *Client {
	client := hcloud.NewClient(options...)
	return &Client{client, TokenClient{client: client}}
}
