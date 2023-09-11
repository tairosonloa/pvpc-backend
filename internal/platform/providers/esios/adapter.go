package esiosapi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dghubble/sling"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/domain/errors"
	"pvpc-backend/pkg/logger"
)

type EsiosAPI struct {
	client *sling.Sling
}

const (
	// pvpcPricesEndpoint is the endpoint to fetch PVPC prices from REE.
	pvpcPricesEndpoint = "/indicators/1001"
)

func NewEsiosAPI(baseUrl string) *EsiosAPI {
	return &EsiosAPI{
		client: sling.New().Base(baseUrl),
	}
}

func (r *EsiosAPI) FetchPVPCPrices(ctx context.Context, zones []domain.Zone, date time.Time) ([]domain.Prices, error) {
	if len(zones) == 0 {
		return nil, nil
	}

	zonesMapByExternalID := make(map[uint16]domain.Zone, len(zones))
	pricesDtoMap := make(map[uint16]domain.PricesDto, len(zones))
	dateTruncated := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dateString := dateTruncated.Format("2006-01-02")
	startDate := dateString + "T00:00:00"
	endDate := dateString + "T23:59:59"

	geoIDs := make([]string, 0, len(zones))
	zonesNames := make([]string, 0, len(zones))
	for _, zone := range zones {
		externalID, err := strconv.ParseUint(zone.ExternalID(), 10, 64)

		if err != nil {
			msg := "error parsing zone external ID to uint16"
			logger.ErrorContext(ctx, msg, "err", err, "zone", zone.Name())
			return nil, errors.WrapIntoDomainError(err, errors.ProviderError, msg)
		}

		geoIDs = append(geoIDs, zone.ExternalID())
		zonesNames = append(zonesNames, zone.Name())
		zonesMapByExternalID[uint16(externalID)] = zone
	}

	resBody := new(fetchPVPCPricesResponse)
	query := fetchPVPCPricesRequest{StartDate: startDate, EndDate: endDate, GeoIds: geoIDs}

	logger.DebugContext(ctx, "fetching PVPC prices from Esios", "zones", zonesNames, "query", query)
	_, err := r.client.Path(pvpcPricesEndpoint).QueryStruct(query).Add("Accept", "application/json").ReceiveSuccess(resBody)

	if err != nil || len(resBody.Indicator.Values) == 0 {
		msg := "error fetching PVPC prices from Esios API"
		logger.ErrorContext(ctx, msg, "err", err, "zones", zonesNames)
		return nil, errors.WrapIntoDomainError(err, errors.ProviderError, msg)
	}

	for _, value := range resBody.Indicator.Values {

		pricesDto, ok := pricesDtoMap[value.GeoID]

		if !ok {
			zone := zonesMapByExternalID[value.GeoID]
			pricesDto = domain.PricesDto{
				ID:   fmt.Sprintf("%s-%s", zone.ID().String(), dateString),
				Date: dateTruncated.Format(time.RFC3339),
				Zone: domain.ZoneDto{
					ID:         zone.ID().String(),
					ExternalID: zone.ExternalID(),
					Name:       zone.Name(),
				},
				Values: make([]domain.HourlyPriceDto, 0, 24),
			}
		}

		pricesDto.Values = append(pricesDto.Values, domain.HourlyPriceDto{
			Datetime: value.Datetime,
			Value:    value.Value,
		})

		pricesDtoMap[value.GeoID] = pricesDto
	}

	prices := make([]domain.Prices, 0, len(pricesDtoMap))

	for _, value := range pricesDtoMap {
		pricesDomain, err := domain.NewPrices(value)
		if err != nil {
			logger.ErrorContext(ctx, "error creating Prices domain object", "err", err, "prices", value)
			continue
		}
		prices = append(prices, pricesDomain)
	}

	return prices, nil

}
