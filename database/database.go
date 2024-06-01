package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDb() (*sql.DB, error) {
	var err error

	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		return nil, errors.New("Missing db connection string")
	}

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables() {
	var err error
	if !tableExists(db, "user") {
		_, err = db.Exec(`
			CREATE TABLE user (
				id VARCHAR(255) PRIMARY KEY,
				username VARCHAR(255) NOT NULL,
				password VARCHAR(255) NOT NULL
			)`)
		if err != nil {
			panic(err)
		}
	}

	if !tableExists(db, "todo") {
		_, err = db.Exec(`
			CREATE TABLE todo (
				id VARCHAR(255) PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				startDate DATETIME NOT NULL,
				endDate DATETIME NOT NULL,
				userId VARCHAR(255) NOT NULL
			)`)
		if err != nil {
			panic(err)
		}
	}
}

func tableExists(db *sql.DB, tableName string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = ?)", tableName).Scan(&exists)
	if err != nil {
		log.Fatalf("Erro ao verificar se a tabela %s existe: %v", tableName, err)
	}
	return exists
}
