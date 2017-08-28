package support

import (
	"context"
	"os"

	"github.com/digitalocean/godo"

	"golang.org/x/oauth2"
)

// Ctx does something.
var Ctx = context.TODO()

// DOClient used for digitalocean api.
var DOClient = GetDOClient()

// TokenSource struct to handle do access token.
type TokenSource struct {
	AccessToken string
}

// Token method of TokenSource returns do access token.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func getToken() string {
	pat := os.Getenv("DOAUTHTOKEN")
	if pat == "" {
		RedLn("No DOAUTHTOKEN env variable! Quiting ...")
		os.Exit(1)
	}
	return pat
}

// GetDOClient return DO api client
func GetDOClient() *godo.Client {
	tokenSource := &TokenSource{
		AccessToken: getToken(),
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	return client

	// 	ctx := context.TODO()

	// 	createRequest := &godo.DropletMultiCreateRequest{
	// 		Names:  []string{"sub-01.example.com", "sub-02.example.com"},
	// 		Region: "nyc3",
	// 		Size:   "512mb",
	// 		Image: godo.DropletCreateImage{
	// 			Slug: "ubuntu-14-04-x64",
	// 		},
	// 		IPv6: true,
	// 		Tags: []string{"web"},
	// 	}

	// 	droplet, _, err := client.Droplets.CreateMultiple(ctx, createRequest)
}
