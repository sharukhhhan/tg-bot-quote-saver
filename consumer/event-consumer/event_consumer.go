package event_consumer

import (
	"log"
	"telegram-bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{fetcher: fetcher, processor: processor, batchSize: batchSize}
}

func (c Consumer) Start() error {
	for {
		recievedEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("error consumer: %w", err.Error())
			continue
		}

		if len(recievedEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(recievedEvents); err != nil {
			log.Print(err)

			continue
		}
	}

}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
