package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"log"
	"os"
	"sync"
	"time"
)

// function for sarama configuration options
func newConfig() *sarama.Config {
	// enable verbose logging for Sarama to debug
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0 // for now default it to v1.0.0
	config.Net.TLS.Enable = false
	config.ClientID = "higgsboson"
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.Timeout = (10 * time.Second)
	config.Net.ReadTimeout = (10 * time.Second)
	config.Net.DialTimeout = (10 * time.Second)
	config.Net.WriteTimeout = (10 * time.Second)
	config.Metadata.Retry.Max = 5
	config.Metadata.Retry.Backoff = (1 * time.Second)
	config.Metadata.RefreshFrequency = (15 * time.Minute)
	return config
}

// setup connects to Kafka broker, and will start consuming messages
func Setup(brokers, version, group, topics string, oldest bool, redis *redis.Client) {

	log.Println("Starting up a new Sarama consumer...")

	ctx, cancel := context.WithCancel(context.Background())

	// create new consumer, set up new sarama consumer group
	config := newConfig()
	client, err := sarama.NewConsumerGroup([]string{brokers}, group, config)
	if err != nil {
		log.Panicf("Oops, there's an error creating consumer group client: %v", err)
	}

	consumer := Consumer{
		ready:       make(chan bool),
		redisclient: redis,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {

		defer wg.Done()
		for {
			log.Println("Connecting to consumer...")
			err := client.Consume(ctx, []string{topics}, &consumer)
			log.Println("Success: Connected to the consumer!")
			if err != nil {
				log.Panicf("Oops, there's an error from consumer: %v", err)
			}
			if ctx.Err() != nil { // check to see if context was cancelled, a signal for consumer to terminate
				return
			}
			consumer.ready = make(chan bool)
		}

	}()
	<-consumer.ready // this will block, until the channel is closed. At this point, nothing closes the channel
	log.Println("Score! Sarama consumer up and running...")

	cancel() // issue a cancellation

	defer func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Not so fast! There's an error closing client: %v", err)
	}
}

// consumer type to implement ConsumerGroupHandler interface
type Consumer struct {
	ready       chan bool
	redisclient *redis.Client
}

// Setup is called at beginning of new sarama session
func (consumer *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("Running setup...")
	return nil
}

// after ConsumerClaim goroutines have exited, cleanup is called at end of sarama session
func (consumer *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Running cleanup...")
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Println("Running consumeClaim...")
	for message := range claim.Messages() { // handle messages
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, offset = %d", string(message.Value), message.Timestamp, message.Topic, message.Offset)
		session.MarkMessage(message, "")
		consumer.redisclient.Set(string(message.Key), string(message.Value), 0)
	}
	return nil
}
