package listing

import (
	"context"

	pvpc "pvpc-backend/internal"
)

// ListingService is the domain service that manages the listing of
// resources, such as Prices and PricesZone's.
type ListingService struct {
	pricesZonesRepository pvpc.PricesZonesRepository
}

// NewListingService returns a new ListingService.
func NewListingService(pricesZonesRepository pvpc.PricesZonesRepository) ListingService {
	return ListingService{
		pricesZonesRepository: pricesZonesRepository,
	}
}

// ListPricesZones returns the list of available PricesZone's.
func (s ListingService) ListPricesZones(ctx context.Context) ([]pvpc.PricesZone, error) {
	return s.pricesZonesRepository.GetAll(ctx)
}
