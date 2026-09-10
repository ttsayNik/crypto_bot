package domain

import (
	"strconv"
	"strings"
)

const scale = 8

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

	if len(cents) < scale {
		cents += strings.Repeat("0", scale-len(cents))
	} else {
		cents = cents[:scale]
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

// TODO: fix low input panic
func (p *Price) String() string {
	str := strconv.Itoa(int(p.price))

	res := str[:len(str)-scale] + "." + str[len(str)-scale:]

	return res
}
