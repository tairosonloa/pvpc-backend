package postgresql

const (
	pricesTableName = "prices"
)

type pricesSchema struct {
	ID      string        `db:"id"`
	Date    string        `db:"date"`
	GeoId   string        `db:"geo_id"`
	GeoName string        `db:"geo_name"`
	Values  []priceSchema `db:"values"`
}

type priceSchema struct {
	Datetime string  `json:"datetime"`
	Value    float32 `json:"value"`
}
