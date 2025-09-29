package tests

import (
	"WebSportwareShop/internal/cache"
	"WebSportwareShop/internal/config"
	"WebSportwareShop/internal/db"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMain(M *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("No .env variables are loaded")
		return
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Doesnt open config")
	}

	// Initialization of Db connection
	db.InitDB(cfg)

	if err := cache.InitRedisConnection(); err != nil {
		log.Fatal("Doesnt open redis")
	}

	code := M.Run()

	db.CloseDB()

	os.Exit(code)
}
