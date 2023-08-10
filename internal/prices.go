package pvpc

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

// PricesDto is the DTO structure that represents the PVPC prices for a day.
type PricesDto struct {
	ID      string
	Date    string
	GeoId   string
	GeoName string
	Values  []PriceDto
}

// PriceDto is the DTO structure that represents a PVPC price for a specific hour.
type PriceDto struct {
	Datetime string
	Value    float32
}

// Prices is the domain structure that represents PVPC prices for a day.
type Prices struct {
	id      PricesID
	date    string
	geoId   string
	geoName string
	values  []Price
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

var ErrInvalidPricesID = errors.New("invalid Prices ID")

// NewPricesID instantiate the VO for PricesID
func NewPricesID(value string) (PricesID, error) {
	if len(value) != 15 { // GEOID-YYYY-MM-DD, where geoid is a 4 digit number
		return PricesID{}, fmt.Errorf("%w: %s", ErrInvalidPricesID, value)
	}

	re := regexp.MustCompile(`\d{4}-\d{4}-\d{2}-\d{2}`)
	matches := re.MatchString(value)

	if !matches {
		return PricesID{}, fmt.Errorf("%w: %s", ErrInvalidPricesID, value)
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
	Save(ctx context.Context, prices []Prices) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=PricesRepository

// NewPrices creates a new Prices struct.
func NewPrices(pricesDto PricesDto) (Prices, error) {
	idVO, err := NewPricesID(pricesDto.ID)
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
		id:      idVO,
		date:    pricesDto.Date,
		geoId:   pricesDto.GeoId,
		geoName: pricesDto.GeoName,
		values:  pricesValues,
	}

	return prices, nil
}

// ID returns the prices unique identifier.
func (c Prices) ID() PricesID {
	return c.id
}

// Date returns the prices date.
func (c Prices) Date() string {
	return c.date
}

// GeoId returns the prices geoId.
func (c Prices) GeoId() string {
	return c.geoId
}

// GeoName returns the prices geoName.
func (c Prices) GeoName() string {
	return c.geoName
}

// Values returns the prices values.
func (c Prices) Values() []Price {
	return c.values
}

// Datetime returns the price datetime.
func (p Price) Datetime() string {
	return p.datetime
}

// Value returns the price value.
func (p Price) Value() float32 {
	return p.value
}
