package database

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDataBase() (db *sql.DB) {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err.Error())
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
