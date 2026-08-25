package usecase

import (
	"telebot/internal/domain"
)

type CoinStateDTO struct {
	Symbol    string
	Price     int64
	HighPrice int64
	LowPrice  int64
}

type CoinStateUsecase struct {
	cr CoinStateRepo
}

type CoinStateRepo interface {
	Update(coinStates []domain.CoinState) error
}

func NewCoinStateUsecase(cr CoinStateRepo) *CoinStateUsecase {
	return &CoinStateUsecase{cr: cr}
}

func (c *CoinStateUsecase) Update(coinStatesDTO []CoinStateDTO) error {
	coinStates := make([]domain.CoinState, 0, len(coinStatesDTO))

	for i := range coinStatesDTO {
		symbol, err := domain.NewSymbol(coinStatesDTO[i].Symbol)
		if err != nil {
			return nil
		}

		price, err := domain.NewPriceFromInt64(coinStatesDTO[i].Price)
		if err != nil {
			return nil
		}

		highPrice, err := domain.NewPriceFromInt64(coinStatesDTO[i].HighPrice)
		if err != nil {
			return nil
		}

		lowPrice, err := domain.NewPriceFromInt64(coinStatesDTO[i].LowPrice)
		if err != nil {
			return nil
		}

		coinState, err := domain.NewCoinState(symbol, price, highPrice, lowPrice)
		if err != nil {
			return nil
		}

		coinStates = append(coinStates, *coinState)
	}

	err := c.cr.Update(coinStates)
	if err != nil {
		return err
	}

	return nil
}
