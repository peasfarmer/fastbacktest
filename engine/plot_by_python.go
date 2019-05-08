package engine

import (
	"encoding/json"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
	"log"
	"os"
)

type PlotByPython struct {
}

func (p *PlotByPython) SavePlotData(outputDir string, name string, configMap map[string]interface{}, data interface{}) {

	configMap["fileName"] = name + ".msgp"

	configStr, err := json.Marshal(configMap)
	if err != nil {
		fmt.Println("生成config json 错误")
		return
	}

	err = os.MkdirAll(outputDir+name, os.ModePerm)
	if err != nil {
		fmt.Println("创建目录错误")
		return
	}
	err = ioutil.WriteFile(outputDir+name+"\\config.json", configStr, os.ModePerm)
	if err != nil {
		fmt.Println("保存config.json错误")
		return
	}

	b, err := msgpack.Marshal(&data)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	err = ioutil.WriteFile(outputDir+name+"\\"+name+".msgp", b, os.ModePerm)
	if err != nil {
		fmt.Println("保存config.json错误")
		return
	}

}

func (p *PlotByPython) SaveBuySell(transcation *[]OrderEvent, path string) error {
	type row struct {
		Ts int64

		Price float64
	}

	transcationArray := make([]row, len(*transcation))
	index := 0

	for _, v := range *transcation {
		if v.Direction() != BOT {
			continue
		}

		transcationArray[index] = row{
			Ts: v.Time().UnixNano() / 1e6,

			Price: v.Price(),
		}
		index++
	}

	configMap := map[string]interface{}{
		"enable":    true,
		"plotIndex": 1,
		"mode":      "markers",
		"marker": map[string]interface{}{
			"symbol": "triangle-up",
			"size":   10,
			"color":  "green",
		},
	}

	p.SavePlotData(path, "buy", configMap, transcationArray)

	//sell
	transcationArray = make([]row, len(*transcation))
	index = 0

	for _, v := range *transcation {
		if v.Direction() != SLD {
			continue
		}

		transcationArray[index] = row{
			Ts: v.Time().UnixNano() / 1e6,

			Price: v.Price(),
		}
		index++
	}

	configMap = map[string]interface{}{
		"enable":    true,
		"plotIndex": 1,
		"mode":      "markers",
		"marker": map[string]interface{}{
			"symbol": "triangle-down",
			"size":   10,
			"color":  "red",
		},
	}

	p.SavePlotData(path, "sell", configMap, transcationArray)
	return nil
}

func (p *PlotByPython) SaveHistoryPrice(history *[]EventHandler, path string) error {
	type row struct {
		Ts    int64
		Price float64
	}

	historyArray := make([]row, len(*history)/60+1)

	for index := 0; index < len(*history); index += 60 {
		v := (*history)[index]

		t, ok := v.(*Tick)
		if !ok {
			return fmt.Errorf("格式错误")
		}

		historyArray[index/60] = row{
			Ts: t.Time().UnixNano() / 1e6,

			Price: t.Last,
		}
	}

	configMap := map[string]interface{}{
		"enable":    true,
		"plotIndex": 1,
	}

	p.SavePlotData(path, "history", configMap, historyArray)

	return nil
}

func (p *PlotByPython) SaveProfit(history *[]equityPoint, path string) error {
	type row struct {
		Ts     int64
		Equity float64
	}

	historyArray := make([]row, len(*history)/60+1)

	for index := 0; index < len(*history); index += 60 {
		v := (*history)[index]

		historyArray[index/60] = row{
			Ts: v.timestamp.UnixNano() / 1e6,

			Equity: v.equity,
		}
	}

	configMap := map[string]interface{}{
		"enable":    true,
		"plotIndex": 2,
	}

	p.SavePlotData(path, "profit", configMap, historyArray)

	return nil
}

func (p *PlotByPython) SaveDefault(app *Backtest, path string) error {

	transcation := app.GetStats().Transactions()
	p.SaveBuySell(&transcation, path)

	//历史价格
	eventHistory := app.GetStats().Events()
	p.SaveHistoryPrice(&eventHistory, path)

	equityHistory := app.GetStats().GetEquity()
	p.SaveProfit(equityHistory, path)

	return nil
}
