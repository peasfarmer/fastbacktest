package engine

import (
	"fmt"
	"log"
	"time"
)

type AppConfig struct {
	Symbol  string	//filename exclude extension name
	DataDir string	//the dir where ticker data file
	Commission float64 //fee 0.0005
	NetDelay int64 //和交易所之间的网络延迟, ms, use in make order,  order status changed notify
}

type AlgoHandler interface {
	Init() *AppConfig
	OnData(t TickerInterface,b *Backtest) (bool, error)
	OnEnd(b *Backtest)
}

func Run(algo AlgoHandler) {
	startTime := time.Now()

	// initiate new backtester
	config := algo.Init()

	app := GetBackTest(config)

	// eos100w  eos1w

	data := &TickerEventFromCSVFile{FileDir: config.DataDir}
	err := data.Load(config.Symbol)
	if err != nil {
		log.Fatal("Load DataProvider error:", err)
		return
	}

	fmt.Println("加载数据文件耗时", time.Now().Unix()-startTime.Unix())

	app.SetDataProvider(data)

	// create a new strategy with an algo stack and load into the backtest
	app.SetAlgo(algo)

	// run the backtest
	err = app.Run()
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	fmt.Println("总耗时", time.Now().Unix()-startTime.Unix())

	algo.OnEnd(app)

}
