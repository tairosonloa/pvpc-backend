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
	mainPricesProvider     domain.PricesProvider
	fallbackPricesProvider domain.PricesProvider
	pricesRepository       domain.PricesRepository
	zonesRepository        domain.ZonesRepository
}

// NewPricesService returns a new ListingService.
func NewPricesService(
	mainPricesProvider domain.PricesProvider,
	fallbackPricesProvider domain.PricesProvider,
	pricesRepository domain.PricesRepository,
	zonesRepository domain.ZonesRepository,
) PricesService {
	return PricesService{
		mainPricesProvider:     mainPricesProvider,
		fallbackPricesProvider: fallbackPricesProvider,
		pricesRepository:       pricesRepository,
		zonesRepository:        zonesRepository,
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
		if len(zonesToFetchToday) == 0 {
			todayCh <- nil
			return
		}
		todayPrices, err := s.mainPricesProvider.FetchPVPCPrices(ctx, zonesToFetchToday, today)
		if err != nil || len(todayPrices) == 0 {
			logger.WarnContext(ctx, "couldn't fetch today prices from main provider. Using fallback", "err", err)
			todayPrices, err = s.fallbackPricesProvider.FetchPVPCPrices(ctx, zonesToFetchToday, today)
			if err != nil || len(todayPrices) == 0 {
				logger.ErrorContext(ctx, "couldn't fetch today prices from fallback provider", "err", err)
			}
		}
		todayCh <- todayPrices
	}()

	go func() {
		if len(zonesToFetchTomorrow) == 0 {
			tomorrowCh <- nil
			return
		}
		tomorrowPrices, err := s.mainPricesProvider.FetchPVPCPrices(ctx, zonesToFetchTomorrow, tomorrow)
		if err != nil || len(tomorrowPrices) == 0 {
			logger.WarnContext(ctx, "couldn't fetch tomorrow prices from main provider. Using fallback", "err", err)
			tomorrowPrices, err = s.fallbackPricesProvider.FetchPVPCPrices(ctx, zonesToFetchTomorrow, tomorrow)
			if err != nil || len(tomorrowPrices) == 0 {
				logger.ErrorContext(ctx, "couldn't fetch tomorrow prices from fallback provider", "err", err)
			}
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
