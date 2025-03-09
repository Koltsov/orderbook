# Data structure and algorithms

Order book consists of two lists: one for buy orders and one for sell orders. Each order is represented as a tuple of (price, quantity). The buy orders are sorted in descending order of price, while the sell orders are sorted in ascending order of price.

Order books's Bids and Asks are treemap Maps, sorted by price of an order.

Each element of the Map is an Order Queue, consists of Order models.

Order queue works as FIFO, first in first out. The first order in the queue is the one that will be

When match between orders is found, the Trade is created.

Trade represents the transaction between a buy order and a sell order. It contains the price, quantity, and timestamp of the trade.

If the order was filled through multiple other orders, multiple trades will be created.

