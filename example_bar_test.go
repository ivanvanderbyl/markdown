package markdown

import (
	"fmt"
	"os"
	"time"
)

// Bar is an aggregate of trades.
type Bar struct {
	Timestamp  time.Time
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     uint64
	TradeCount uint64
	VWAP       float64
}

// ExampleMarkdown_Table_bars shows how to turn seven daily Bar values into a markdown table.
func ExampleMarkdown_Table_bars() {
	bars := []Bar{
		{Timestamp: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC), Open: 101.25, High: 105.50, Low: 100.90, Close: 104.20, Volume: 1200345, TradeCount: 3456, VWAP: 103.45},
		{Timestamp: time.Date(2024, 10, 2, 0, 0, 0, 0, time.UTC), Open: 104.20, High: 106.80, Low: 103.75, Close: 105.10, Volume: 980456, TradeCount: 2980, VWAP: 104.95},
		{Timestamp: time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC), Open: 105.10, High: 107.20, Low: 104.10, Close: 106.75, Volume: 1100456, TradeCount: 3104, VWAP: 106.15},
		{Timestamp: time.Date(2024, 10, 4, 0, 0, 0, 0, time.UTC), Open: 106.75, High: 108.90, Low: 105.30, Close: 108.40, Volume: 1023400, TradeCount: 2890, VWAP: 107.85},
		{Timestamp: time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC), Open: 108.40, High: 109.25, Low: 106.80, Close: 107.10, Volume: 954320, TradeCount: 2605, VWAP: 107.35},
		{Timestamp: time.Date(2024, 10, 6, 0, 0, 0, 0, time.UTC), Open: 107.10, High: 108.75, Low: 106.40, Close: 108.20, Volume: 876540, TradeCount: 2400, VWAP: 107.95},
		{Timestamp: time.Date(2024, 10, 7, 0, 0, 0, 0, time.UTC), Open: 108.20, High: 110.15, Low: 107.95, Close: 109.60, Volume: 1132050, TradeCount: 3250, VWAP: 109.05},
	}

	rows := make([][]string, len(bars))
	for i, bar := range bars {
		rows[i] = []string{
			bar.Timestamp.Format("2006-01-02"),
			fmt.Sprintf("%.2f", bar.Open),
			fmt.Sprintf("%.2f", bar.High),
			fmt.Sprintf("%.2f", bar.Low),
			fmt.Sprintf("%.2f", bar.Close),
			fmt.Sprintf("%d", bar.Volume),
			fmt.Sprintf("%d", bar.TradeCount),
			fmt.Sprintf("%.2f", bar.VWAP),
		}
	}

	md := NewMarkdown(os.Stdout)
	md.H2("Daily Bars")
	md.Table(TableSet{
		Header: []string{"Day", "Open", "High", "Low", "Close", "Volume", "Trades", "VWAP"},
		Rows:   rows,
	})

	if err := md.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "Error building markdown: %v\n", err)
		return
	}

	// Output:
	// ## Daily Bars
	// | Day        | Open   | High   | Low    | Close  | Volume  | Trades | VWAP   |
	// | ---------- | ------ | ------ | ------ | ------ | ------- | ------ | ------ |
	// | 2024-10-01 | 101.25 | 105.50 | 100.90 | 104.20 | 1200345 | 3456   | 103.45 |
	// | 2024-10-02 | 104.20 | 106.80 | 103.75 | 105.10 | 980456  | 2980   | 104.95 |
	// | 2024-10-03 | 105.10 | 107.20 | 104.10 | 106.75 | 1100456 | 3104   | 106.15 |
	// | 2024-10-04 | 106.75 | 108.90 | 105.30 | 108.40 | 1023400 | 2890   | 107.85 |
	// | 2024-10-05 | 108.40 | 109.25 | 106.80 | 107.10 | 954320  | 2605   | 107.35 |
	// | 2024-10-06 | 107.10 | 108.75 | 106.40 | 108.20 | 876540  | 2400   | 107.95 |
	// | 2024-10-07 | 108.20 | 110.15 | 107.95 | 109.60 | 1132050 | 3250   | 109.05 |
}
