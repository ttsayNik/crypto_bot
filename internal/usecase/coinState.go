package usecase

import (
	"fmt"
	"log/slog"
	"telebot/internal/domain"
	"time"
)

type CoinStateDTO struct {
	Symbol    string
	Price     int64
	HighPrice int64
	LowPrice  int64
}

type CoinStateUsecase struct {
	cr     CoinStateRepo
	logger *slog.Logger

	baselines map[string]domain.CoinState
}

type CoinStateRepo interface {
	Update(coinStates []domain.CoinState) error
	Get() ([]domain.CoinState, error)
	GetCoin(symbol string) (*domain.CoinState, error)
}

func NewCoinStateUsecase(cr CoinStateRepo, logger *slog.Logger) (*CoinStateUsecase, error) {
	if cr == nil || logger == nil {
		return nil, domain.ErrNilPointer
	}

	return &CoinStateUsecase{
		cr:        cr,
		logger:    logger,
		baselines: make(map[string]domain.CoinState),
	}, nil
}

func (c *CoinStateUsecase) Update(coinStatesDTO []CoinStateDTO) error {
	coinStates := make([]domain.CoinState, 0, len(coinStatesDTO))

	for i := range coinStatesDTO {
		symbol, err := domain.NewSymbol(coinStatesDTO[i].Symbol)
		if err != nil {
			return fmt.Errorf("failed creating symbol from %s: %w", coinStatesDTO[i].Symbol, err)
		}

		price, err := domain.NewPriceFromInt64(coinStatesDTO[i].Price)
		if err != nil {
			return fmt.Errorf("error creating price from %d: %w", coinStatesDTO[i].Price, err)
		}

		highPrice, err := domain.NewPriceFromInt64(coinStatesDTO[i].HighPrice)
		if err != nil {
			return fmt.Errorf("error creating price from %d: %w", coinStatesDTO[i].HighPrice, err)
		}

		lowPrice, err := domain.NewPriceFromInt64(coinStatesDTO[i].LowPrice)
		if err != nil {
			return fmt.Errorf("error creating price from %d: %w", coinStatesDTO[i].LowPrice, err)
		}

		coinState, err := domain.NewCoinState(symbol, price, highPrice, lowPrice)
		if err != nil {
			return fmt.Errorf("error creating coin state: %w", err)
		}

		if err := coinState.SetLastUpdate(time.Now()); err != nil {
			return fmt.Errorf("error setting last update: %w", err)
		}

		base, err := c.baseline(symbol.String())
		if err != nil {
			return err
		}

		if base == nil {
			c.baselines[symbol.String()] = *coinState
			coinStates = append(coinStates, *coinState)

			continue
		}

		percent, err := domain.NewPercent(base.GetPrice(), coinState.GetPrice())
		if err != nil {
			return fmt.Errorf("error creating percent for %s: %w", symbol.String(), err)
		}

		if err := coinState.SetPercentChange(percent); err != nil {
			return fmt.Errorf("error setting percent for %s: %w", symbol.String(), err)
		}

		if time.Now().Sub(base.GetLastUpdate()) >= time.Hour {
			c.baselines[symbol.String()] = *coinState
		}

		coinStates = append(coinStates, *coinState)
	}

	err := c.cr.Update(coinStates)
	if err != nil {
		return err
	}

	return nil
}

// sameHour сообщает, относятся ли две отметки времени к одному и тому же часу.
// Сравниваются номер часа и дата, иначе 14:00 вчера и 14:00 сегодня совпали бы.
func sameHour(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()

	return aYear == bYear && aMonth == bMonth && aDay == bDay && a.Hour() == b.Hour()
}

// baseline отдаёт точку отсчёта по символу: из памяти, а после рестарта —
// из БД. nil означает, что записей по символу ещё нет.
func (c *CoinStateUsecase) baseline(symbol string) (*domain.CoinState, error) {
	if base, ok := c.baselines[symbol]; ok {
		return &base, nil
	}

	base, err := c.cr.GetCoin(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get the baseline for %s: %w", symbol, err)
	}

	if base == nil {
		return nil, nil
	}

	c.baselines[symbol] = *base

	return base, nil
}

func (c *CoinStateUsecase) Get() ([]domain.CoinState, error) {
	coins, err := c.cr.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve current cryptocurrency prices: %w", err) // find the repeat this error
	}

	return coins, nil
}

func (c *CoinStateUsecase) GetLastUpdate(symbol string) (*domain.CoinState, error) {
	coin, err := c.cr.GetCoin(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a coin: %w", err)
	}

	return coin, nil
}
