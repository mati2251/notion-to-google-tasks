package utils

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

const GOOGLE_AUTH_URI_KEY = "google.auth_uri"
const GOOGLE_TOKEN_URI_KEY = "google.token_uri"
const GOOGLE_AUTH_PROVIDER_X509_CERT_URL_KEY = "google.auth_provider_x509_cert_url"
const GOOGLE_REDIRECT_URIS_KEY = "google.redirect_uris"
const GOOGLE_CLIENT_ID_KEY = "google.client_id"
const GOOGLE_CLIENT_SECRET_KEY = "google.client_secret"

func GoogleConfig() {
	setDefaults()
	getClientIdAndSecret()
}

func setDefaults() {
	viper.SetDefault(GOOGLE_AUTH_URI_KEY, "https://accounts.google.com/o/oauth2/auth")
	viper.SetDefault(GOOGLE_TOKEN_URI_KEY, "https://oauth2.googleapis.com/token")
	viper.SetDefault(GOOGLE_AUTH_PROVIDER_X509_CERT_URL_KEY, "https://www.googleapis.com/oauth2/v1/certs")
	viper.SetDefault(GOOGLE_REDIRECT_URIS_KEY, "http://localhost")
}

func getClientIdAndSecret() {
	prompt := promptui.Prompt{
		Label: "Client ID",
	}

	clientId, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	prompt = promptui.Prompt{
		Label: "Client Secret",
	}

	clientSecret, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	viper.Set(GOOGLE_CLIENT_ID_KEY, clientId)
	viper.Set(GOOGLE_CLIENT_SECRET_KEY, clientSecret)
}
