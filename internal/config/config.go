package config

import (
	"errors"
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	oidc "github.com/markbates/goth/providers/openidConnect"
	"log"
	"os"
)

func InitOAuthConfig() {
	discoveryURL := "https://accounts.google.com/.well-known/openid-configuration"
	openIdOfGoogle, err := oidc.NewNamed("google", os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:8080/auth/google/callback", discoveryURL,
		"email", "profile", "openid")
	if err != nil {
		log.Fatalf("failed to init OIDC provider: %v", err)
	}
	openIdOfGoogle.SetName("google")
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://localhost:8080/auth/github/callback", "read:user", "user:email"),
		openIdOfGoogle,
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
