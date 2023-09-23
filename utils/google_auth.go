package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/tasks/v1"
)

const GOOGLE_AUTH_URI_KEY = "google.auth_uri"
const GOOGLE_TOKEN_URI_KEY = "google.token_uri"
const GOOGLE_AUTH_PROVIDER_X509_CERT_URL_KEY = "google.auth_provider_x509_cert_url"
const GOOGLE_REDIRECT_URIS_KEY = "google.redirect_uris"
const GOOGLE_CLIENT_ID_KEY = "google.client_id"
const GOOGLE_CLIENT_SECRET_KEY = "google.client_secret"
const GOOGLE_REFRESH_TOKEN_KEY = "google.refresh_token"
const GOOGLE_ACCESS_TOKEN_KEY = "google.access_token"
const GOOGLE_TOKEN_TYPE_KEY = "google.token_type"
const GOOGLE_EXPIRY_KEY = "google.expiry"

func GoogleConfig() (*http.Client, error) {
	setDefaults()
	getClientIdAndSecret()
	return setGoogleToken()
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

func getGoogleOauth2Conf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString(GOOGLE_CLIENT_ID_KEY),
		ClientSecret: viper.GetString(GOOGLE_CLIENT_SECRET_KEY),
		RedirectURL:  viper.GetString(GOOGLE_REDIRECT_URIS_KEY),
		Scopes: []string{
			tasks.TasksReadonlyScope,
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  viper.GetString(GOOGLE_AUTH_URI_KEY),
			TokenURL: viper.GetString(GOOGLE_TOKEN_URI_KEY),
		},
	}

}

func setGoogleToken() (*http.Client, error) {
	conf := getGoogleOauth2Conf()
	tok, err := getTokenFromConfig()
	if err != nil {
		tok, err = getTokenFromWeb(conf)
		if err != nil {
			log.Fatalf("Unable to get token from web: %v", err)
			return nil, err
		}
	}

	viper.Set(GOOGLE_REFRESH_TOKEN_KEY, tok.RefreshToken)
	viper.Set(GOOGLE_ACCESS_TOKEN_KEY, tok.AccessToken)
	viper.Set(GOOGLE_TOKEN_TYPE_KEY, tok.TokenType)
	viper.Set(GOOGLE_EXPIRY_KEY, tok.Expiry)
	return conf.Client(context.Background(), tok), nil
}

func GetGoogleToken() (*http.Client, error) {
	conf := getGoogleOauth2Conf()
	tok, err := getTokenFromConfig()
	if err != nil {
		log.Fatalf("Unable to get token from web: %v", err)
		return nil, err
	}

	return conf.Client(context.Background(), tok), nil
}

func getTokenFromConfig() (*oauth2.Token, error) {
	if viper.GetString(GOOGLE_REFRESH_TOKEN_KEY) == "" {
		return nil, errors.New("no refresh token found")
	}
	tok := &oauth2.Token{
		RefreshToken: viper.GetString(GOOGLE_REFRESH_TOKEN_KEY),
		AccessToken:  viper.GetString(GOOGLE_ACCESS_TOKEN_KEY),
		TokenType:    viper.GetString(GOOGLE_TOKEN_TYPE_KEY),
		Expiry:       viper.GetTime(GOOGLE_EXPIRY_KEY),
	}
	return tok, nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	prompt := promptui.Prompt{
		Label: "Authorization code",
	}
	authCode, err := prompt.Run()
	if err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
		return nil, err
	}
	return tok, nil
}
