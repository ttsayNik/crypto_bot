package usecase

import (
	"fmt"
	"log/slog"
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
	logger *slog.Logger
}

type CoinStateRepo interface {
	Update(coinStates []domain.CoinState) error
	Get() ([]domain.CoinState, error)
}

func NewCoinStateUsecase(cr CoinStateRepo, logger *slog.Logger) *CoinStateUsecase {
	return &CoinStateUsecase{cr: cr, logger: logger}
}

func (c *CoinStateUsecase) Update(coinStatesDTO []CoinStateDTO) error {
	coinStates := make([]domain.CoinState, 0, len(coinStatesDTO))

	for i := range coinStatesDTO {
		symbol, err := domain.NewSymbol(coinStatesDTO[i].Symbol)
		if err != nil {
			return fmt.Errorf("failed creating symbol from %s: %w", &symbol, err)
		}

		price, err := domain.NewPriceFromInt64(coinStatesDTO[i].Price)
		if err != nil {
			return fmt.Errorf("error creating price from %d: %w", price.Int64(), err)
		}

		highPrice, err := domain.NewPriceFromInt64(coinStatesDTO[i].HighPrice)
		if err != nil {
			return fmt.Errorf("error creating price from %d: %w", highPrice.Int64(), err)
		}

		lowPrice, err := domain.NewPriceFromInt64(coinStatesDTO[i].LowPrice)
		if err != nil {
			return fmt.Errorf("error creating price from %d: %w", lowPrice.Int64(), err)
		}

		coinState, err := domain.NewCoinState(symbol, price, highPrice, lowPrice)
		if err != nil {
			return fmt.Errorf("error creating coin state: %w", err)
		}

		coinStates = append(coinStates, *coinState)
	}

	err := c.cr.Update(coinStates)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoinStateUsecase) Get() ([]domain.CoinState, error) {
	coins, err := c.cr.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve current cryptocurrency prices") // find the repeat this error
	}

	return coins, nil
}
