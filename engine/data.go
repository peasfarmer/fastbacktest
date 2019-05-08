package engine

// DataHandler is the combined dataProvider interface.
type DataHandler interface {
	Next() (TickerInterface, bool)
	Stream() []TickerInterface
	History() []TickerInterface
	Latest() TickerInterface

	Reseter
}

// DataProvider is a basic dataProvider provider struct.
type DataProvider struct {
	latest TickerInterface
	stream []TickerInterface
	offset int
}

// Reset implements Reseter to reset the dataProvider struct to a clean state with loaded dataProvider events.
func (d *DataProvider) Reset() error {
	d.latest = nil
	d.offset = 0
	return nil
}

// Stream returns the dataProvider stream.
func (d *DataProvider) Stream() []TickerInterface {
	return d.stream[d.offset:]
}

// SetStream sets the dataProvider stream.
func (d *DataProvider) SetStream(stream []TickerInterface) {
	d.stream = stream
}

// Next returns the first element of the dataProvider stream,
// deletes it from the dataProvider stream and appends it to the historic dataProvider stream.
func (d *DataProvider) Next() (dh TickerInterface, ok bool) {
	// check for element in datastream
	if len(d.stream) <= d.offset {
		return dh, false
	}

	dh = d.stream[d.offset]
	d.offset++

	d.latest = dh

	return dh, true
}

// History returns the historic dataProvider stream.
func (d *DataProvider) History() []TickerInterface {
	return d.stream[:d.offset]
}

// Latest returns the last known dataProvider event for a symbol.
func (d *DataProvider) Latest() TickerInterface {
	return d.latest
}

// TickerInterface declares a dataProvider event interface
type TickerInterface interface {
	EventHandler
	Price() float64
}

// Tick declares a dataProvider event for a price tick.
type Tick struct {
	Event
	Bid       float64
	Ask       float64
	BidVolume int64
	AskVolume int64
	Last      float64
}

// Price returns the middle of Bid and Ask.
func (t Tick) Price() float64 {
	return t.Last
}

// Spread returns the difference or spread of Bid and Ask.
func (t Tick) Spread() float64 {
	return t.Bid - t.Ask
}
