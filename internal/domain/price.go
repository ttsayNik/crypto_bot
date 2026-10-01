package domain

import (
	"strconv"
	"strings"
)

const priceScale = 8

type Price struct {
	price int64
}

func NewPriceFromInt64(price int64) (Price, error) {
	if price < 0 {
		return Price{}, ErrPriceLessThanNull
	}

	return Price{price: price}, nil
}

func NewPriceFromString(price string) (Price, error) {
	price = strings.TrimSpace(price)

	dollar, cents, _ := strings.Cut(price, ".")

	if len(cents) < priceScale {
		cents += strings.Repeat("0", priceScale-len(cents))
	} else {
		cents = cents[:priceScale]
	}

	priceWithoutDot, err := strconv.ParseInt(dollar+cents, 10, 64)
	if err != nil {
		return Price{}, ErrInvalidPriceFormat
	}

	if priceWithoutDot < 0 {
		return Price{}, ErrPriceLessThanNull
	}

	result, err := NewPriceFromInt64(priceWithoutDot)
	if err != nil {
		return Price{}, err
	}

	return result, nil
}

func (p *Price) LessThan(price Price) bool {
	return p.price < price.price
}

func (p *Price) Int64() int64 {
	return p.price
}

func (p *Price) String() string {
	str := strconv.FormatInt(p.price, 10)

	if len(str) <= priceScale {
		str = strings.Repeat("0", priceScale-len(str)+1) + str
	}

	return str[:len(str)-priceScale] + "." + str[len(str)-priceScale:]
}
