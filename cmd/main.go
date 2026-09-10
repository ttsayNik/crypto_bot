package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"telebot/internal/adapters/api/binance"
	repo "telebot/internal/adapters/db/sqlite"
	handler "telebot/internal/adapters/delivery/http"
	"telebot/internal/usecase"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// init db
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

	// init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))

	// init repositories
	coinRepo := repo.NewCoinStateRepo(db)

	// init usecases
	coinUsecase := usecase.NewCoinStateUsecase(coinRepo, logger)

	// init clients
	apiClient := binance.New(&http.Client{}, coinUsecase, logger)
	handler := handler.NewTGClient(&http.Client{}, coinUsecase, logger)

	// register handlers
	http.HandleFunc("GET /telebot", handler.GetCoinStates)

	// init producer
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			apiClient.GetCoinStates()
		}
	}()

	fmt.Println("server starts at localhost:8080")

	http.ListenAndServe(":8080", nil)
}

// func DropTable(db *sql.DB, table string) error {
// 	query := "drop table " + table

// 	_, err := db.Exec(query)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
