package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

func NewConf() *oauth2.Config {
	res := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost",
		Scopes:       []string{"https://picasaweb.google.com/data/"},
		Endpoint:     google.Endpoint,
	}

	return res
}

func GMailNewConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost",
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https: //mail.google.com/",
		},
	}
}
