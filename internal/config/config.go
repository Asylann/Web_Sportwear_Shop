package config

import (
	"errors"
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"os"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func InitGithubConfig() {
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://localhost:8080/auth/github/callback", "read:user", "user:email"),
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:8080/auth/google/callback", "email", "profile"),
	)
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
