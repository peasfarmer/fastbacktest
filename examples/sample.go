package main

import (
	gbt "github.com/peasfarmer/fastbacktest/engine"
)

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

	if last.Price() < 2.83 {
		portfolio.MakeOrderByFixNum(-1, 500)
	}
	if last.Price() > 2.88 {
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
