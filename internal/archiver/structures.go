package archiver

import (
	"time"
)

type Entry struct {
	R float64 `json:"r"`
	Q float64 `json:"q"`
}

type Orderbook struct {
	Ts     time.Time `json:"ts"`
	Market string    `json:"market"`
	Buy    []Entry   `json:"buy"`
	Sell   []Entry   `json:"sell"`
}

type Trade struct {
	Ts         time.Time `json:"ts"`
	Market     string    `json:"market"`
	Id         string    `json:"id"`
	ExecutedAt time.Time `json:"executedAt"`
	R          float64   `json:"r"`
	Q          float64   `json:"q"`
	TakerSide  string    `json:"takerSide"`
}
