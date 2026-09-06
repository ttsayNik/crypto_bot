package domain

import "errors"

var (
	// CoinState
	ErrHighLessThanLow = errors.New("highPrice can't be less than lowPrice")

	// Symbol
	ErrInvalidSymbolLen = errors.New("the string length must be 7-20 characters.")
	ErrInvalidSymbol    = errors.New("invalid character in the string")

	// Price
	ErrPriceLessThanNull  = errors.New("the set price is less than zero")
	ErrInvalidPriceFormat = errors.New("invalid price format")

	// Time
	ErrTimeNil = errors.New("time not specified")

	// Nil
	ErrNilPointer = errors.New("nil pointer recieved")
)
