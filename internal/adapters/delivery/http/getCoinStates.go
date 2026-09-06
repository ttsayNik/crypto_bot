package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type coinStateResponse struct {
	Symbol         string    `json:"symbol"`
	Price          string    `json:"price"`
	HighPrice      string    `json:"high"`
	LowPrice       string    `json:"low"`
	LastUpdate     time.Time `json:"last_update"`
	ProcentChanges string    `json:"procent_changes"`
}

func (h *TGClient) GetCoinStates(w http.ResponseWriter, r *http.Request) {
	resp, err := h.cu.Get()
	if err != nil {
		fmt.Println("TRUBAAAAAAAAAAAAAAAAAAAAAAAAA")
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
		log.Println("error here")
		return
	}
}
