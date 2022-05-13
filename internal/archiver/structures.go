package archiver

import (
	"time"
)

type BotConfig struct {
	Market            string
	ArchiveOrderbooks bool
	ArchiveTrades     bool
}

type BotConfigs struct {
	BotConfigs []BotConfig
}

type Entry struct {
	Rate     string `json:"rate"`
	Quantity string `json:"quantity"`
}

// Describes the structure of a Bittrex Orderbook
type BittrexOrderbook struct {
	Bid []Entry `json:"bid"`
	Ask []Entry `json:"ask"`
}

type MongoOrderbookEntry struct {
	R float64 `bson:"r"`
	Q float64 `bson:"q"`
}

// Structure to be saved to the "orderbooks" collection
type MongoOrderbook struct {
	Ts     time.Time             `bson:"ts"`
	Market string                `bson:"market"`
	Bid    []MongoOrderbookEntry `bson:"bid"`
	Ask    []MongoOrderbookEntry `bson:"ask"`
}

// Describes the structure of a Bittrex Trade
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

// Structure to be saved to the "trades" collection
type MongoTrades struct {
	Ts     time.Time    `bson:"ts"`
	Market string       `bson:"market"`
	Trades []MongoTrade `bson:"trades"`
}

// Structure that is routed through RabbitMQ
type BrokerMsg struct {
	Market   string
	HttpBody []byte
}
