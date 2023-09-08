package services

import (
	"context"
	"errors"
	"time"

	"pvpc-backend/internal/domain"
)

// PricesService is the domain service that manages operations over Price's.
type PricesService struct {
	pricesRepository domain.PricesRepository
	zonesRepository  domain.ZonesRepository
}

// NewPricesService returns a new ListingService.
func NewPricesService(pricesRepository domain.PricesRepository, zonesRepository domain.ZonesRepository) PricesService {
	return PricesService{
		pricesRepository: pricesRepository,
		zonesRepository:  zonesRepository,
	}
}

// FetchAndStorePricesFromREE calls REE APIs to fetch prices and stores them in the database.
func (s PricesService) FetchAndStorePricesFromREE(ctx context.Context) error {
	var zonesToFetchToday []domain.Zone
	var zonesToFetchTomorrow []domain.Zone

	prices, err := s.pricesRepository.Query(ctx, nil, nil)
	if err != nil {
		return err
	}

	today := time.Now()

	if today.Hour() > 20 {
		zonesToFetchTomorrow, err = s.zonesRepository.GetAll(ctx)
		if err != nil {
			return err
		}
	}

	today = today.Truncate(24 * time.Hour)

	for _, price := range prices {
		if price.Date().Before(today) {
			zonesToFetchToday = append(zonesToFetchToday, price.Zone())
		}
	}

	// TODO: Fetch prices from REE
	// function that receives zonesToFetchToday and zonesToFetchTomorrow and returns a slice of prices

	// TODO: Store prices in database
	// Store the prices slice returned by the previous function into database

	return errors.New("not implemented")
}
