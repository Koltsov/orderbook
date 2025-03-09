package order_book

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shopspring/decimal"
	"orderbook/internal/domain/order_queue"
	"sync"
)

// DecimalComparator for sorting TreeMap by decimal price
func DecimalComparator(a, b interface{}) int {
	ad := a.(decimal.Decimal)
	bd := b.(decimal.Decimal)
	return ad.Cmp(bd) // -1 if a < b, 1 if a > b, 0 if equal
}

type OrderBook struct {
	bids      *treemap.Map // Sorted descending (highest bid first)
	asks      *treemap.Map // Sorted ascending (lowest ask first)
	ordersMap sync.Map     // Map of orderID to Order
}

type OrderBookInterface interface {
	AddBid(orderId string, price decimal.Decimal, amount decimal.Decimal) error
	AddAsk(orderId string, price decimal.Decimal, amount decimal.Decimal) error
	BestBid() decimal.Decimal
	BestAsk() decimal.Decimal
	Bids() *treemap.Map
	Asks() *treemap.Map
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		bids:      treemap.NewWith(DecimalComparator),
		asks:      treemap.NewWith(DecimalComparator),
		ordersMap: sync.Map{},
	}
}

func (ob *OrderBook) AddBid(orderId string, price decimal.Decimal, amount decimal.Decimal) error {
	fmt.Printf("AddBid: orderId=%s, price=%s, amount=%s\n", orderId, price.String(), amount.String())
	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("amount must be greater than zero")
	}

	bidQueue, found := ob.bids.Get(price)
	if !found {
		fmt.Printf("AddBid: bid not found for price=%s\n", price.String())
		bidQueue = order_queue.NewOrderQueue()
		ob.bids.Put(price, bidQueue)
	}

	bidQueue.(order_queue.OrderQueueInterface).Enqueue(orderId)

	return nil
}

func (ob *OrderBook) AddAsk(orderId string, price decimal.Decimal, amount decimal.Decimal) error {
	fmt.Printf("AddAsk: orderId=%s, price=%s, amount=%s\n", orderId, price.String(), amount.String())
	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("amount must be greater than zero")
	}

	askQueue, found := ob.asks.Get(price)
	if !found {
		fmt.Printf("AddAsk: ask not found for price=%s\n", price.String())
		askQueue = order_queue.NewOrderQueue()
		ob.asks.Put(price, askQueue)
	}

	askQueue.(order_queue.OrderQueueInterface).Enqueue(orderId)

	return nil
}

func (ob *OrderBook) BestBid() decimal.Decimal {
	if ob.bids.Empty() {
		return decimal.Zero
	}

	bestBid, _ := ob.bids.Max()
	return bestBid.(decimal.Decimal)
}

func (ob *OrderBook) BestAsk() decimal.Decimal {
	if ob.asks.Empty() {
		return decimal.Zero
	}

	bestAsk, _ := ob.asks.Min()
	return bestAsk.(decimal.Decimal)
}

func (ob *OrderBook) Bids() *treemap.Map {
	return ob.bids
}

func (ob *OrderBook) Asks() *treemap.Map {
	return ob.asks
}

//func (ob *OrderBook) PlaceOrder(order dorder.Order) error {
//	if !order.IsPending() {
//		return fmt.Errorf("dorder is not pending")
//	}
//
//	if order.IsBuy() {
//		// Get/create the order queue for the buy side
//		buyQueue, found := ob.bids.Get(order.Price)
//		if !found {
//			buyQueue = order_queue.NewOrderQueue()
//			ob.bids.Put(order.Price, buyQueue)
//		}
//
//		// Enqueue the order
//		buyQueue.(order_queue.OrderQueueInterface).Enqueue(order)
//	} else {
//		// Get/create the order queue for the sell side
//		sellQueue, found := ob.asks.Get(order.Price)
//		if !found {
//			sellQueue = order_queue.NewOrderQueue()
//			ob.asks.Put(order.Price, sellQueue)
//		}
//
//		// Enqueue the order
//		sellQueue.(order_queue.OrderQueueInterface).Enqueue(order)
//	}
//
//	// Store the order in the orders map
//	ob.ordersMap.Store(order.ID(), order)
//
//	return nil
//}
//
//func (ob *OrderBook) CancelOrder(orderID string) error {
//	// Check if the order exists in the orders map
//	order, found := ob.ordersMap.Load(orderID)
//	if !found {
//		return fmt.Errorf("order not found")
//	}
//}
