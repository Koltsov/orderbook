package trade

import "github.com/shopspring/decimal"

type Trade struct {
	ID     string
	Amount decimal.Decimal
	Price  decimal.Decimal
}
