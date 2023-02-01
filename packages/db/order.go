package db

import (
	"context"

	"github.com/Petr0752/L0/packages/model"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (db *Database) CreateOrder(ctx context.Context, order *model.Order) error {
	query := "INSERT INTO orders (id, data) VALUES ($1, $2);"

	_, err := db.pool.Exec(ctx, query, order.ID, order)
	if err != nil {
		return dbError(err)
	}

	return nil
}

func (db *Database) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	query := "SELECT data FROM orders WHERE id=$1;"

	order := &model.Order{}
	err := pgxscan.Get(ctx, db.pool, order, query, id)
	if err != nil {
		return nil, dbError(err)
	}

	return order, nil
}

func (db *Database) GetOrders(ctx context.Context) ([]*model.Order, error) {
	query := "SELECT data FROM orders;"

	var orders []*model.Order
	err := pgxscan.Select(ctx, db.pool, &orders, query)
	if err != nil {
		return nil, dbError(err)
	}

	return orders, nil
}
