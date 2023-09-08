package redataapi

type fetchPVPCPricesRequest struct {
	StartDate string `url:"start_date"`
	EndDate   string `url:"end_date"`
	TimeTrunc string `url:"time_trunc"`
	GeoIds    string `url:"geo_ids"`
}
