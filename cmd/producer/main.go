package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flohansen/kafka-random-producer/internal/cli"
)

func main() {
	flags := cli.Flags{}
	flag.StringVar(&flags.Brokers, "brokers", "localhost:9092", "The Kafka brokers list (comma separated)")
	flag.StringVar(&flags.Topic, "topic", "", "The Kafka topic for which the producer should produce messages")
	flag.StringVar(&flags.Username, "username", "", "The SASL username to authenticate to the brokers")
	flag.StringVar(&flags.Password, "password", "", "The SASL password to authenticate to the brokers")
	flag.Parse()

	app := cli.Producer{}
	if err := app.Run(cli.SignalContext(), flags); err != nil {
		fmt.Printf("error running producer: %s", err)
		os.Exit(1)
	}
}
