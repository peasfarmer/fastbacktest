package engine

import (
	// "fmt"
	"fmt"
	"time"
)

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	OnData(TickerInterface, *Backtest) (OrderEvent, error)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Commission  CommissionHandler
	ExchangeFee ExchangeFeeHandler
}

// NewExchange creates a default exchange with sensible defaults ready for use.
func NewExchange() *Exchange {
	return &Exchange{
		Commission:  &PercentageCommission{Commission: GetBackTest(nil).getConfig().Commission},
		ExchangeFee: &FixedExchangeFee{ExchangeFee: 0.0},
	}
}

// OnData executes any open order on new dataProvider
func (e *Exchange) OnData(data TickerInterface, t *Backtest) (OrderEvent, error) {
	portfolio, ok := t.portfolio.(*Portfolio)
	if !ok {
		return nil, fmt.Errorf("OnData error")
	}

	ticker := data.(*Tick)

	orders := &portfolio.orderBook
	for i := len(*orders) - 1; i >= 0; i-- {
		v := (*orders)[i]

		order, _ := v.(*Order)

		price := 0.0

		switch order.orderType {
		case LimitBuy:
			if order.limitPrice < ticker.Ask {
				continue
			}
			price = order.limitPrice
			break
		case LimitSell:
			if order.limitPrice > ticker.Bid {
				continue
			}
			price = order.limitPrice
			break
		case MarketBuy:
			price = ticker.Ask
			break
		case MarketSell:
			price = ticker.Bid
			break
		}


		order.qtyFilled = order.qty
		order.avgFillPrice = price
		order.fillTime = time.Now()
		order.status = OrderFilled




		commission, err := e.Commission.Calculate(float64(order.qtyFilled), order.avgFillPrice)
		if err != nil {
			return order, err
		}
		order.commission = commission

		exchangeFee, err := e.ExchangeFee.Fee()
		if err != nil {
			return order, err
		}
		order.exchangeFee = exchangeFee

		order.cost = commission+ exchangeFee

		transaction, err := portfolio.OnFill(order, nil)
		if err != nil {
			continue
		}
		t.statistic.TrackTransaction(transaction)

		//删除这个order
		*orders = append((*orders)[:i], (*orders)[i+1:]...)
	}

	return nil, nil
}



