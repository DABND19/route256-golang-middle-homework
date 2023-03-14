-- +goose Up
-- +goose StatementBegin
CREATE TABLE "carts_items" (
  user_id BIGINT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE,
  sku BIGINT NOT NULL,
  count BIGINT NOT NULL,
  PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "carts_items";
-- +goose StatementEnd
