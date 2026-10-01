package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"telebot/internal/domain"
	"time"
)

type CoinStateRepo struct {
	db *sql.DB
}

func NewCoinStateRepo(db *sql.DB) (*CoinStateRepo, error) {
	if db == nil {
		return nil, domain.ErrNilPointer
	}

	return &CoinStateRepo{
		db: db,
	}, nil
}

func (c *CoinStateRepo) Update(coinStates []domain.CoinState) error {
	query := `
		insert into coin_state(name, price, high_price, low_price, percent_change, last_update)
		values(?, ?, ?, ?, ?, ?);
	`

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start the transaction: %w", err)
	}
	defer tx.Rollback()

	for _, cs := range coinStates {
		symbol := cs.GetSymbol()
		price := cs.GetPrice()
		highPrice := cs.GetHighPrice()
		lowPrice := cs.GetLowPrice()
		percent := cs.GetPercantChange()
		lastUpdate := cs.GetLastUpdate()

		_, err := tx.Exec(query,
			symbol.String(),
			price.Int64(),
			highPrice.Int64(),
			lowPrice.Int64(),
			percent.Int64(),
			lastUpdate.UTC(),
		)

		if err != nil {
			return fmt.Errorf("failed to insert the coin state %s: %w", symbol.String(), err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("the transaction could not be completed: %w", err)
	}

	return nil
}

func (c *CoinStateRepo) Get() ([]domain.CoinState, error) {
	coinStates := make([]domain.CoinState, 0, 5)

	query := `
		select * from coin_state
		order by last_update desc
		limit 5;
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cryptocurrency exchange rates: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var symbol string
		var price, highPrice, lowPrice, percent int64
		var lastUpdate time.Time
		err = rows.Scan(&symbol, &price, &highPrice, &lowPrice, &percent, &lastUpdate)
		if err != nil {
			return nil, fmt.Errorf("failed scan rows: %w", err)
		}

		var coinState domain.CoinState

		err := coinState.SetSymbol(symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to asign symbol: %w", err)
		}
		err = coinState.SetPrice(price)
		if err != nil {
			return nil, fmt.Errorf("error to asign price: %w", err)
		}
		err = coinState.SetHighPrice(highPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to asign high price: %w", err)
		}
		err = coinState.SetLowPrice(lowPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to asign low price: %w", err)
		}
		err = coinState.SetPercentChange(domain.NewPercentFromInt64(percent))
		if err != nil {
			return nil, fmt.Errorf("failed to asign percent change: %w", err)
		}
		err = coinState.SetLastUpdate(lastUpdate)
		if err != nil {
			return nil, fmt.Errorf("failed to asign last update: %w", err)
		}

		coinStates = append(coinStates, coinState)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %w", err)
	}

	return coinStates, nil
}

func (c *CoinStateRepo) GetCoin(symbol string) (*domain.CoinState, error) {
	query := `
		select name, price, high_price, low_price, percent_change, last_update
		from coin_state
		where name = ?
		order by last_update desc
		limit 1;
	`

	var name string
	var price, highPrice, lowPrice, percent int64
	var lastUpdate time.Time

	err := c.db.QueryRow(query, symbol).Scan(&name, &price, &highPrice, &lowPrice, &percent, &lastUpdate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to retrieve the coin %s: %w", symbol, err)
	}

	coin := &domain.CoinState{}

	if err := coin.SetSymbol(name); err != nil {
		return nil, fmt.Errorf("failed to asign symbol: %w", err)
	}
	if err := coin.SetPrice(price); err != nil {
		return nil, fmt.Errorf("failed to asign price: %w", err)
	}
	if err := coin.SetHighPrice(highPrice); err != nil {
		return nil, fmt.Errorf("failed to asign high price: %w", err)
	}
	if err := coin.SetLowPrice(lowPrice); err != nil {
		return nil, fmt.Errorf("failed to asign low price: %w", err)
	}
	if err := coin.SetPercentChange(domain.NewPercentFromInt64(percent)); err != nil {
		return nil, fmt.Errorf("failed to asign percent change: %w", err)
	}
	if err := coin.SetLastUpdate(lastUpdate); err != nil {
		return nil, fmt.Errorf("failed to asign last update: %w", err)
	}

	return coin, nil
}
