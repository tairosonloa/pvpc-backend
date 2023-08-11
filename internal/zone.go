package pvpc

import (
	"errors"
	"fmt"
	"regexp"
)

// PricesZoneDto is the DTO structure that represents a PVPC prices zone.
type PricesZoneDto struct {
	ID         string
	ExternalId string
	Name       string
}

// PricesZone represents a PVPC prices zone.
type PricesZone struct {
	id         PricesZoneID
	externalId string
	name       string
}

// PricesZoneID represents the PVPC prices zone unique identifier.
type PricesZoneID struct {
	value string
}

var ErrInvalidPricesZoneID = errors.New("invalid PricesZone ID. It must be three capital letters")

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

// NewPricesZone creates a new PricesZone struct.
func NewPricesZone(pricesZoneDto PricesZoneDto) (PricesZone, error) {
	idVO, err := NewPricesZoneID(pricesZoneDto.ID)
	if err != nil {
		return PricesZone{}, err
	}

	pricesZone := PricesZone{
		id:         idVO,
		externalId: pricesZoneDto.ExternalId,
		name:       pricesZoneDto.Name,
	}

	return pricesZone, nil
}

// ID returns the PricesZone unique identifier.
func (c PricesZone) ID() PricesZoneID {
	return c.id
}

// Date returns the PricesZone external ID.
func (c PricesZone) ExternalId() string {
	return c.externalId
}

// GeoId returns the PricesZone Name.
func (c PricesZone) Name() string {
	return c.name
}
