package main

import (
	"encoding/json"
	"kaf-interface/cmd/producer/models"
	"kaf-interface/internal/producer"
	"log"
)

var (
	address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}
)

func main() {
	p, err := producer.NewProducer(address)
	if err != nil {
		log.Fatal(err)
	}
	_ = p

	orders, err := models.OredersLoad()
	if err != nil {
		log.Fatal(err)
	}

	for _, order := range orders {

		msg, _ := json.Marshal(order)

		if err := p.Produce(msg, "orders"); err != nil {
			log.Fatal(err)
		}
	}
}
