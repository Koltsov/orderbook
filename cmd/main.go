package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"orderbook/internal/application/matching_engine"
	order2 "orderbook/internal/domain/order"
)

func main() {
	// Test orderbook

	matchingEngine := matching_engine.NewMatchingEngine("BTC-USD")

	// Create a new order
	buyOrder, err := order2.NewOrder("customer1", decimal.NewFromInt(1001),
		decimal.NewFromFloat(10000.0), order2.OrderSideBuy)
	if err != nil {
		panic(err)
	}

	// Place the order in the matching engine
	err = matchingEngine.PlaceOrder(*buyOrder)
	if err != nil {
		panic(err)
	}

	// Print the order book
	matchingEngine.PrintOrderBook()

	// fetch order from engine
	order, err := matchingEngine.GetOrder(buyOrder.ID())
	if err != nil {
		panic(err)
	}

	// Print the order details
	printOrderDetails(order)

	// Create sell order
	sellOrder, err := order2.NewOrder("customer2", decimal.NewFromInt(500),
		decimal.NewFromFloat(9000.0), order2.OrderSideSell)

	if err != nil {
		panic(err)
	}

	// Place the order in the matching engine
	err = matchingEngine.PlaceOrder(*sellOrder)
	if err != nil {
		panic(err)
	}

	// Print the order book
	matchingEngine.PrintOrderBook()

	// fetch order from engine
	order, err = matchingEngine.GetOrder(sellOrder.ID())
	if err != nil {
		panic(err)
	}

	// Print the order details
	printOrderDetails(order)

	// Print details of buy order
	order, err = matchingEngine.GetOrder(buyOrder.ID())
	if err != nil {
		panic(err)
	}

	printOrderDetails(order)

	// Print the order book again
	matchingEngine.PrintOrderBook()

	return

}

func printOrderDetails(order *order2.Order) {
	fmt.Println("-- Order Details --")
	fmt.Println("Order ID:", order.ID())
	fmt.Println("Customer ID:", order.CustomerID())
	fmt.Println("Amount:", order.Amount().String())
	fmt.Println("Price:", order.Price().String())
	fmt.Println("Status:", order.Status())
	fmt.Println("Created At:", order.CreatedAt().String())
	fmt.Println("Executed Amount:", order.ExecutedAmount().String())
	fmt.Println("Executed Price:", order.ExecutedPrice().String())
	fmt.Println("Total:", order.Total().String())
	fmt.Println("-------------------")
}
