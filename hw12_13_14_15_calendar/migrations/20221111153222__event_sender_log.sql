-- +goose Up
-- +goose StatementBegin
CREATE TABLE sender_logs (
                        id VARCHAR,
                        name VARCHAR,
                        time TIMESTAMP,
                        owner_id varchar NOT NULL
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE sender_logs;