package repo

import (
	"database/sql"
	"fmt"
	"telebot/internal/domain"
	"time"
)

type CoinStateRepo struct {
	db *sql.DB
}

func NewCoinStateRepo(db *sql.DB) *CoinStateRepo {
	return &CoinStateRepo{
		db: db,
	}
}

func (c *CoinStateRepo) Update(coinStates []domain.CoinState) error {
	query := `
		insert into coin_state(name, price, high_price, low_price, last_update)
		values(?, ?, ?, ?, ?);
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
		lastUpdate := cs.GetLastUpdate()

		_, err := tx.Exec(query,
			symbol.String(),
			price.Int64(),
			highPrice.Int64(),
			lowPrice.Int64(),
			lastUpdate.UTC(),
		)

		if err != nil {
			break
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
		var price, highPrice, lowPrice int64
		var lastUpdate time.Time
		err = rows.Scan(&symbol, &price, &highPrice, &lowPrice, &lastUpdate)
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
