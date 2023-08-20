package services

import (
	"context"

	"pvpc-backend/internal/domain"
)

// ZonesService is the domain service that manages operations
// over Zone's.
type ZonesService struct {
	zonesRepository domain.ZonesRepository
}

// NewZonesService returns a new ListingService.
func NewZonesService(zonesRepository domain.ZonesRepository) ZonesService {
	return ZonesService{
		zonesRepository: zonesRepository,
	}
}

// ListZones returns the list of available Zone's.
func (s ZonesService) ListZones(ctx context.Context) ([]domain.Zone, error) {
	return s.zonesRepository.GetAll(ctx)
}
