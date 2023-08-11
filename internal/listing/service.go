package listing

import (
	"context"

	pvpc "go-pvpc/internal"
)

// ListingService is the domain service that manages the listing of
// resources, such as Prices and PricesZone's.
type ListingService struct {
	pricesZoneRepository pvpc.PricesZoneRepository
}

// NewListingService returns a new ListingService.
func NewListingService(pricesZoneRepository pvpc.PricesZoneRepository) ListingService {
	return ListingService{
		pricesZoneRepository: pricesZoneRepository,
	}
}

// ListPricesZones returns the list of available PricesZone's.
func (s ListingService) ListPricesZones(ctx context.Context) ([]pvpc.PricesZone, error) {
	return s.pricesZoneRepository.GetAll(ctx)
}
