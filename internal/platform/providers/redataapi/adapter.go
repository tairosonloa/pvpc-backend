package redataapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dghubble/sling"

	"pvpc-backend/internal/domain"
	"pvpc-backend/pkg/logger"
)

type REDataAPI struct {
	client *sling.Sling
}

const (
	// pvpcPricesEndpoint is the endpoint to fetch PVPC prices from REE.
	pvpcPricesEndpoint = "/es/datos/mercados/precios-mercados-tiempo-real"
)

func NewREDataAPI(baseUrl string) *REDataAPI {
	return &REDataAPI{
		client: sling.New().Base(baseUrl),
	}
}

func (r *REDataAPI) FetchPVPCPrices(ctx context.Context, zones []domain.Zone, date time.Time) ([]domain.Prices, error) {
	if len(zones) == 0 {
		return nil, nil
	}

	prices := make([]domain.Prices, 0, len(zones))
	startDate := date.Format("2006-01-02T00:00")
	endDate := date.Format("2006-01-02T23:59")

	for _, zone := range zones {
		req, err := r.client.Path(pvpcPricesEndpoint).
			QueryStruct(fetchPVPCPricesRequest{StartDate: startDate, EndDate: endDate, TimeTrunc: "hour", GeoIds: zone.ExternalID()}).
			Request()

		if err != nil {
			logger.ErrorContext(ctx, "error building request to fetch PVPC prices from REE", "err", err, "zone", zone.Name())
			continue
		}

		var res fetchPVPCPricesResponse
		err = json.NewDecoder(req.Body).Decode(&res)

		if err != nil {
			logger.ErrorContext(ctx, "error decoding response from REE", "err", err, "zone", zone.Name())
			continue
		}

		pricesDto := domain.PricesDto{
			ID:   fmt.Sprintf("%s-%s", zone.ID().String(), date.Format("2006-01-02")),
			Date: date.Truncate(24 * time.Hour).Format(time.RFC3339),
			Zone: domain.ZoneDto{
				ID:         zone.ID().String(),
				ExternalID: zone.ExternalID(),
				Name:       zone.Name(),
			},
			Values: make([]domain.HourlyPriceDto, 0, 24),
		}

		for _, v := range res.Included[0].Attributes.Values {
			pricesDto.Values = append(pricesDto.Values, domain.HourlyPriceDto{
				Datetime: v.Datetime,
				Value:    v.Value,
			})
		}

		pricesDomain, err := domain.NewPrices(pricesDto)
		if err != nil {
			logger.ErrorContext(ctx, "error creating Prices domain object", "err", err, "zone", zone.Name())
			continue
		}
		prices = append(prices, pricesDomain)

	}

	return prices, nil

}
