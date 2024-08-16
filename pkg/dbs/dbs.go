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

var (
	dbname   = os.Getenv("DB_NAME")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New(ctx context.Context) *Service {
	var err error
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname),
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
	log.Printf("Disconnected from database: %s", dbname)
	return s.DB.Close()
}
