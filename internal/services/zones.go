package services

import (
	"context"

	pvpc "pvpc-backend/internal"
)

// ZonesService is the domain service that manages operations
// over Zone's.
type ZonesService struct {
	pricesZonesRepository pvpc.PricesZonesRepository
}

// NewZonesService returns a new ListingService.
func NewZonesService(pricesZonesRepository pvpc.PricesZonesRepository) ZonesService {
	return ZonesService{
		pricesZonesRepository: pricesZonesRepository,
	}
}

// ListZones returns the list of available PricesZone's.
func (s ZonesService) ListZones(ctx context.Context) ([]pvpc.PricesZone, error) {
	return s.pricesZonesRepository.GetAll(ctx)
}
