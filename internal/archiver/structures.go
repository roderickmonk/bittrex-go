package archiver

import (
	"time"
)

type Entry struct {
	Rate     string `json:"rate"`
	Quantity string `json:"quantity"`
}

type BittrexOrderbook struct {
	Bid []Entry `json:"bid"`
	Ask []Entry `json:"ask"`
}

/*
{
  "bid": [
    {
      "quantity": "number (double)",
      "rate": "number (double)"
    }
  ],
  "ask": [
    {
      "quantity": "number (double)",
      "rate": "number (double)"
    }
  ]
}
*/

type MongoOrderbookEntry struct {
	R float64 `bson:"r"`
	Q float64 `bson:"q"`
}

type MongoOrderbook struct {
	Ts     time.Time             `bson:"ts"`
	Market string                `bson:"market"`
	Bid    []MongoOrderbookEntry `bson:"bid"`
	Ask    []MongoOrderbookEntry `bson:"ask"`
}

type BittrexTrade struct {
	Id         string `json:"id"`
	ExecutedAt string `json:"executedAt"`
	R          string `json:"rate"`
	Q          string `json:"quantity"`
	TakerSide  string `json:"takerSide"`
}

type MongoTrade struct {
	Id         string    `bson:"id"`
	ExecutedAt time.Time `bson:"executedAt"`
	R          float64   `bson:"r"`
	Q          float64   `bson:"q"`
	TakerSide  string    `bson:"takerSide"`
}

type MongoTrades struct {
	Ts     time.Time    `bson:"ts"`
	Market string       `bson:"market"`
	Trades []MongoTrade `bson:"trades"`
}

type BrokerMsg struct {
	Market   string
	HttpBody []byte
}
