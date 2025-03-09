package matching_engine

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	dorder "orderbook/internal/domain/order"
	"orderbook/internal/domain/order_book"
	"orderbook/internal/domain/order_queue"
	dtrade "orderbook/internal/domain/trade"
	"sync"
)

type MatchingEngine struct {
	symbol string
	// Order book for buy and sell orders
	orderBook order_book.OrderBookInterface
	orders    sync.Map
}

// NewMatchingEngine creates a new instance of MatchingEngine with the given symbol.
func NewMatchingEngine(symbol string) *MatchingEngine {
	ob := order_book.NewOrderBook()

	return &MatchingEngine{
		symbol:    symbol,
		orderBook: ob,
	}
}

// GetOrder retrieves an order by its ID from the matching engine.
func (e *MatchingEngine) GetOrder(orderId string) (*dorder.Order, error) {
	order, found := e.orders.Load(orderId)
	if !found {
		return nil, fmt.Errorf("order not found: %s", orderId)
	}

	return order.(*dorder.Order), nil
}

// PrintOrderBook prints the current state of the order book.
func (e *MatchingEngine) PrintOrderBook() {
	fmt.Println("-----------------------------")
	fmt.Println("Order Book for symbol:", e.symbol)
	fmt.Println("Bids:")
	e.orderBook.Bids().Each(func(key, value interface{}) {
		price := key.(decimal.Decimal)
		queue := value.(order_queue.OrderQueueInterface)
		fmt.Printf("Price: %s, Orders: %d\n", price.String(), queue.Size())
	})

	fmt.Println("++++")
	fmt.Println("Asks:")
	e.orderBook.Asks().Each(func(key, value interface{}) {
		price := key.(decimal.Decimal)
		queue := value.(order_queue.OrderQueueInterface)
		fmt.Printf("Price: %s, Orders: %v\n", price.String(), queue)
	})
	fmt.Println("-----------------------------")
}

// PlaceOrder places a new order in the order book.
func (e *MatchingEngine) PlaceOrder(order dorder.Order) error {
	// place order in the order book
	var obErr error
	if order.IsBuy() {
		obErr = e.orderBook.AddBid(order.ID(), order.Price(), order.Amount())
	} else if order.IsSell() {
		obErr = e.orderBook.AddAsk(order.ID(), order.Price(), order.Amount())
	}

	if obErr != nil {
		return fmt.Errorf("failed to place order in order book: %w", obErr)
	}

	// Store the order in the orders map
	e.orders.Store(order.ID(), &order)

	// Match orders after placing a new order
	err := e.matchOrders()
	if err != nil {
		return fmt.Errorf("failed to match orders: %w", err)
	}

	return nil
}

func (e *MatchingEngine) matchOrders() error {
	// Match orders in the order book

	bestBid := e.orderBook.BestBid()
	bestAsk := e.orderBook.BestAsk()

	if bestBid.IsZero() || bestAsk.IsZero() {
		// No orders to match
		return nil
	}

	if bestBid.GreaterThan(bestAsk) {
		// Match orders
		bq, _ := e.orderBook.Bids().Get(bestBid)
		aq, _ := e.orderBook.Asks().Get(bestAsk)

		bidQueue := bq.(order_queue.OrderQueueInterface)
		askQueue := aq.(order_queue.OrderQueueInterface)

		for !bidQueue.IsEmpty() && !askQueue.IsEmpty() {
			bidOrderId, _ := bidQueue.Peek()
			askOrderId, _ := askQueue.Peek()

			// Create trades and update order status
			bo, found := e.orders.Load(bidOrderId)
			if !found {
				return errors.New("critical: no bid order found")
			}
			ao, found := e.orders.Load(askOrderId)
			if !found {
				return errors.New("critical: no ask order found")
			}

			bidOrder := bo.(*dorder.Order)
			askOrder := ao.(*dorder.Order)

			// Execute trade
			tradeAmount := decimal.Min(bidOrder.Amount(), askOrder.Amount())
			tradePrice := bestAsk

			tradeId, err := uuid.NewRandom()
			if err != nil {
				return fmt.Errorf("critical: failed to generate trade id: %w", err)
			}
			trade := dtrade.Trade{
				ID:     tradeId.String(),
				Amount: tradeAmount,
				Price:  tradePrice,
			}

			bidOrder.AddTrade(trade)
			askOrder.AddTrade(trade)

			e.orders.Store(bidOrderId, bidOrder)
			e.orders.Store(askOrderId, askOrder)

			if bidOrder.IsCompleted() {
				bidQueue.Pop() // Remove the order from the queue
			}

			if askOrder.IsCompleted() {
				askQueue.Pop() // Remove the order from the queue
			}

			// Update order book
			if bidQueue.IsEmpty() {
				e.orderBook.Bids().Remove(bestBid)
			}

			if askQueue.IsEmpty() {
				e.orderBook.Asks().Remove(bestAsk)
			}
		}
	}

	return nil
}
