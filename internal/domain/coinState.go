package domain

import (
	"time"
)

type CoinState struct {
	symbol     Symbol // name is a primary key
	price      Price
	highPrice  Price
	lowPrice   Price
	lastUpdate time.Time
}

func NewCoinState(symbol Symbol, price, highPrice, lowPrice Price) (*CoinState, error) {
	if highPrice.LessThan(lowPrice) {
		return nil, ErrHighLessThanLow
	}

	if price.LessThan(lowPrice) {
		lowPrice = price
	}

	if highPrice.LessThan(price) {
		highPrice = price
	}

	return &CoinState{
		symbol:     symbol,
		price:      price,
		highPrice:  highPrice,
		lowPrice:   lowPrice,
		lastUpdate: time.Now(),
	}, nil
}

func (c *CoinState) GetSymbol() Symbol {
	return c.symbol
}

func (c *CoinState) GetPrice() Price {
	return c.price
}

func (c *CoinState) GetHighPrice() Price {
	return c.highPrice
}

func (c *CoinState) GetLowPrice() Price {
	return c.lowPrice
}

func (c *CoinState) GetLastUpdate() time.Time {
	return c.lastUpdate
}

func (c *CoinState) SetSymbol(symbol string) error {
	sym, err := NewSymbol(symbol)
	if err != nil {
		return err
	}

	c.symbol = sym

	return nil
}

func (c *CoinState) SetPrice(price int64) error {
	prc, err := NewPriceFromInt64(price)
	if err != nil {
		return err
	}

	c.price = prc
	return nil
}

func (c *CoinState) SetHighPrice(highPrice int64) error {
	prc, err := NewPriceFromInt64(highPrice)
	if err != nil {
		return err
	}

	c.highPrice = prc
	return nil
}

func (c *CoinState) SetLowPrice(lowPrice int64) error {
	prc, err := NewPriceFromInt64(lowPrice)
	if err != nil {
		return err
	}

	c.lowPrice = prc
	return nil
}

func (c *CoinState) SetLastUpdate(lastUpdate time.Time) error {
	if lastUpdate.IsZero() {
		return ErrTimeNil
	}

	c.lastUpdate = lastUpdate

	return nil
}
