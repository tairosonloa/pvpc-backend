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
	ID      string           `db:"id"`
	Date    string           `db:"date"`
	GeoId   string           `db:"geo_id"`
	GeoName string           `db:"geo_name"`
	Prices  priceSchemaSlice `db:"values"`
}

type priceSchemaSlice []priceSchema

type priceSchema struct {
	Datetime string  `json:"datetime"`
	Price    float32 `json:"value"`
}

// Make the Attrs struct implement the driver.Value interface. This method
// simply returns the JSON-encoded representation of the struct.
func (ps priceSchemaSlice) Value() (driver.Value, error) {
	return json.Marshal(ps)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (ps *priceSchemaSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ps)
}
