-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS zones
(
    id           CHAR(3)  PRIMARY KEY, -- THREE UPPERCASE LETTERS
    external_id  TEXT     UNIQUE NOT NULL,
    name         TEXT     NOT NULL
);

CREATE TABLE IF NOT EXISTS prices
(
    id        CHAR(14)  PRIMARY KEY, -- ZONE_ID-YYYY-MM-DD
    date      DATE      NOT NULL,
    zone_id   CHAR(3)   NOT NULL REFERENCES zones (id),
    values    JSONB     NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS zone_external_id_uniq_index ON zones (external_id);

INSERT INTO zones (id, external_id, name) VALUES
('PEN', '8741', 'Península'),
('CAN', '8742', 'Canarias'),
('BAL', '8743', 'Baleares'),
('CEU', '8744', 'Ceuta'),
('MEL', '8745', 'Melilla');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS zone_external_id_uniq_index;
DROP TABLE IF EXISTS prices CASCADE;
DROP TABLE IF EXISTS zones CASCADE;
-- +goose StatementEnd
