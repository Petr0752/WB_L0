package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Petr0752/L0/packages/cache"
	"github.com/Petr0752/L0/packages/config"
	"github.com/Petr0752/L0/packages/model"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

type MsgHandler struct {
	sc        stan.Conn
	cacheImpl *cache.Cache

	isStarted bool
	sub       stan.Subscription
}

func New(cacheImpl *cache.Cache) (*MsgHandler, error) {
	cfg := config.Get().NATSConfig
	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID,
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		return nil, errors.Wrap(err, "stan.Connect")
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", cfg.URL, cfg.ClusterID, cfg.ClientID)

	return &MsgHandler{
		sc:        sc,
		cacheImpl: cacheImpl,
	}, nil
}

func (handler *MsgHandler) Start() error {
	cfg := config.Get().NATSConfig
	if !handler.isStarted {
		startOpt := stan.StartWithLastReceived()
		sub, err := handler.sc.QueueSubscribe(
			cfg.Subject, cfg.QueueGroup, handler.HandleOrderMsg, startOpt, stan.DurableName(cfg.Durable),
		)
		if err != nil {
			return errors.Wrap(err, "stan.QueueSubscribe")
		}

		handler.sub = sub
	}

	return nil
}

func (handler *MsgHandler) Close() {
	if handler.isStarted {
		err := handler.sub.Unsubscribe()
		if err != nil {
			log.Printf("failed to unsubscribe from nats: %v", err)
		}
	}

	err := handler.sc.Close()
	if err != nil {
		log.Printf("failed to close nats connection: %v", err)
	}
}

func (handler *MsgHandler) HandleOrderMsg(msg *stan.Msg) {
	order := &model.Order{}
	err := json.Unmarshal(msg.Data, order)
	if err != nil {
		log.Println(err)
	}

	err = handler.cacheImpl.CreateOrder(context.Background(), order)
	if err != nil {
		log.Println(err)
	}

	log.Printf("create order with id %s\n", order.ID)
}
