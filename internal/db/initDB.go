package db

import (
	"WebSportwareShop/internal/config"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var db *sql.DB

func InitDB(cfg config.Config) {
	var err error
	db, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
		return
	}
	RunMigrations(db)
	initStmt(db)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(4 * time.Minute)
	log.Println("DB is connectedQ!!!")
}

func CloseDB() {
	db.Close()
}
