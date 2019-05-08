package engine

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"

	//"encoding/csv"
	"errors"
	"os"

	//"io/ioutil"
	"log"
	//"os"
	//"path/filepath"
	"strconv"
	"time"
)

// TickerEventFromCSVFile loads the market dataProvider from csv files.
// It expands the underlying dataProvider struct.
type TickerEventFromCSVFile struct {
	DataProvider
	FileDir string
}

type tickCacheData struct {
	Ts        int64
	Bid       float64
	Ask       float64
	BidVolume int64
	AskVolume int64
	Last      float64
}

func (d *TickerEventFromCSVFile) saveData(symbol string, data []*Tick) (err error) {
	filePtr, err := os.Create(d.FileDir + symbol + ".db")
	if err != nil {
		fmt.Printf("Create file failed %s", err.Error())
		return
	}
	defer filePtr.Close()

	p3 := make([]*tickCacheData, len(data))
	for index, v := range data {
		p3[index] = &tickCacheData{
			Ts:        v.Time().UnixNano(),
			Bid:       v.Bid,
			Ask:       v.Ask,
			BidVolume: v.BidVolume,
			AskVolume: v.AskVolume,
			Last:      v.Last,
		}
	}

	enc := gob.NewEncoder(filePtr)
	err = enc.Encode(p3)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return err
}

func (d *TickerEventFromCSVFile) LoadCacheData(symbol string) (err error) {
	filePtr, err := os.Open(d.FileDir + symbol + ".db")
	if err != nil {
		return
	}
	defer filePtr.Close()

	var p2 []tickCacheData
	enc := gob.NewDecoder(filePtr)
	err = enc.Decode(&p2)
	if err != nil {
		return err
	}

	interfaceSlice := make([]TickerInterface, len(p2))
	for i, d := range p2 {
		date := time.Unix(0, d.Ts) //time.Parse("2006-01-02 15:04:05.999999999", d.Ts)

		event := &Event{}
		event.SetTime(date)

		ticker := &Tick{
			Event:     *event,
			Bid:       d.Bid,
			Ask:       d.Ask,
			BidVolume: d.BidVolume,
			AskVolume: d.AskVolume,
			Last:      d.Last,
		}

		interfaceSlice[i] = ticker
	}
	d.DataProvider.SetStream(interfaceSlice)
	fmt.Println("load data ok, data len:", len(interfaceSlice))

	return nil
}

func (d *TickerEventFromCSVFile) LoadCSVFile(symbol string) (err error) {
	file, err := os.Open(d.FileDir + symbol + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// create scanner on top of file
	reader := csv.NewReader(file)
	// set delimeter
	reader.Comma = ' '
	// read first line for keys and fill in array
	//reader.Read()

	var eventList []*Tick

	// read each line and create a map of values combined to the keys
	for line, err := reader.Read(); err == nil; line, err = reader.Read() {

		event, err := d.createTickerEventFromLine(line, symbol)
		if err != nil {
			// what happens if line could not be parsed - needs logging
			// log.Println(line)
			// log.Println(err)
			continue
		}

		eventList = append(eventList, event)
	}

	_ = d.saveData(symbol, eventList)

	interfaceSlice := make([]TickerInterface, len(eventList))
	for i, d := range eventList {
		interfaceSlice[i] = d
	}
	d.DataProvider.SetStream(interfaceSlice)

	fmt.Println("load data ok, data len:", len(interfaceSlice))
	return nil
}

// Load single dataProvider events into the stream ordered by date (latest first).
func (d *TickerEventFromCSVFile) Load(symbol string) (err error) {
	// check file location
	if len(d.FileDir) == 0 {
		return errors.New("no directory for dataProvider provided: ")
	}

	//try read from cache first
	err = d.LoadCacheData(symbol)
	if err == nil {
		return nil
	}

	err = d.LoadCSVFile(symbol)
	if err != nil {
		fmt.Println("load file fail", symbol, err)
		return nil
	}

	return nil
}

// createBarEventFromLine takes a key/value map and a string and builds a bar struct.
func (d *TickerEventFromCSVFile) createTickerEventFromLine(line []string, symbol string) (ticker *Tick, err error) {
	if len(line) < 6 {
		return nil, fmt.Errorf("行数据异常")
	}

	ts, err := strconv.ParseFloat(line[5], 64)
	if err != nil {
		return ticker, err
	}

	date := time.Unix(int64(ts), int64(ts*1e9)%1e9)
	//fmt.Print(date.Format("2006-01-02 15:04:05.999999999"))

	Bid, err := strconv.ParseFloat(line[0], 64)
	if err != nil {
		return ticker, err
	}

	tmp, err := strconv.ParseFloat(line[1], 64)
	if err != nil {
		return ticker, err
	}
	BidVolume := int64(tmp)

	Ask, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return ticker, err
	}

	tmp2, err := strconv.ParseFloat(line[3], 64)
	if err != nil {
		return ticker, err
	}
	AskVolume := int64(tmp2)

	Last, err := strconv.ParseFloat(line[4], 64)
	if err != nil {
		return ticker, err
	}

	// create and populate new event
	event := &Event{}
	event.SetTime(date)

	ticker = &Tick{
		Event:     *event,
		Bid:       Bid,
		Ask:       Ask,
		BidVolume: BidVolume,
		AskVolume: AskVolume,
		Last:      Last,
	}

	return ticker, nil
}
