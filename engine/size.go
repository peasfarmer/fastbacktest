package engine

import (
	"errors"
	"math"
)

// SizeHandler is the basic interface for setting the size of an order
type SizeHandler interface {
	SizeOrder(OrderEvent, TickerInterface, PortfolioHandler) (*Order, error)
}

// Size is a basic size handler implementation
type Size struct {
	DefaultSize  int64
	DefaultValue float64
}

// SizeOrder adjusts the size of an order
func (s *Size) SizeOrder(order OrderEvent, data TickerInterface, pf PortfolioHandler) (*Order, error) {
	// assert interface to concrete Type
	o := order.(*Order)
	// no default set, no sizing possible, order rejected
	if (s.DefaultSize == 0) || (s.DefaultValue == 0) {
		return o, errors.New("cannot size order: no defaultSize or defaultValue set,")
	}

	// decide on order direction
	switch o.Direction() {
	case BOT:
		o.SetQty(s.setDefaultSize(data.Price()))
	case SLD:
		o.SetQty(s.setDefaultSize(data.Price()))
	case EXT: // all shares should be sold or bought, depending on position
		// poll postions
		if _, ok := pf.IsInvested(); !ok {

			return o, errors.New("cannot exit order: no position to symbol in portfolio,")
		}
		if pos, ok := pf.IsLong(); ok {
			o.SetQty(pos.Qty)
		}
		if pos, ok := pf.IsShort(); ok {
			o.SetQty(pos.Qty * -1)
		}
	}

	return o, nil
}

func (s *Size) setDefaultSize(price float64) int64 {
	if (float64(s.DefaultSize) * price) > s.DefaultValue {
		correctedQty := int64(math.Floor(s.DefaultValue / price))
		return correctedQty
	}
	return s.DefaultSize
}
