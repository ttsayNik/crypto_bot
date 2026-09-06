package handler

import (
	"log/slog"
	"net/http"
	"telebot/internal/usecase"
)

type TGClient struct {
	client *http.Client
	cu     *usecase.CoinStateUsecase
	logger *slog.Logger
}

func NewTGClient(client *http.Client, cu *usecase.CoinStateUsecase, logger *slog.Logger) *TGClient {
	return &TGClient{
		client: client,
		cu:     cu,
		logger: logger,
	}
}
