package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"telebot/internal/usecase"
	"time"
)

type Handler struct {
	client *http.Client
	cu     *usecase.CoinStateUsecase
	logger *slog.Logger
}

func NewTGClient(client *http.Client, cu *usecase.CoinStateUsecase, logger *slog.Logger) *Handler {
	return &Handler{
		client: client,
		cu:     cu,
		logger: logger,
	}
}

type coinStateResponse struct {
	Symbol         string    `json:"symbol"`
	Price          string    `json:"price"`
	HighPrice      string    `json:"high"`
	LowPrice       string    `json:"low"`
	LastUpdate     time.Time `json:"last_update"`
	ProcentChanges string    `json:"procent_changes"`
}

func (h *Handler) GetCoinStates(w http.ResponseWriter, r *http.Request) {
	resp, err := h.cu.Get()
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed try to get the coin states", err))
		return
	}

	coinStates := make([]coinStateResponse, len(resp))
	for i := range resp {
		symbol := resp[i].GetSymbol()
		price := resp[i].GetPrice()
		high := resp[i].GetHighPrice()
		low := resp[i].GetLowPrice()

		coinStates[len(resp)-(i+1)] = coinStateResponse{
			Symbol:     symbol.String(),
			Price:      price.String(),
			HighPrice:  high.String(),
			LowPrice:   low.String(),
			LastUpdate: resp[i].GetLastUpdate(),
		}
	}

	if err := json.NewEncoder(w).Encode(coinStates); err != nil {
		h.logger.Error(fmt.Sprintf("failed try to encode coin states", err))
		return
	}
}
