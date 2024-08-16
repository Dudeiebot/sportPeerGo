package dbs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Service struct {
	DB *sql.DB
}

type DBConfig struct {
	DBName     string
	DBPassword string
	DBUsername string
	DBPort     string
	DBHost     string
}

var dbConfig = &DBConfig{
	DBName:     os.Getenv("DB_NAME"),
	DBPassword: os.Getenv("DB_PASSWORD"),
	DBUsername: os.Getenv("DB_USERNAME"),
	DBPort:     os.Getenv("DB_PORT"),
	DBHost:     os.Getenv("DB_HOST"),
}

func New(ctx context.Context) *Service {
	var err error
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			dbConfig.DBName,
			dbConfig.DBPassword,
			dbConfig.DBHost,
			dbConfig.DBPort,
			dbConfig.DBName,
		),
	)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Db Connected")

	return &Service{DB: db}
}

func (s *Service) Close() error {
	log.Printf("Disconnected from database: %s", dbConfig.DBName)
	return s.DB.Close()
}
