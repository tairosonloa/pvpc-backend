package esiosapi

type fetchPVPCPricesResponse struct {
	Indicator struct {
		Values []struct {
			Value       float64 `json:"value"`
			Datetime    string  `json:"datetime"`
			DatetimeUTC string  `json:"datetime_utc"`
			TzTime      string  `json:"tz_time"`
			GeoID       uint16  `json:"geo_id"`
			GeoName     string  `json:"geo_name"`
		} `json:"values"`
	} `json:"indicator"`
}
