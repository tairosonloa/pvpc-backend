package pvpc

import (
	"context"
	"regexp"

	"go-pvpc/internal/errors"
)

// PricesDto is the DTO structure that represents the PVPC prices for a day.
type PricesDto struct {
	ID     string
	Date   string
	Zone   PricesZoneDto
	Values []PriceDto
}

// PriceDto is the DTO structure that represents a PVPC price for a specific hour.
type PriceDto struct {
	Datetime string
	Value    float32
}

// Prices is the domain structure that represents PVPC prices for a day.
type Prices struct {
	id     PricesID
	date   string
	zone   PricesZone
	values []Price
}

// Price is the domain structure that represents a PVPC price for a specific hour.
type Price struct {
	datetime string
	value    float32
}

// PricesID represents the prices unique identifier.
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

	zone, err := NewPricesZone(pricesDto.Zone)
	if err != nil {
		return Prices{}, err
	}

	pricesValues := make([]Price, len(pricesDto.Values))
	for i, v := range pricesDto.Values {
		pricesValues[i] = Price{
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

// ID returns the Prices unique identifier.
func (c Prices) ID() PricesID {
	return c.id
}

// Date returns the Prices date.
func (c Prices) Date() string {
	return c.date
}

// Zone returns the PricesZone for this Prices.
func (c Prices) Zone() PricesZone {
	return c.zone
}

// Values returns the Prices values.
func (c Prices) Values() []Price {
	return c.values
}

// Datetime returns the Price datetime.
func (p Price) Datetime() string {
	return p.datetime
}

// Value returns the Price value.
func (p Price) Value() float32 {
	return p.value
}
