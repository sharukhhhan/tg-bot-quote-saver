package main

import (
	"log"
	"telegram-bot/client"
	event_consumer "telegram-bot/consumer/event-consumer"
	"telegram-bot/events"
)

const batchSize = 1000

func main() {
	//token, err := token()
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	token := "5903539225:AAG_743NMK4eSHOQu0d96HhBaE1t7upi8mc"
	Client := client.NewClient("api.telegram.org", token)
	eventProcessor := events.NewHandler(Client)

	log.Println("service started")
	consumer := event_consumer.NewConsumer(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

//func token() (string, error) {
//	token := flag.String("tg-bot-token", "", "token for the telegram bot access")
//	flag.Parse()
//
//	if *token == "" {
//		return "", fmt.Errorf("token is not specified")
//	}
//
//	return *token, nil
//}
