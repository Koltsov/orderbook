package order_queue

import (
	"container/list"
)

type OrderQueue struct {
	orders *list.List
}

type OrderQueueInterface interface {
	Enqueue(orderId string) // Enqueue adds an order id to the queue
	Pop() (string, bool)    // Pop removes and returns the first order id from the queue
	IsEmpty() bool          // IsEmpty checks if the queue is empty
	Size() int              // Size returns the number of orders in the queue
	Peek() (string, bool)   // Peek returns the first order id without removing it from the queue
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		orders: list.New(),
	}
}

func (oq *OrderQueue) Enqueue(orderId string) {
	oq.orders.PushBack(orderId)
}

func (oq *OrderQueue) Pop() (string, bool) {
	if oq.orders.Len() == 0 {
		return "", false
	}
	front := oq.orders.Front()
	oq.orders.Remove(front)
	return front.Value.(string), true
}

func (oq *OrderQueue) IsEmpty() bool {
	return oq.orders.Len() == 0
}

func (oq *OrderQueue) Size() int {
	return oq.orders.Len()
}

func (oq *OrderQueue) Peek() (string, bool) {
	if oq.orders.Len() == 0 {
		return "", false
	}
	return oq.orders.Front().Value.(string), true
}
