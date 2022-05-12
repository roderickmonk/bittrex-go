# Bittrex-Go

This repo uses golang to ingest the responses to a couple of HTTP Get requests to the Bittrex cyptocurrency exchange - specifically get Trades and Orderbooks - and then stores the results into two Mongo collections, namely "trades" and "orderbooks" respectively.

* The software is capable of doing this for a configurable number of markets.  Each injeston Bot runs in its own market-specific goroutine.  
* The ingestion effort is detached from the archiving effort via a local instance of RabbitMQ.
* Two further goroutines provide services to receive messages from the source bot goroutines and then store the content of the messages to MongoDB.  One such goroutine does this for orderbooks and another for trades.
* The software exploits RabbitMQ's exchange services.

* A test configuration can be found in `bot_configs.yaml`.  Some other configuration file can be provided at the command line via a `-config` flag.

* The responses from Bittrex are in pure text; hence before storing data to MongoDB, two data type conversions are applied where appropriate: 
    i. string to float64 and 
    ii. string to time 







