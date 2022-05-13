# Bittrex-Go

## Introduction
This repo uses golang to ingest the responses to a couple of HTTP Get requests to the Bittrex cyptocurrency exchange - specifically get Trades and Orderbooks - and then stores the results into two Mongo collections, "trades" and "orderbooks".

* The software is capable of doing this for a configurable number of markets.  Each ingeston bot runs in its own market-specific goroutine.  
* The ingestion effort is detached from the archiving effort via a local instance of RabbitMQ.
* Two further goroutines provide services to receive messages from the source bot goroutines.  Each stores the content of the messages to MongoDB.
* The software exploits RabbitMQ's exchange services.

* A test (and default) configuration can be found in `bot_config.yaml`.  Some other configuration file can be provided at the command line via a `-config` flag.

* The responses from Bittrex are in pure text; hence before storing data to MongoDB, two data type conversions are applied where appropriate: 
    i. string to float64 and 
    ii. string to time 


## Requires

1. golang >1.18
1. Local Mongo instance, version >= 5.0 (port 27017)
2. Local RabbitMQ instance (port 5672)


## Installation

    git clone git@github.com:roderickmonk/bittrex-go.git
    cd bittrex-go
    go get -u ./...   # Install all required packages


## Run
    cd bittrex-go
    go run cmd/archiver/main.go












