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


```