package config

import (
	"errors"
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func InitGithubConfig() {
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://localhost:8080/auth/github/callback", "read:user", "user:email"),
	)
}

var GoogleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
}

func getConnStr() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DATABASE"),
	)
}
func LoadConfig() (Config, error) {
	connStr := getConnStr()
	fmt.Println(connStr)
	if connStr == "host=  port= user= password= dbname= sslmode=disable" {
		return Config{}, errors.New("invalid database connection string!!!")
	}
	config := Config{DatabaseURL: connStr, Port: os.Getenv("PORT"), JWTSecret: "FHF7g3&*D4"}
	return config, nil
}
