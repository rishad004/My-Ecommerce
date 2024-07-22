package helper

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Google() *oauth2.Config {

	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  "https://http://adf92f645b67e4a27bca9dfc175ed059-1876988238.eu-north-1.elb.amazonaws.com//user/google/callback",
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
