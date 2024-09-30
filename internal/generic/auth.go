package generic

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
