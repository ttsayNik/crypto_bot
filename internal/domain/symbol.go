package domain

import (
	"fmt"
	"strings"
)

type Symbol struct {
	symbol string
}

func NewSymbol(symbol string) (Symbol, error) {
	if len(symbol) < 7 || len(symbol) > 20 {
		return Symbol{}, ErrInvalidSymbolLen
	}

	uppercase := strings.ToUpper(symbol)

	fmt.Println(uppercase)

	for i := range uppercase {
		if (uppercase[i] < 'A' || uppercase[i] > 'Z') && (uppercase[i] < '0' || uppercase[i] > '9' ) {
			fmt.Println(uppercase[i])
			return Symbol{}, ErrInvalidSymbol
		}
	}

	return Symbol{symbol: symbol}, nil
}

func (s *Symbol) String() string {
	return s.symbol
}