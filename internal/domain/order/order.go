package order

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"orderbook/internal/domain/trade"
	"time"
)

// Order represents a single order in the order book.
// Fields in the struct are private to enforce encapsulation.
type Order struct {
	id             string
	side           OrderSide
	customerID     string
	amount         decimal.Decimal
	price          decimal.Decimal
	status         OrderStatus
	executedAmount decimal.Decimal
	executedPrice  decimal.Decimal
	total          decimal.Decimal
	createdAt      time.Time
	trades         []trade.Trade
}

type OrderStatus string
type OrderSide string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

func NewOrder(customerID string, amount, price decimal.Decimal, side OrderSide) (*Order, error) {
	newOrderId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	return &Order{
		id:         newOrderId.String(),
		customerID: customerID,
		amount:     amount,
		price:      price,
		status:     OrderStatusPending,
		side:       side,
		createdAt:  time.Now(),
	}, nil
}

func (o *Order) ID() string {
	return o.id
}

func (o *Order) CustomerID() string {
	return o.customerID
}

func (o *Order) Amount() decimal.Decimal {
	return o.amount
}

func (o *Order) Price() decimal.Decimal {
	return o.price
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) ExecutedAmount() decimal.Decimal {
	return o.executedAmount
}

func (o *Order) ExecutedPrice() decimal.Decimal {
	return o.executedPrice
}

func (o *Order) Total() decimal.Decimal {
	return o.total
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) Trades() []trade.Trade {
	return o.trades
}

func (o *Order) Side() OrderSide {
	return o.side
}

func (o *Order) IsPending() bool {
	return o.status == OrderStatusPending
}

func (o *Order) IsBuy() bool {
	return o.side == OrderSideBuy
}

func (o *Order) IsSell() bool {
	return o.side == OrderSideSell
}

func (o *Order) IsCompleted() bool {
	return o.status == OrderStatusCompleted
}

func (o *Order) AddTrade(trade trade.Trade) {
	// Add trade to the order and calculate the executed amount and price, total, and status
	o.trades = append(o.trades, trade)
	o.executedAmount = o.executedAmount.Add(trade.Amount)
	o.executedPrice = trade.Price
	o.total = o.executedAmount.Mul(o.executedPrice)
	if o.executedAmount.Equal(o.amount) {
		o.status = OrderStatusCompleted
	}
}
