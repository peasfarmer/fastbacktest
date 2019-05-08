package engine

import "fmt"

// PortfolioHandler is the combined interface building block for a portfolio.
type PortfolioHandler interface {
	OnSignal(SignalEvent, DataHandler) (*Order, error)
	OnFill(*Order, DataHandler) (*Order, error)
	IsInvested() (Position, bool)
	IsLong() (Position, bool)
	IsShort() (Position, bool)
	Update(TickerInterface)

	InitialCash() float64
	SetInitialCash(float64)
	Cash() float64
	SetCash(float64)

	Value() float64

	Reset() error

	//make order
	MakeOrderByFixNum(price float64, num int64)

	GetPosition() Position
}

// Portfolio represent a simple portfolio struct.
type Portfolio struct {
	initialCash float64
	cash        float64
	holdings    Position
	orderBook   []OrderEvent
	//transactions []OrderEvent
	sizeManager SizeHandler
	riskManager RiskHandler
}

// NewPortfolio creates a default portfolio with sensible defaults ready for use.
func NewPortfolio() *Portfolio {
	return &Portfolio{
		initialCash: 100000,
		sizeManager: &Size{DefaultSize: 100, DefaultValue: 1000},
		riskManager: &Risk{},
	}
}

// SizeManager return the size manager of the portfolio.
func (p Portfolio) SizeManager() SizeHandler {
	return p.sizeManager
}

// SetSizeManager sets the size manager to be used with the portfolio.
func (p *Portfolio) SetSizeManager(size SizeHandler) {
	p.sizeManager = size
}

// RiskManager returns the risk manager of the portfolio.
func (p Portfolio) RiskManager() RiskHandler {
	return p.riskManager
}

// SetRiskManager sets the risk manager to be used with the portfolio.
func (p *Portfolio) SetRiskManager(risk RiskHandler) {
	p.riskManager = risk
}

// Reset the portfolio into a clean state with set initial cash.
func (p *Portfolio) Reset() error {
	p.cash = 0
	p.orderBook = []OrderEvent{}
	p.holdings = Position{}
	return nil
}

// OnSignal handles an incomming signal event
func (p *Portfolio) OnSignal(signal SignalEvent, data DataHandler) (*Order, error) {
	// fmt.Printf("Portfolio receives Signal: %#v \n", signal)
	return nil, fmt.Errorf("OnSignal不应该被执行")
	/*
		// set order type
		orderType := MarketBuy // default Market, should be set by risk manager
		var limit float64

		initialOrder := &Order{
			Event: Event{
				timestamp: signal.Time(),
				symbol:    signal.Symbol(),
			},
			//direction: signal.Direction(),
			// Qty should be set by PositionSizer
			orderType:  orderType,
			limitPrice: limit,
		}

		// fetch latest known price for the symbol
		latest := dataProvider.Latest(signal.Symbol())

		sizedOrder, err := p.sizeManager.SizeOrder(initialOrder, latest, p)
		if err != nil {
		}

		order, err := p.riskManager.EvaluateOrder(sizedOrder, latest, p.holdings)
		if err != nil {
		}

		return order, nil*/
}

// OnFill handles an incomming fill event
func (p *Portfolio) OnFill(order *Order, data DataHandler) (*Order, error) {
	// Check for nil map, else initialise the map

	p.holdings.Update(order)

	// update cash
	if order.Direction() == BOT {
		p.cash = p.cash - order.NetValue()
	} else {
		// direction is "SLD"
		p.cash = p.cash + order.NetValue()
	}

	// add order to transactions
	//p.transactions = append(p.transactions, order)

	return order, nil
}

// IsInvested checks if the portfolio has an open position on the given symbol
func (p Portfolio) IsInvested() (pos Position, ok bool) {
	pos = p.holdings
	if pos.Qty != 0 {
		return pos, true
	}
	return pos, false
}

// IsLong checks if the portfolio has an open long position on the given symbol
func (p Portfolio) IsLong() (pos Position, ok bool) {
	pos = p.holdings
	if pos.Qty > 0 {
		return pos, true
	}
	return pos, false
}

// IsShort checks if the portfolio has an open short position on the given symbol
func (p Portfolio) IsShort() (pos Position, ok bool) {
	pos = p.holdings
	if pos.Qty < 0 {
		return pos, true
	}
	return pos, false
}

// Update updates the holding on a dataProvider event
func (p *Portfolio) Update(d TickerInterface) {
	if pos, ok := p.IsInvested(); ok {
		pos.UpdateValue(d)
		p.holdings = pos
	}
}

// SetInitialCash sets the initial cash value of the portfolio
func (p *Portfolio) SetInitialCash(initial float64) {
	p.initialCash = initial
}

// InitialCash returns the initial cash value of the portfolio
func (p Portfolio) InitialCash() float64 {
	return p.initialCash
}

// SetCash sets the current cash value of the portfolio
func (p *Portfolio) SetCash(cash float64) {
	p.cash = cash
}

// Cash returns the current cash value of the portfolio
func (p Portfolio) Cash() float64 {
	return p.cash
}

// Value return the current total value of the portfolio
func (p Portfolio) Value() float64 {
	var holdingValue float64 = p.holdings.marketValue

	value := p.cash + holdingValue
	return value
}

// Holdings returns the holdings of the portfolio
func (p Portfolio) GetPosition() Position {
	return p.holdings
}

// OrderBook returns the order book of the portfolio
func (p Portfolio) OrderBook() ([]OrderEvent, bool) {
	if len(p.orderBook) == 0 {
		return p.orderBook, false
	}

	return p.orderBook, true
}

// OrdersBySymbol returns the order of a specific symbol from the order book.
func (p Portfolio) OrdersBySymbol() ([]OrderEvent, bool) {
	var orders = []OrderEvent{}

	for _, order := range p.orderBook {
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return orders, false
	}

	return orders, true
}

func (p *Portfolio) MakeOrderByFixNum(price float64, num int64) {
	position := p.holdings

	data, _ := GetBackTest(nil).GetDataProvider()
	nowTime := data.Latest().Time()

	if position.Qty > num {
		//需要卖出
		//1.撤销买单
		var orderType OrderType
		if price < 0 {
			orderType = MarketSell
		} else {
			orderType = LimitSell
		}

		//开卖
		initialOrder := &Order{
			Event: Event{
				timestamp: nowTime,
			},
			qty:        position.Qty - num,
			orderType:  orderType,
			limitPrice: price,
		}
		p.orderBook = []OrderEvent{initialOrder}

	} else if position.Qty < num {
		//需要买入
		//1.撤销卖单
		var orderType OrderType
		if price < 0 {
			orderType = MarketBuy
		} else {
			orderType = LimitBuy
		}

		//开卖
		initialOrder := &Order{
			Event: Event{
				timestamp: nowTime,
			},
			qty:        num - position.Qty,
			orderType:  orderType,
			limitPrice: price,
		}
		p.orderBook = []OrderEvent{initialOrder}
	}

}
