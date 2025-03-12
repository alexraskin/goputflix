package internal

import (
	"context"

	"github.com/putdotio/go-putio"
	"golang.org/x/oauth2"
)

type PutIoOptions struct {
	Token string
}

func InitPutio(options PutIoOptions) *putio.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: options.Token})

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	client := putio.NewClient(oauthClient)

	return client
}
