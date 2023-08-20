package services

import (
	"context"

	"pvpc-backend/internal/domain"
)

// ZonesService is the domain service that manages operations
// over Zone's.
type ZonesService struct {
	pricesZonesRepository domain.PricesZonesRepository
}

// NewZonesService returns a new ListingService.
func NewZonesService(pricesZonesRepository domain.PricesZonesRepository) ZonesService {
	return ZonesService{
		pricesZonesRepository: pricesZonesRepository,
	}
}

// ListZones returns the list of available PricesZone's.
func (s ZonesService) ListZones(ctx context.Context) ([]domain.PricesZone, error) {
	return s.pricesZonesRepository.GetAll(ctx)
}
