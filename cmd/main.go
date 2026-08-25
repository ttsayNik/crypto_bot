package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"telebot/internal/adapters/api"
	repo "telebot/internal/adapters/db/sqlite"
	"telebot/internal/usecase"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ticker := time.NewTicker(10 * time.Second)

	db, err := sql.Open("sqlite3", "telebot.db")
	if err != nil {
		log.Println("error opening:", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS coin_state(
		name VARCHAR(20) NOT NULL,
		price BIGINT NOT NULL,
		high_price BIGINT NOT NULL,
		low_price BIGINT NOT NULL,
		last_update TIMESTAMP NOT NULL
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Println("error execute the table")
	}

	coinRepo := repo.NewCoinStateRepo(db)
	coinUsecase := usecase.NewCoinStateUsecase(coinRepo)
	client := api.NewApiClient(&http.Client{}, coinUsecase)

	go func() {
		for range ticker.C {
			client.GetCoinState()
		}
	}()

	fmt.Println("thats okey")

	http.ListenAndServe(":8080", nil)
}

func DropTable(db *sql.DB, table string) error {
	query := "drop table " + table

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
