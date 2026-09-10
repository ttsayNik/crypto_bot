package binance

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"telebot/internal/domain"
	"telebot/internal/usecase"
)

var coins = []string{
	"BTCUSDT",
	"ETHUSDT",
	"SOLUSDT",
	"BNBUSDT",
	"XRPUSDT",
}

var url = "https://api.binance.com/api/v3/ticker/24hr?symbol=[\"BTCUSDT\", \"ETHUSDT\", \"SOLUSDT\", \"BNBUSDT\", \"XRPUSDT\"]"

type BinanceClient struct {
	client *http.Client
	cu     *usecase.CoinStateUsecase
	logger *slog.Logger
}

func New(client *http.Client, cu *usecase.CoinStateUsecase, logger *slog.Logger) *BinanceClient {
	return &BinanceClient{
		client: client,
		cu:     cu,
		logger: logger,
	}
}

type PreCoinStateDTO struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"lastPrice"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
}

func (a *BinanceClient) GetCoinStates() {
	coinStates := make([]usecase.CoinStateDTO, 0, 5)

	resp, err := a.client.Get(url) // we do batch request
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed batch request: %v", err))
		return
	}
	defer resp.Body.Close()

	preCoinDTOs := make([]PreCoinStateDTO, 0, 5)
	err = json.NewDecoder(resp.Body).Decode(preCoinDTOs)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed decode reponse body: %v", err))
		return
	}

	for i := range preCoinDTOs {
		coinDTO, err := ToDTO(&preCoinDTOs[i])
		if err != nil {
			a.logger.Error(fmt.Sprintf("failed prase predto to dto: %v", err))
			return
		}

		coinStates = append(coinStates, *coinDTO)
	}

	err = a.cu.Update(coinStates)
	if err != nil {
		a.logger.Error(fmt.Sprintf("failed to update cryptocurrency rates: %v", err))
		return
	}
}

func ToDTO(preCoinDTO *PreCoinStateDTO) (*usecase.CoinStateDTO, error) {
	price, err := domain.NewPriceFromString(preCoinDTO.Price)
	if err != nil {
		return nil, fmt.Errorf("error creating price from string %s: %w", &price, err)
	}

	highPrice, err := domain.NewPriceFromString(preCoinDTO.HighPrice)
	if err != nil {
		return nil, fmt.Errorf("error creating price from string %s: %w", &highPrice, err)
	}

	lowPrice, err := domain.NewPriceFromString(preCoinDTO.LowPrice)
	if err != nil {
		return nil, fmt.Errorf("error creating price from string %s: %w", &lowPrice, err)
	}

	coinStateDTO := &usecase.CoinStateDTO{
		Symbol:    preCoinDTO.Symbol,
		Price:     price.Int64(),
		HighPrice: highPrice.Int64(),
		LowPrice:  lowPrice.Int64(),
	}

	return coinStateDTO, nil
}
