package engine

// DP sets the the precision of rounded floating numbers
// used after calculations to format
const DP = 4 // DP

// Reseter provides a resting interface.
type Reseter interface {
	Reset() error
}

// Backtest is the main struct which holds all elements.
type Backtest struct {
	dataProvider DataHandler
	portfolio    PortfolioHandler
	exchange     ExecutionHandler
	statistic    StatisticHandler

	algos AlgoHandler

	config *AppConfig
}

var backTest *Backtest

// New creates a default backtest with sensible defaults ready for use.
func GetBackTest(config *AppConfig) *Backtest {
	if backTest != nil {
		return backTest
	}
	backTest = &Backtest{
		config: config,
	}

	backTest.portfolio = NewPortfolio()
	backTest.exchange = NewExchange()
	backTest.statistic = &Statistic{}

	return backTest
}

// SetAlgo sets the algo stack for the Strategy
func (b *Backtest) SetAlgo(algo AlgoHandler) *Backtest {
	b.algos = algo
	return b
}

// SetDataProvider sets the dataProvider provider to be used within the backtest.
func (b *Backtest) SetDataProvider(data DataHandler) {
	b.dataProvider = data
}

func (b *Backtest) GetDataProvider() (DataHandler, bool) {
	if b.dataProvider == nil {
		return nil, false
	}

	return b.dataProvider, true
}

// SetPortfolio sets the portfolio provider to be used within the backtest.
func (b *Backtest) SetPortfolio(portfolio PortfolioHandler) {
	b.portfolio = portfolio
}

func (b *Backtest) GetPortfolio() (portfolio PortfolioHandler) {
	return b.portfolio
}

// SetExchange sets the execution provider to be used within the backtest.
func (b *Backtest) SetExchange(exchange ExecutionHandler) {
	b.exchange = exchange
}

// SetStatistic sets the statistic provider to be used within the backtest.
func (b *Backtest) SetStatistic(statistic StatisticHandler) {
	b.statistic = statistic
}

// GetStats returns the statistic handler of the backtest.
func (b *Backtest) GetStats() StatisticHandler {
	return b.statistic
}

// Reset the backtest into a clean state with loaded dataProvider.
func (b *Backtest) Reset() error {
	b.dataProvider.Reset()
	b.portfolio.Reset()
	b.statistic.Reset()
	return nil
}

// Run starts the backtest.
func (b *Backtest) Run() error {
	// setup before the backtest runs
	err := b.setup()
	if err != nil {
		return err
	}

	for ticker, ok := b.dataProvider.Next(); ok; ticker, ok = b.dataProvider.Next() {

		b.portfolio.Update(ticker)
		// update statistics
		b.statistic.Update(ticker, b.portfolio)
		// check if any orders are filled before proceding

		b.exchange.OnData(ticker, b)

		b.statistic.TrackEvent(ticker)

		b.algos.OnData(ticker, b)

	}

	// teardown at the end of the backtest
	err = b.teardown()
	if err != nil {
		return err
	}

	return nil
}

// setup runs at the beginning of the backtest to perfom preparing operations.
func (b *Backtest) setup() error {
	// before first run, set portfolio cash
	b.portfolio.SetCash(b.portfolio.InitialCash())

	return nil
}

// teardown performs any cleaning operations at the end of the backtest.
func (b *Backtest) teardown() error {
	// no implementation yet
	return nil
}

func (b *Backtest) getConfig() *AppConfig {
	return b.config
}
