package main

import (
	"encoding/json"
	"log"

	"github.com/Petr0752/L0/packages/config"
	"github.com/Petr0752/L0/packages/model"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

// clientID - client for publishing
const clientID = "stan-pub"

var testOrder = `
{
  "id": "some id",
  "track_number": "some track number",
  "entry": "some entry",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "some transaction id",
    "request_id": "some request id",
    "currency": "USD",
    "provider": "some provider",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "some track number",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "packages_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
`

func main() {
	// load config
	err := config.Load("configs/")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Get().NATSConfig

	sc, err := stan.Connect(cfg.ClusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, stan.DefaultNatsURL)
	}
	defer sc.Close()

	var order model.Order
	err = json.Unmarshal([]byte(testOrder), &order)
	if err != nil {
		log.Fatal(err)
	}

	// set random uid
	order.ID = uuid.NewString()

	orderJson, err := json.Marshal(&order)
	if err != nil {
		log.Fatal(err)
	}

	// send message to nats-streaming
	msg := orderJson
	err = sc.Publish(cfg.Subject, msg)
	if err != nil {
		log.Fatalf("Error during publish: %v\n", err)
	}
	log.Printf("Published [%s] : '%s'\n", cfg.Subject, msg)
}
