package repo

import (
	"database/sql"
	"fmt"
	"telebot/internal/domain"
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
		return err
	}
	defer tx.Rollback()

	for _, cs := range coinStates {
		symbol := cs.GetSymbol()
		price := cs.GetPrice()
		highPrice := cs.GetHighPrice()
		lowPrice := cs.GetLowPrice()
		lastUpdate := cs.GetLastUpdate()

		tx.Exec(query,
			symbol.String(),
			price.Int64(),
			highPrice.Int64(),
			lowPrice.Int64(),
			lastUpdate.UTC(),
		)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	for i := range coinStates {
		fmt.Println(coinStates[i])
	}

	fmt.Println()
	fmt.Println()

	return nil
}
