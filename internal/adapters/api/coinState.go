package api

import (
	"encoding/json"
	"log"
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

var url = "https://api.binance.com/api/v3/ticker/24hr?symbol="

type ApiClient struct {
	client *http.Client
	cu     *usecase.CoinStateUsecase
}

func NewApiClient(client *http.Client, cu *usecase.CoinStateUsecase) *ApiClient {
	return &ApiClient{
		client: client,
		cu:     cu,
	}
}

type PreCoinStateDTO struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"lastPrice"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
}

func (a *ApiClient) GetCoinState() {
	coinStates := make([]usecase.CoinStateDTO, 0, 5)

	for i := range coins {
		resp, err := a.client.Get(url + coins[i])
		if err != nil {
			log.Println(err)
			return
		}

		preCoinDTO := &PreCoinStateDTO{}
		err = json.NewDecoder(resp.Body).Decode(preCoinDTO)
		if err != nil {
			log.Println(err)
			return
		}

		coinDTO, err := ToDTO(preCoinDTO)
		if err != nil {
			log.Println(err)
			return
		}

		coinStates = append(coinStates, *coinDTO)
	}

	err := a.cu.Update(coinStates)
	if err != nil {
		log.Println(err)
		return
	}
}

func ToDTO(preCoinDTO *PreCoinStateDTO) (*usecase.CoinStateDTO, error) {
	price, err := domain.NewPriceFromString(preCoinDTO.Price)
	if err != nil {
		return nil, err
	}

	highPrice, err := domain.NewPriceFromString(preCoinDTO.HighPrice)
	if err != nil {
		return nil, err
	}

	lowPrice, err := domain.NewPriceFromString(preCoinDTO.LowPrice)
	if err != nil {
		return nil, err
	}

	coinStateDTO := &usecase.CoinStateDTO{
		Symbol:    preCoinDTO.Symbol,
		Price:     price.Int64(),
		HighPrice: highPrice.Int64(),
		LowPrice:  lowPrice.Int64(),
	}

	return coinStateDTO, nil
}
