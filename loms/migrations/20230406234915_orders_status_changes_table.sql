-- +goose Up
-- +goose StatementBegin
CREATE TABLE "orders_status_changes" (
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  submitted_at TIMESTAMP WITHOUT TIME ZONE,
  order_id BIGINT NOT NULL,
  status order_status NOT NULL,
  PRIMARY KEY (order_id, status)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "orders_status_changes";
-- +goose StatementEnd
