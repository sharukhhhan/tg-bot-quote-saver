package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"telegram-bot/client"
	event_consumer "telegram-bot/consumer/event-consumer"
	"telegram-bot/events"
	"telegram-bot/storage"
)

const batchSize = 1000

//var tkn string = "6171814935:AAEGumrFAw42zqPdnkfcTcNzyvPQbXHNoN4"
func main() {
	token, err := token()
	if err != nil {
		log.Fatal(err.Error())
	}

	Client := client.NewClient("api.telegram.org", token)
	s, err := storage.NewStorage("data")
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init the storage: ", err)
	}

	eventProcessor := events.NewHandler(Client, s)

	log.Println("service started")
	consumer := event_consumer.NewConsumer(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func token() (string, error) {
	token := flag.String("tg-bot-token", "", "token for the telegram bot access")
	flag.Parse()

	if *token == "" {
		return "", fmt.Errorf("token is not specified")
	}

	return *token, nil
}
