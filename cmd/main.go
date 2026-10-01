package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"telebot/internal/adapters/api/binance"
	repo "telebot/internal/adapters/db/sqlite"
	handler "telebot/internal/adapters/delivery/http"
	"telebot/internal/usecase"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: add telegram token and add telegram bot api comands with answers

// const token = "8940190856:AAGkgAzOjqpL8LACSLO6CDQqs5bRT_bXsI4"

func main() {
	// init gracefull off
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// init db
	db, err := sql.Open("sqlite3", "telebot.db")
	if err != nil {
		log.Fatalf("error opening: %v", err)
	}

	// query := `
	// CREATE TABLE IF NOT EXISTS coin_state(
	// 	name VARCHAR(20) NOT NULL,
	// 	price BIGINT NOT NULL,
	// 	high_price BIGINT NOT NULL,
	// 	low_price BIGINT NOT NULL,
	//  percent_change INT NOT NULL
	// 	last_update TIMESTAMP NOT NULL
	// 	);
	// 	`

	// _, err = db.Exec(query)
	// if err != nil {
	// 	log.Println("error execute the table")
	// }

	// init logger
	logFile, err := os.OpenFile("loggs.txt", os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		log.Fatalf("error opening loggs file: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))

	// init repositories
	coinRepo, err := repo.NewCoinStateRepo(db)
	if err != nil {
		log.Fatalf("error creating coin state repo: %v", err)
	}

	// init usecases
	coinUsecase, err := usecase.NewCoinStateUsecase(coinRepo, logger)
	if err != nil {
		log.Fatalf("error creating coin state usecase: %v", err)
	}

	// init clients
	binanceClient, err := binance.New(&http.Client{}, coinUsecase, logger)
	if err != nil {
		log.Fatalf("error creating new binance: %v", err)
	}
	handler, err := handler.New(&http.Client{}, coinUsecase, logger)
	if err != nil {
		log.Fatalf("error creating new handler: %v", err)
	}

	// register handlers
	http.HandleFunc("GET /telebot", handler.GetCoinStates)

	// init producer
	ticker := time.NewTicker(5 * time.Minute)
	go func(ctx context.Context) {
		for range ticker.C {
			select {
			case <-ctx.Done():
				break
			default:
				binanceClient.GetCoinStates()
			}
		}
	}(ctx)

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
