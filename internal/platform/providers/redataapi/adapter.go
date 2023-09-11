package redataapi

import (
	"context"
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

	dateTruncated := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dateString := dateTruncated.Format("2006-01-02")
	startDate := dateString + "T00:00"
	endDate := dateString + "T23:59"
	prices := make([]domain.Prices, 0, len(zones))

	for _, zone := range zones {
		resBody := new(fetchPVPCPricesResponse)
		query := fetchPVPCPricesRequest{StartDate: startDate, EndDate: endDate, TimeTrunc: "hour", GeoIds: zone.ExternalID()}

		logger.DebugContext(ctx, "fetching PVPC prices from REData API", "zone", zone.Name(), "query", query)
		_, err := r.client.Path(pvpcPricesEndpoint).QueryStruct(query).Add("Accept", "application/json").ReceiveSuccess(resBody)

		if err != nil || len(resBody.Included) == 0 {
			logger.ErrorContext(ctx, "error fetching PVPC prices from REData API", "err", err, "zone", zone.Name())
			continue
		}

		pricesDto := domain.PricesDto{
			ID:   fmt.Sprintf("%s-%s", zone.ID().String(), dateString),
			Date: dateTruncated.Format(time.RFC3339),
			Zone: domain.ZoneDto{
				ID:         zone.ID().String(),
				ExternalID: zone.ExternalID(),
				Name:       zone.Name(),
			},
			Values: make([]domain.HourlyPriceDto, 0, 24),
		}

		for _, v := range resBody.Included[0].Attributes.Values {
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
