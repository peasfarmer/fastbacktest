package engine

import (
	"math"
	"time"
	// "github.com/shopspring/decimal"
)


/*
Qty = {int64} 100
 qtyBOT = {int64} 200
 qtySLD = {int64} 100
 avgPrice = {float64} 7.413
 avgPriceNet = {float64} 7.4167
 avgPriceBOT = {float64} 7.408
 avgPriceSLD = {float64} 7.164
 value = {float64} -765.2
 valueBOT = {float64} 1481.6
 valueSLD = {float64} 716.4
 netValue = {float64} -766.299
 netValueBOT = {float64} 1482.3407
 netValueSLD = {float64} 716.0418
 marketValue = {float64} 741.3
 commission = {float64} 1.099
 exchangeFee = {float64} 0
 cost = {float64} 1.099
 costBasis = {float64} 741.6706
 realProfitLoss = {float64} -24.6282
 unrealProfitLoss = {float64} -0.3706
 totalProfitLoss = {float64} -24.9988

*/
// Position represents the holdings position
type Position struct {
	timestamp   time.Time
	Qty         int64   // 当前持仓 current Qty of the position, positive on BOT position, negativ on SLD position
	qtyBOT      int64   // how many BOT 累计多数量
	qtySLD      int64   // how many SLD 累计卖数量

	avgPrice    float64 // average price without cost
	avgPriceNet float64 // average price including cost
	avgPriceBOT float64 // average price BOT, without cost
	avgPriceSLD float64 // average price SLD, without cost

	value       float64 // Qty * price
	valueBOT    float64 // Qty BOT * price
	valueSLD    float64 // Qty SLD * price

	netValue    float64 // current value - cost
	netValueBOT float64 // current BOT value + cost
	netValueSLD float64 // current SLD value - cost
	//marketPrice float64 // last known market price
	marketValue float64 // Qty * price
	commission  float64	// 所有交易总的交易费
	exchangeFee float64	//无用
	cost        float64 // 所有交易总的(commission + fees)
	costBasis   float64 // absolute Qty * avgPriceNet

	realProfitLoss   float64 //已成交盈亏
	unrealProfitLoss float64 //浮盈
	totalProfitLoss  float64 //总盈亏=realProfitLoss+unrealProfitLoss
}

// Create a new position based on a fill event
/*func (p *Position) Create(fill *Order) {
	p.timestamp = fill.Time()

	p.update(fill)
}*/

// Update a position on a new fill event
func (p *Position) Update(fill *Order) {
	p.timestamp = fill.Time()

	p.update(fill)
}

// UpdateValue updates the current market value of a position
func (p *Position) UpdateValue(data TickerInterface) {
	p.timestamp = data.Time()

	latest := data.Price()
	p.updateValue(latest)
}

// internal function to update a position on a new fill event
func (p *Position) update(order *Order) {
	// convert order to internally used decimal numbers
	fillQty := float64(order.GetQtyFilled())
	fillPrice := order.GetAvgFillPrice()
	fillCommission := order.Commission()
	fillExchangeFee := order.ExchangeFee()
	fillCost := order.Cost()
	fillNetValue := order.NetValue()

	// convert position to internally used decimal numbers
	qty := float64(p.Qty)
	qtyBot := float64(p.qtyBOT)
	qtySld := float64(p.qtySLD)
	avgPrice := p.avgPrice
	avgPriceNet := p.avgPriceNet
	avgPriceBot := p.avgPriceBOT
	avgPriceSld := p.avgPriceSLD
	value := p.value
	valueBot := p.valueBOT
	valueSld := p.valueSLD
	netValue := p.netValue
	netValueBot := p.netValueBOT
	netValueSld := p.netValueSLD
	commission := p.commission
	exchangeFee := p.exchangeFee
	cost := p.cost
	costBasis := p.costBasis
	realProfitLoss := p.realProfitLoss

	switch order.Direction() {
	case BOT:
		if p.Qty >= 0 { // position is long, adding to position
			costBasis += fillNetValue
		} else { // position is short, closing partially out
			// costBasis + abs(fillQty) / Qty * costBasis
			costBasis += math.Abs(fillQty) / qty * costBasis
			// realProfitLoss + fillQty * (avgPriceNet - fillPrice) - fillCost
			realProfitLoss += fillQty*(avgPriceNet-fillPrice) - fillCost
		}

		// update average price for bought stock without cost
		// ( (abs(Qty) * avgPrice) + (fillQty * fillPrice) ) / (abs(Qty) + fillQty)
		avgPrice = ((math.Abs(qty) * avgPrice) + (fillQty * fillPrice)) / (math.Abs(qty) + fillQty)
		// (abs(Qty) * avgPriceNet + fillNetValue) / (abs(Qty) * fillQty)
		avgPriceNet = (math.Abs(qty)*avgPriceNet + fillNetValue) / (math.Abs(qty) + fillQty)
		// ( (Qty + avgPriceBot) + (fillQty * fillPrice) ) / fillQty
		avgPriceBot = ((qtyBot * avgPriceBot) + (fillQty * fillPrice)) / (qtyBot + fillQty)

		// update position Qty
		qty += fillQty
		qtyBot += fillQty

		// update bought value
		valueBot = qtyBot * avgPriceBot
		netValueBot += fillNetValue

	case SLD:
		if p.Qty > 0 { // position is long, closing partially out
			costBasis -= math.Abs(fillQty) / qty * costBasis
			// realProfitLoss + fillQty * (fillPrice - avgPriceNet) - fillCost
			realProfitLoss += math.Abs(fillQty)*(fillPrice-avgPriceNet) - fillCost
		} else { // position is short, adding to position
			costBasis -= fillNetValue
		}

		// update average price for bought stock without cost
		// ( (abs(Qty) * avgPrice) + (fillQty * fillPrice) ) / (abs(Qty) + fillQty)
		avgPrice = (math.Abs(qty)*avgPrice + fillQty*fillPrice) / (math.Abs(qty) + fillQty)
		// (abs(Qty) * avgPriceNet + fillNetValue) / (abs(Qty) * fillQty)
		avgPriceNet = (math.Abs(qty)*avgPriceNet + fillNetValue) / (math.Abs(qty) + fillQty)
		// avgPriceSld + (fillQty * fillPrice) / fillQty
		avgPriceSld = (qtySld*avgPriceSld + fillQty*fillPrice) / (qtySld + fillQty)

		// update position Qty
		qty -= fillQty
		qtySld += fillQty

		// update sold value
		valueSld = qtySld * avgPriceSld
		netValueSld += fillNetValue
	}

	commission += fillCommission
	exchangeFee += fillExchangeFee
	cost += fillCost
	value = valueSld - valueBot
	netValue = value - cost

	// convert from internal decimal to float
	p.Qty = int64(qty)
	p.qtyBOT = int64(qtyBot)
	p.qtySLD = int64(qtySld)
	p.avgPrice = math.Round(avgPrice*math.Pow10(DP)) / math.Pow10(DP)
	p.avgPriceBOT = math.Round(avgPriceBot*math.Pow10(DP)) / math.Pow10(DP)
	p.avgPriceSLD = math.Round(avgPriceSld*math.Pow10(DP)) / math.Pow10(DP)
	p.avgPriceNet = math.Round(avgPriceNet*math.Pow10(DP)) / math.Pow10(DP)
	p.value = math.Round(value*math.Pow10(DP)) / math.Pow10(DP)
	p.valueBOT = math.Round(valueBot*math.Pow10(DP)) / math.Pow10(DP)
	p.valueSLD = math.Round(valueSld*math.Pow10(DP)) / math.Pow10(DP)
	p.netValue = math.Round(netValue*math.Pow10(DP)) / math.Pow10(DP)
	p.netValueBOT = math.Round(netValueBot*math.Pow10(DP)) / math.Pow10(DP)
	p.netValueSLD = math.Round(netValueSld*math.Pow10(DP)) / math.Pow10(DP)
	p.commission = commission
	p.exchangeFee = exchangeFee
	p.cost = cost
	p.costBasis = math.Round(costBasis*math.Pow10(DP)) / math.Pow10(DP)
	p.realProfitLoss = math.Round(realProfitLoss*math.Pow10(DP)) / math.Pow10(DP)

	p.updateValue(order.Price())
}

// internal function to updates the current market value and profit/loss of a position
func (p *Position) updateValue(l float64) {
	// convert to internally used decimal numbers
	latest := l
	qty := float64(p.Qty)
	costBasis := p.costBasis

	// update market value
	//p.marketPrice = latest
	// abs(Qty) * current
	p.marketValue = math.Abs(qty) * latest

	// Qty * current - costBasis
	unrealProfitLoss := qty*latest - costBasis
	p.unrealProfitLoss = math.Round(unrealProfitLoss*math.Pow10(DP)) / math.Pow10(DP)

	realProfitLoss := p.realProfitLoss
	totalProfitLoss := realProfitLoss + unrealProfitLoss
	p.totalProfitLoss = math.Round(totalProfitLoss*math.Pow10(DP)) / math.Pow10(DP)
}
