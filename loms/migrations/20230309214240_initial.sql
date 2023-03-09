-- +goose Up
-- +goose StatementBegin
CREATE TYPE "order_status" AS ENUM (
  'NEW',
  'AWAITING_PAYMENT',
  'FAILED',
  'PAYED',
  'CANCELLED'
);

CREATE TABLE "orders" (
  order_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  user_id BIGINT NOT NULL,
  status order_status NOT NULL
);

CREATE TABLE "orders_items" (
  order_id BIGINT NOT NULL,
  sku BIGINT NOT NULL,
  count INTEGER,
  PRIMARY KEY (order_id, sku)
);

CREATE TABLE "items_stocks" (
  warehouse_id BIGINT NOT NULL,
  sku BIGINT NOT NULL,
  count BIGINT NOT NULL,
  PRIMARY KEY (warehouse_id, sku)
);

CREATE TABLE "items_bookings" (
  created_at TIMESTAMP WITHOUT TIME ZONE,
  order_id BIGINT NOT NULL,
  warehouse_id BIGINT NOT NULL,
  sku BIGINT NOT NULL,
  count INTEGER NOT NULL,
  PRIMARY KEY (order_id, warehouse_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "items_bookings";
DROP TABLE "items_stocks";
DROP TABLE "orders_items";
DROP TABLE "orders";
DROP TYPE "order_status";
-- +goose StatementEnd
