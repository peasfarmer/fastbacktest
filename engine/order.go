package engine

import (
	"time"
)

// OrderStatus defines an order status
type OrderStatus int

// different types of order status
const (
	OrderNone OrderStatus = iota // 0
	OrderNew
	OrderSubmitted
	OrderPartiallyFilled
	OrderFilled
	OrderCanceled
	OrderCancelPending
	OrderInvalid
)

// OrderType defines which type an order is
type OrderType int

// different types of orders
const (
	MarketBuy OrderType = iota // 0
	MarketSell
	//MarketLongClose
	//MarketShortClose

	LimitBuy
	LimitSell
	//LimitLongClose
	//LimitShortClose

	//MarketOnOpenOrder
	//MarketOnCloseOrder
	//StopMarketOrder
	//LimitOrder
	//StopLimitOrder
)

// Direction defines which direction a signal indicates
type Direction int

// different types of order directions
const (
	// Buy
	BOT Direction = iota // 0
	// Sell
	SLD
	// Hold
	HLD
	// Exit
	EXT
)

// OrderEvent declares the order event interface.
type OrderEvent interface {
	EventHandler
	Directioner
	Quantifier


	IDer
	Status() OrderStatus

	Price() float64
	Commission() float64
	ExchangeFee() float64
	Cost() float64
	Value() float64
	NetValue() float64
}

// Order declares a basic order event.
type Order struct {
	Event
	id        int
	orderType OrderType // market or limit
	status    OrderStatus
	qty          int64 // quantity of the order
	qtyFilled    int64
	avgFillPrice float64
	limitPrice   float64 // limit for the order

	fillTime time.Time
	commission  float64
	exchangeFee float64
	cost        float64 // the total cost of the filled order incl commission and fees
}

// ID returns the id of the Order.
func (o Order) ID() int {
	return o.id
}

// SetID of the Order.
func (o *Order) SetID(id int) {
	o.id = id
}

// Direction returns the Direction of an Order
func (o Order) Direction() Direction {
	if o.orderType == LimitBuy || o.orderType == MarketBuy {
		return BOT
	} else {
		return SLD
	}
}

func (o *Order) SetOrderType(orderType OrderType) {
	o.orderType = orderType
}
func (o *Order) GetOrderType() (orderType OrderType) {
	return o.orderType
}

// Qty returns the Qty field of an Order
func (o Order) Qty() int64 {
	return o.qty
}

// SetQty sets the Qty field of an Order
func (o *Order) SetQty(i int64) {
	o.qty = i
}

// Status returns the status of an Order
func (o Order) Status() OrderStatus {
	return o.status
}


// Cancel cancels an order
func (o *Order) Cancel() {
	o.status = OrderCancelPending
}

// Update updates an order on a fill event
func (o *Order) Update(fill OrderEvent) {
	// not implemented
}


func (o Order) GetQtyFilled() int64 {
	return o.qtyFilled
}

func (o Order) GetAvgFillPrice() float64 {
	return o.avgFillPrice
}



// Price returns the Price field of a fill
func (o Order) Price() float64 {
	return o.avgFillPrice
}

// Commission returns the Commission field of a fill.
func (o Order) Commission() float64 {
	return o.commission
}

// ExchangeFee returns the ExchangeFee Field of a fill
func (o Order) ExchangeFee() float64 {
	return o.exchangeFee
}

// Cost returns the Cost field of a Fill
func (o Order) Cost() float64 {
	return o.cost
}

// Value returns the value without cost.
func (o Order) Value() float64 {
	value := float64(o.qtyFilled) * o.avgFillPrice
	return value
}

// NetValue returns the net value including cost.
func (o *Order) NetValue() float64 {
	if o.Direction() == BOT {
		// Qty * price + cost
		netValue := float64(o.qtyFilled)*o.avgFillPrice + o.cost
		return netValue
	}
	// SLD
	//Qty * price - cost
	netValue := float64(o.qtyFilled)*o.avgFillPrice - o.cost
	return netValue
}
