package domain

import (
	"context"
	"regexp"

	"pvpc-backend/internal/domain/errors"
)

// PricesDto is the main DTO struct used to build a Prices domain entity by calling domain.NewPrices().
type PricesDto struct {
	ID     string
	Date   string
	Zone   ZoneDto
	Values []HourlyPriceDto
}

// HourlyPriceDto is the DTO struct that represents a PVPC price for a specific hour.
// Used as a part of PricesDto and only to build a Prices domain entity.
type HourlyPriceDto struct {
	Datetime string
	Value    float32
}

// Prices is the domain entity that represents PVPC prices for a day.
type Prices struct {
	id     PricesID
	date   string
	zone   Zone
	values []HourlyPrice
}

// HourlyPrice is the domain entity that represents a PVPC price for a specific hour.
// As prices for the same hour varies between zones, this entity has not meaning without a Zone,
// which is linked to the parent Prices entity.
type HourlyPrice struct {
	datetime string
	value    float32
}

// PricesID represents the Prices' unique identifier.
type PricesID struct {
	value string
}

// NewPricesID instantiate the VO for PricesID
func NewPricesID(value string) (PricesID, error) {
	err := errors.NewDomainError(errors.InvalidPricesID, "invalid Prices ID: %s. It must be in the shape of ZONE_ID-YYYY-MM-DD", value)

	if len(value) != 14 {
		return PricesID{}, err
	}

	if !regexp.MustCompile(`[A-Z]{3}-\d{4}-\d{2}-\d{2}`).MatchString(value) {
		return PricesID{}, err
	}

	return PricesID{
		value: value,
	}, nil
}

// String converts the PricesID into string.
func (id PricesID) String() string {
	return id.value
}

// PricesRepository defines the expected behavior from a prices storage.
type PricesRepository interface {
	// Save persists the given prices.
	Save(ctx context.Context, prices []Prices) error
}

// NewPrices creates a new Prices struct.
func NewPrices(pricesDto PricesDto) (Prices, error) {
	idVO, err := NewPricesID(pricesDto.ID)
	if err != nil {
		return Prices{}, err
	}

	zone, err := NewZone(pricesDto.Zone)
	if err != nil {
		return Prices{}, err
	}

	pricesValues := make([]HourlyPrice, len(pricesDto.Values))
	for i, v := range pricesDto.Values {
		pricesValues[i] = HourlyPrice{
			datetime: v.Datetime,
			value:    v.Value,
		}
	}

	prices := Prices{
		id:     idVO,
		date:   pricesDto.Date,
		zone:   zone,
		values: pricesValues,
	}

	return prices, nil
}

// ID returns the Prices' unique identifier.
func (c Prices) ID() PricesID {
	return c.id
}

// Date returns the Prices' date.
func (c Prices) Date() string {
	return c.date
}

// Zone returns the Zone' for this Prices.
func (c Prices) Zone() Zone {
	return c.zone
}

// Values returns the Prices' HourlyPrice values.
func (c Prices) Values() []HourlyPrice {
	return c.values
}

// Datetime returns the HourlyPrice's datetime.
func (p HourlyPrice) Datetime() string {
	return p.datetime
}

// Value returns the HourlyPrice's value.
func (p HourlyPrice) Value() float32 {
	return p.value
}
