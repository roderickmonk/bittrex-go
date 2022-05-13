package archiver

import (
	"errors"
	"flag"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

func NewArchiver() error {

	mongoClient := NewMongoClient()

	RabbitConnect()
	defer RabbitClose()

	var fileName string
	flag.StringVar(&fileName, "config", "bot_config.yaml", "Configuration File")
	flag.Parse()

	if fileName == "" {
		failOnError(errors.New("File Not Found"), "Please provide yaml file by using -config option")
	}

	yamlFile2, err := ioutil.ReadFile(fileName)
	failOnError(err, "Error reading YAML file")

	var botConfigs BotConfigs
	err = yaml.Unmarshal(yamlFile2, &botConfigs)
	failOnError(err, "Error parsing YAML file")

	var wg sync.WaitGroup

	// Orderbooks consumer (orderbooks to be archived)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ArchiveBot(mongoClient, "orderbooks")
	}()

	// Trades consumer (trades to be archived)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ArchiveBot(mongoClient, "trades")
	}()

	// Launch the bots
	for _, botConfig := range botConfigs.BotConfigs {

		wg.Add(1)

		botConfig := botConfig

		// Bots can ingest orderbooks and trades for its allocated market
		go func() {
			defer wg.Done()
			BittrexBot(&botConfig)
		}()

	}

	wg.Wait()

	return nil
}
