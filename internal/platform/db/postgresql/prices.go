package postgresql

const (
	sqlPricesTable = "prices"
)

type sqlPrices struct {
	ID      string     `db:"id"`
	Date    string     `db:"date"`
	GeoId   string     `db:"geo_id"`
	GeoName string     `db:"geo_name"`
	Values  []sqlPrice `db:"values"`
}

type sqlPrice struct {
	Datetime string  `db:"datetime"`
	Value    float32 `db:"value"`
}
