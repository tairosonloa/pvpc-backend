-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS prices
(
    id        CHAR(15)  PRIMARY KEY, -- GEOID-YYYY-MM-DD
    date      DATE      NOT NULL,
    geo_id    TEXT      NOT NULL,
    geo_name  TEXT      NOT NULL,
    values    JSONB     NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS prices CASCADE;
-- +goose StatementEnd
