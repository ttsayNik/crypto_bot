package domain

import (
	"math"
	"strconv"
	"strings"
)

const percentScale = 4

type Percent struct {
	percent int64
}

func NewPercent(before, after Price) (Percent, error) {
	old := before.Int64()
	if old <= 0 {
		return Percent{}, ErrZeroBasePrice
	}

	fresh := after.Int64()

	percent := float64(fresh-old) / float64(old) * float64(100*pow(percentScale))

	return Percent{percent: int64(math.Round(percent))}, nil
}

func NewPercentFromInt64(percent int64) Percent {
	return Percent{percent: percent}
}

func pow(n int) int {
	res := 1
	for i := 0; i < n; i++ {
		res *= 10
	}

	return res
}

func (p Percent) Int64() int64 {
	return p.percent
}

func (p Percent) String() string {
	percent := p.percent

	sign := ""
	if percent < 0 {
		sign = "-"
		percent = -percent
	}

	unit := int64(pow(percentScale))

	whole := strconv.FormatInt(percent/unit, 10)
	frac := strconv.FormatInt(percent%unit, 10)
	frac = strings.Repeat("0", percentScale-len(frac)) + frac

	return sign + whole + "." + frac
}
