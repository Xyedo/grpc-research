package pgrepository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func OpenDb() (*sqlx.DB, error) {
	db := sqlx.MustOpen("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("HOST"),
			os.Getenv("PG_PORT"),
			os.Getenv("PG_USER"),
			os.Getenv("PG_PASSWORD"),
			os.Getenv("PG_DB"),
		),
	)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	duration, err := time.ParseDuration("15m")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, err
}
