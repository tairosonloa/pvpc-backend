package services

import (
	"context"
	"fmt"
	"time"

	"pvpc-backend/internal/domain"
	"pvpc-backend/pkg/logger"
)

var now = time.Now

// PricesService is the domain service that manages operations over Price's.
type PricesService struct {
	pricesProvider   domain.PricesProvider
	pricesRepository domain.PricesRepository
	zonesRepository  domain.ZonesRepository
}

// NewPricesService returns a new ListingService.
func NewPricesService(
	pricesProvider domain.PricesProvider,
	pricesRepository domain.PricesRepository,
	zonesRepository domain.ZonesRepository,
) PricesService {
	return PricesService{
		pricesProvider:   pricesProvider,
		pricesRepository: pricesRepository,
		zonesRepository:  zonesRepository,
	}
}

// FetchAndStorePricesFromREE calls REE APIs to fetch prices and stores them in the database.
func (s PricesService) FetchAndStorePricesFromREE(ctx context.Context) ([]domain.PricesID, error) {
	var today, tomorrow time.Time
	var zonesToFetchToday []domain.Zone
	var zonesToFetchTomorrow []domain.Zone

	allZones, err := s.zonesRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	pricesMapByZoneID := make(map[domain.ZoneID]domain.Prices)

	allPrices, err := s.pricesRepository.Query(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	for _, prices := range allPrices {
		pricesMapByZoneID[prices.Zone().ID()] = prices
	}

	now := now()
	locationStr := "Europe/Madrid"
	loc, err := time.LoadLocation(locationStr)

	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("error loading %s timezone. Using server default: %s", locationStr, now.Location().String()), "err", err)
		today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	}
	tomorrow = today.AddDate(0, 0, 1)

	for _, zone := range allZones {
		if _, ok := pricesMapByZoneID[zone.ID()]; !ok {
			zonesToFetchToday = append(zonesToFetchToday, zone)
			if now.Hour() > 20 {
				zonesToFetchTomorrow = append(zonesToFetchTomorrow, zone)
			}
		} else {
			if pricesMapByZoneID[zone.ID()].Date().Before(today) {
				zonesToFetchToday = append(zonesToFetchToday, zone)
			}
			if now.Hour() > 20 && pricesMapByZoneID[zone.ID()].Date().Before(tomorrow) {
				zonesToFetchTomorrow = append(zonesToFetchTomorrow, zone)
			}
		}

	}

	todayCh := make(chan []domain.Prices)
	tomorrowCh := make(chan []domain.Prices)

	go func() {
		todayPrices, err := s.pricesProvider.FetchPVPCPrices(ctx, zonesToFetchToday, today)
		if err != nil {
			todayCh <- nil
			logger.ErrorContext(ctx, "error fetching today prices", "err", err)
		}
		todayCh <- todayPrices
	}()
	go func() {
		tomorrowPrices, err := s.pricesProvider.FetchPVPCPrices(ctx, zonesToFetchTomorrow, today.AddDate(0, 0, 1))
		if err != nil {
			tomorrowCh <- nil
			logger.ErrorContext(ctx, "error fetching tomorrow prices", "err", err)
		}
		tomorrowCh <- tomorrowPrices
	}()

	pricesToStore := append(<-todayCh, <-tomorrowCh...)
	if len(pricesToStore) == 0 {
		return nil, nil
	}

	err = s.pricesRepository.Save(ctx, pricesToStore)

	if err != nil {
		return nil, err
	}

	pricesIDs := make([]domain.PricesID, len(pricesToStore))
	for i, price := range pricesToStore {
		pricesIDs[i] = price.ID()
	}
	return pricesIDs, nil
}
