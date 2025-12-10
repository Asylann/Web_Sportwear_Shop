package db

import (
	"WebSportwareShop/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var db *sqlx.DB

func InitDB(cfg config.Config) {
	var err error
	db, err = sqlx.Open("postgres", cfg.DatabaseConnection)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
		return
	}
	/*RunMigrations(db)*/
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(4 * time.Minute)
}

func CloseDB() {
	db.Close()
}
