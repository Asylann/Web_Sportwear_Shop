package config

import (
	"errors"
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/openidConnect"
	"log"
	"os"
)

func InitOAuthProviders() {
	discoveryURL := "https://accounts.google.com/.well-known/openid-configuration"
	GoogleOpenid, err := openidConnect.NewNamed("google", os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"),
		"https://localhost:8080/auth/google/callback", discoveryURL, "email", "profile", "openid")
	if err != nil {
		log.Printf("Error during init OIDC of google: %v", err.Error())
		return
	}
	GoogleOpenid.SetName("google")
	goth.UseProviders(
		GoogleOpenid,
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"),
			"https://localhost:8080/auth/github/callback", "read:user", "user:email"),
	)
}

type Config struct {
	DatabaseConnection string
	Port               string
	JWT_Secret         string
}

func getConnection() string {
	res := fmt.Sprintf("host=%v port=%v user=%v password=%v database=%v sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGPORT"),
		os.Getenv("PGUSER"), os.Getenv("PASSWORD"),
		os.Getenv("DATABASE"))
	return res
}

func LoadConfig() (Config, error) {
	connectionStr := getConnection()
	if connectionStr == "host= port= user= password= db_name= sslmode=disable" {
		log.Fatal("Connection configuration is null")
		return Config{}, errors.New("Cannot get configuration")
	}
	port := os.Getenv("PORT")
	config := Config{DatabaseConnection: connectionStr, Port: port, JWT_Secret: "Btokhm23f"}
	return config, nil
}
