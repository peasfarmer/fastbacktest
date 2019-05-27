package main

import (
	gbt "github.com/peasfarmer/fastbacktest/engine"
)

type example1 struct {
}

func (s *example1) Init() *gbt.AppConfig {
	return &gbt.AppConfig{
		Symbol:     "sample",
		DataDir:    "./testdata/",
		Commission: 0.0005,
	}
}
func (s *example1) OnData(last gbt.TickerInterface, b *gbt.Backtest) (bool, error) {
	portfolio := b.GetPortfolio()

	if last.Price() < 2.83 {
		portfolio.MakeOrderByFixNum(-1, 500)
	}
	if last.Price() > 2.88 {
		portfolio.MakeOrderByFixNum(-1, 0)
	}

	return true, nil
}
func (s *example1) OnEnd(b *gbt.Backtest) {
	b.GetStats().PrintResult()

	outputDir := "./draw_data/"

	plotByPython := gbt.PlotByPython{}

	_ = plotByPython.SaveDefault(b, outputDir)
}

func main() {
	alg := &example1{}
	gbt.Run(alg)
}
