-- +goose Up
ALTER TABLE feeds
ADD last_fetched_at TIMESTAMP NULL;


-- +goose Down
ALTER TABLE feeds
DROP last_fetched_at;

