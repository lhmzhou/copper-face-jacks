package main

import (
	"flag"
	"github.com/go-redis/redis"
	"log"
	"copper-face-jacks/consumer"
)

var (
	brokers = ""
	version = ""
	group   = ""
	topics  = ""
	oldest  = true
)

func parseFlags() {
	flag.StringVar(&brokers, "brokers", "localhost:9092", "Kafka bootstrap brokers to connect to")
	flag.StringVar(&version, "version", "2.1.1", "kafka cluster version")
	flag.StringVar(&group, "group", "test-consumer-group", "kafka consumer group definition")
	flag.StringVar(&topics, "topics", "test_topic", "kafka topics to be consumed")
	flag.BoolVar(&oldest, "oldest", true, "kafka consumer consume initial offset from oldest")
	flag.Parse()

}

func main() {
	parseFlags()

	// make connection to redis image instance
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Panic("Oops, sorry: Unable to ping redis", err)
	}

	consumer.Setup(brokers, version, group, topics, oldest, client)

}
