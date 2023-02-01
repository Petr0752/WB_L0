-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    "id"   text PRIMARY KEY NOT NULL,
    "data" json             NOT NULL
);

COMMENT ON TABLE orders IS 'Заказы';
COMMENT ON COLUMN orders.id IS 'ID заказа';
COMMENT ON COLUMN orders.data IS 'Данные заказа';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
