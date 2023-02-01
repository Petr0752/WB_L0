package cache

import (
	"context"
	"sync"

	"github.com/Petr0752/L0/packages/db"
	"github.com/Petr0752/L0/packages/model"
	"github.com/pkg/errors"
)

type Cache struct {
	database *db.Database
	orders   map[string]*model.Order
	mutex    sync.RWMutex
}

func New(ctx context.Context, database *db.Database) (*Cache, error) {
	orders, err := database.GetOrders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "database.GetOrders")
	}

	ordersMap := make(map[string]*model.Order, len(orders))
	for _, order := range orders {
		ordersMap[order.ID] = order
	}

	return &Cache{
		database: database,
		orders:   ordersMap,
		mutex:    sync.RWMutex{},
	}, nil
}

func (cache *Cache) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	cache.mutex.RLock()

	order, ok := cache.orders[id]
	cache.mutex.RUnlock()
	if ok {
		return order, nil
	}

	order, err := cache.database.GetOrder(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "database.GetOrder")
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.orders[order.ID] = order

	return order, nil
}

func (cache *Cache) CreateOrder(ctx context.Context, order *model.Order) error {
	err := cache.database.CreateOrder(ctx, order)
	if err != nil {
		return errors.Wrap(err, "cache.database.CreateOrder")
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.orders[order.ID] = order

	return nil
}
