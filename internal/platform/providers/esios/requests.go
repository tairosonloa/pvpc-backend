package esios

type fetchPVPCPricesRequest struct {
	StartDate string   `url:"start_date"`
	EndDate   string   `url:"end_date"`
	GeoIds    []string `url:"geo_ids[]"`
}
