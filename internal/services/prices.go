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
	var today time.Time
	var allZones []domain.Zone
	var zonesToFetchToday []domain.Zone
	var zonesToFetchTomorrow []domain.Zone

	prices, err := s.pricesRepository.Query(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	if len(prices) == 0 {
		allZones, err = s.zonesRepository.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		zonesToFetchToday = allZones
	}

	now := now()

	if now.Hour() > 20 {
		if len(allZones) == 0 {
			allZones, err = s.zonesRepository.GetAll(ctx)
			if err != nil {
				return nil, err
			}
			zonesToFetchTomorrow = allZones
		} else {
			zonesToFetchTomorrow = allZones
		}
	}

	locationStr := "Europe/Madrid"
	loc, err := time.LoadLocation(locationStr)

	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("error loading %s timezone. Using server default: %s", locationStr, now.Location().String()), "err", err)
		today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	}

	for _, price := range prices {
		if price.Date().Before(today) {
			zonesToFetchToday = append(zonesToFetchToday, price.Zone())
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
