package postgresql

const (
	zonesTableName = "zones"
)

type zoneSchema struct {
	ID         string `db:"id"`
	ExternalId string `db:"external_id"`
	Name       string `db:"name"`
}
