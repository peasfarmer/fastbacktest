package engine

import (
//"reflect"
)

// RiskHandler is the basic interface for accessing risks of a portfolio
type RiskHandler interface {
	EvaluateOrder(OrderEvent, TickerInterface, map[string]Position) (*Order, error)
}

// Risk is a basic risk handler implementation
type Risk struct {
}

// EvaluateOrder handles the risk of an order, refines or cancel it
func (r *Risk) EvaluateOrder(order OrderEvent, data TickerInterface, positions map[string]Position) (*Order, error) {
	// simple implementation, just gives the received order back
	// no risk management
	o := order.(*Order)
	return o, nil
}
