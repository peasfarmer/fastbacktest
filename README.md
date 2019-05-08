# fastbacktest
fork from https://github.com/dirkolbrich/gobacktest, and I changed a lot.
This is a framework in development, with only basic functionality.

# Usage:

Basic example:
```

type sample1 struct {
}

func (s *sample1) Init() *gbt.AppConfig {
	return &gbt.AppConfig{
		Symbol:     "sample",
		DataDir:    "./testdata/",
		Commission: 0.0005,
	}
}
func (s *sample1) OnData(last gbt.TickerInterface, b *gbt.Backtest) (bool, error) {
	portfolio := b.GetPortfolio()

	if last.Price() < 2.2 {
		portfolio.MakeOrderByFixNum(-1, 500)
	}
	if last.Price() > 2.8 {
		portfolio.MakeOrderByFixNum(-1, 0)
	}

	return true, nil
}
func (s *sample1) OnEnd(b *gbt.Backtest) {
	b.GetStats().PrintResult()

	outputDir := "./examples/draw_data/"

	plotByPython := gbt.PlotByPython{}

	_ = plotByPython.SaveDefault(b, outputDir)
}

func main() {
	alg := &sample1{}
	gbt.Run(alg)
}


```

# plot
golang plot is difficult than python, so I write a python read backtest result,and plot it use plotly(https://plot.ly/),that can show chart in browser.
you only need do is call SavePlotData function save the data you want show, the python script will read data auto. 
Basic example:
```
type row struct {
    Ts int64

    MA1 float64
    MA2 float64
    MA3 float64
}

historyArray := make([]row, len(t.tickerList))
index := 0

for _, v := range t.tickerList {
    if v.maList[2] == 0.0 {
        continue
    }
    historyArray[index] = row{
        Ts: v.t.UnixNano() / 1e6, //单位是毫秒

        MA1: v.maList[0],
        MA2: v.maList[1],
        MA3: v.maList[2],
    }
    index++
}

configMap := map[string]interface{}{
    "enable":    true,
    "plotIndex": 1,
}

plot.SavePlotData(path, "ma", configMap, historyArray)

```

![sample](/doc/plot1.png)
