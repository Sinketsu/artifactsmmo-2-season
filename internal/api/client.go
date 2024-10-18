package api

import (
	"context"

	api "github.com/Sinketsu/artifactsmmo/gen/oas"
)

type Auth struct {
	Token string
}

func (a *Auth) HTTPBasic(_ context.Context, _ string) (api.HTTPBasic, error) {
	return api.HTTPBasic{}, nil
}

func (a *Auth) JWTBearer(_ context.Context, _ string) (api.JWTBearer, error) {
	return api.JWTBearer{Token: a.Token}, nil
}

type Params struct {
	ServerUrl   string
	ServerToken string
}

type Client struct {
	*api.Client
}

func NewClient(params Params) (*Client, error) {
	client, err := api.NewClient(params.ServerUrl, &Auth{Token: params.ServerToken})
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil
}
