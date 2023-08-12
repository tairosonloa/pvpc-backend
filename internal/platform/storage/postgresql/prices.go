package postgresql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

const (
	pricesTableName = "prices"
)

type pricesSchema struct {
	ID     string           `db:"id"`
	Date   string           `db:"date"`
	ZoneID string           `db:"zone_id"`
	Prices priceSchemaSlice `db:"values"`
}

type priceSchemaSlice []priceSchema

type priceSchema struct {
	Datetime string  `json:"datetime"`
	Price    float32 `json:"value"`
}

// Make the priceSchemaSlice type implement the driver.Value interface.
// This method simply returns the JSON-encoded representation of the struct.
func (ps priceSchemaSlice) Value() (driver.Value, error) {
	return json.Marshal(ps)
}

// Make the priceSchemaSlice type implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (ps *priceSchemaSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("sql.Scanner Scan() custom implementation: type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ps)
}
