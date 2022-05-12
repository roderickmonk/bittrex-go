package archiver

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type BotConfig struct {
	Market            string
	ArchiveOrderbooks bool
	ArchiveTrades     bool
}

type BotConfigs struct {
	BotConfigs []BotConfig
}

func NewArchiver() error {

	NewMongoClient()

	RabbitConnect()
	defer RabbitClose()

	var fileName string
	flag.StringVar(&fileName, "config", "bot_configs.yaml", "Configuration File")
	flag.Parse()

	if fileName == "" {
		errMsg := "Please provide yaml file by using -config option"
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	yamlFile2, err := ioutil.ReadFile(fileName)
	if err != nil {
		errMsg :=
			fmt.Sprintf("Error reading YAML file: %s\n", err)
		return errors.New(errMsg)
	}

	var botConfigs BotConfigs
	err = yaml.Unmarshal(yamlFile2, &botConfigs)
	if err != nil {
		errMsg :=
			fmt.Sprintf("Error parsing YAML file: %s\n", err)
		return errors.New(errMsg)
	}

	var wg sync.WaitGroup

	// Orderbooks consumer (orderbooks to be archived to Mongo)
	wg.Add(1)
	go func() {
		defer wg.Done()
		BrokerReceiver(MongoClient, "orderbooks")
	}()

	// Trades consumer (trades to be archived to Mongo)
	wg.Add(1)
	go func() {
		defer wg.Done()
		BrokerReceiver(MongoClient, "trades")
	}()

	for _, botConfig := range botConfigs.BotConfigs {
		// fmt.Println(i, botConfig)
		// fmt.Println(i, botConfig.Market)
		// fmt.Println(i, botConfig.ArchiveOrderbooks)
		// fmt.Println(i, botConfig.ArchiveTrades)

		wg.Add(1)

		botConfig := botConfig

		// Bots can ingest orderbooks and trades for its allocated market
		go func() {
			defer wg.Done()
			Bot(&botConfig)
		}()

	}

	wg.Wait()

	return nil
}
