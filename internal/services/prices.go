package services

import (
	"context"
	"errors"

	"pvpc-backend/internal/domain"
)

// PricesService is the domain service that manages operations over Price's.
type PricesService struct {
	pricesRepository domain.PricesRepository
}

// NewPricesService returns a new ListingService.
func NewPricesService(pricesRepository domain.PricesRepository) PricesService {
	return PricesService{
		pricesRepository: pricesRepository,
	}
}

// FetchAndStorePricesFromREE calls REE APIs to fetch prices and stores them in the database.
func (s PricesService) FetchAndStorePricesFromREE(ctx context.Context) error {
	/*
		1. Get last prices from database
		2. Check if last prices date is less than today
		3. If true, fetch today prices from REE
		4. Store prices in database
		5. Check if time is before 20:30
		6. If true, fetch tomorrow prices from REE
		7. Store prices in database
	*/
	return errors.New("not implemented")
}
