package domain

import (
	"strconv"
	"strings"
)

const Scale = 8

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

	if len(cents) < Scale {
		cents += strings.Repeat("0", Scale-len(cents))
	} else {
		cents = cents[:Scale]
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
	str := strconv.Itoa(int(p.price))

	res := str[:len(str)-Scale] + "." + str[len(str)-Scale:]

	return res
}
