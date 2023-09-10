package redataapi

type fetchPVPCPricesResponse struct {
	Included []struct {
		Attributes struct {
			Values []struct {
				Value    float32 `json:"value"`
				Datetime string  `json:"datetime"`
			} `json:"values"`
		} `json:"attributes"`
	} `json:"included"`
}
