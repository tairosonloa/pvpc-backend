package pvpc

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

// PricesZoneDto is the DTO structure that represents a PVPC prices zone.
type PricesZoneDto struct {
	ID         string
	ExternalID string
	Name       string
}

// PricesZone represents a PVPC prices zone.
type PricesZone struct {
	id         PricesZoneID
	externalID string
	name       string
}

// PricesZoneID represents the PVPC prices zone unique identifier.
type PricesZoneID struct {
	value string
}

var ErrInvalidPricesZoneID = errors.New("invalid PricesZone ID. It must be three capital letters")
var ErrPricesZoneNotFound = errors.New("prices zone not found")

// NewPricesZoneID instantiate the VO for PricesZoneID.
func NewPricesZoneID(value string) (PricesZoneID, error) {
	if len(value) != 3 { // ID consist of 3 uppercase letters (A-Z) representing the zone.
		return PricesZoneID{}, fmt.Errorf("%w: %s", ErrInvalidPricesZoneID, value)
	}

	re := regexp.MustCompile(`[A-Z]{3}`)
	matches := re.MatchString(value)

	if !matches {
		return PricesZoneID{}, fmt.Errorf("%w: %s", ErrInvalidPricesZoneID, value)
	}

	return PricesZoneID{
		value: value,
	}, nil
}

// String converts the PricesZoneID into string.
func (id PricesZoneID) String() string {
	return id.value
}

// PricesZonesRepository defines the expected behavior from a prices storage.
type PricesZonesRepository interface {
	// GetAll returns all the prices zones.
	GetAll(ctx context.Context) ([]PricesZone, error)
	// GetByID returns the prices zone with the given ID.
	GetByID(ctx context.Context, id PricesZoneID) (PricesZone, error)
	// GetByExternalID returns the prices zone with the given external ID.
	GetByExternalID(ctx context.Context, externalID string) (PricesZone, error)
}

// NewPricesZone creates a new PricesZone struct.
func NewPricesZone(pricesZoneDto PricesZoneDto) (PricesZone, error) {
	idVO, err := NewPricesZoneID(pricesZoneDto.ID)
	if err != nil {
		return PricesZone{}, err
	}

	pricesZone := PricesZone{
		id:         idVO,
		externalID: pricesZoneDto.ExternalID,
		name:       pricesZoneDto.Name,
	}

	return pricesZone, nil
}

// ID returns the PricesZone unique identifier.
func (c PricesZone) ID() PricesZoneID {
	return c.id
}

// ExternalID returns the PricesZone external ID.
func (c PricesZone) ExternalID() string {
	return c.externalID
}

// Name returns the PricesZone Name.
func (c PricesZone) Name() string {
	return c.name
}
