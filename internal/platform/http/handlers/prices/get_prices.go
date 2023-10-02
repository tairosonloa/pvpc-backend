package prices

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"pvpc-backend/internal/domain"
	"pvpc-backend/internal/platform/http/responses"
	"pvpc-backend/internal/services"
	"pvpc-backend/pkg/logger"
)

type getPricesResponse struct {
	Prices []pricesResponse `json:"prices"`
}

type pricesResponse struct {
	Date   string                `json:"date"`
	ZoneID string                `json:"zone_id"`
	Values []hourlyPriceResponse `json:"values"`
}

type hourlyPriceResponse struct {
	Datetime string  `json:"datetime"`
	Value    float64 `json:"value"`
}

// GetPricesV1 returns a gin.HandlerFunc to retrieve prices from storage.
func GetPricesV1(pricesService services.PricesService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		zoneID, date := parseGetPricesParams(ctx, ctx.Request.URL.Query())

		prices, err := pricesService.GetPrices(ctx, zoneID, date)
		if err != nil {
			statusCode, response := responses.NewAPIErrorResponse(err)
			ctx.JSON(statusCode, response)
			return
		}

		response := getPricesResponse{
			Prices: make([]pricesResponse, len(prices)),
		}

		for i, price := range prices {
			response.Prices[i] = pricesResponse{
				Date:   price.Date().Format("2006-01-02"),
				ZoneID: price.Zone().ID().String(),
				Values: make([]hourlyPriceResponse, len(price.Values())),
			}

			for j, value := range price.Values() {
				response.Prices[i].Values[j] = hourlyPriceResponse{
					Datetime: value.Datetime().Format(time.RFC3339),
					Value:    value.Value(),
				}
			}
		}

		if len(response.Prices) == 0 {
			ctx.JSON(http.StatusNotFound, response)
		} else {
			ctx.JSON(http.StatusCreated, response)
		}
	}
}

func parseGetPricesParams(ctx context.Context, params url.Values) (*domain.ZoneID, *time.Time) {
	var zoneID *domain.ZoneID
	var date *time.Time

	for key, value := range params {
		switch key {
		case "zone_id":
			zoneID = parseZoneIDParamValue(ctx, value)
		case "date":
			date = parseDateParamValue(ctx, value)
		}
	}

	return zoneID, date
}

func parseZoneIDParamValue(ctx context.Context, zoneID []string) *domain.ZoneID {
	if len(zoneID) == 0 {
		return nil
	}
	parsedZoneID, err := domain.NewZoneID(zoneID[0])
	if err != nil {
		logger.DebugContext(ctx, "Invalid zoneID", "zoneID", zoneID[0], "err", err)
		return nil
	}
	return &parsedZoneID
}

func parseDateParamValue(ctx context.Context, date []string) *time.Time {
	if len(date) == 0 {
		return nil
	}
	parsedDate, err := time.Parse("2006-01-02", date[0])
	if err != nil {
		logger.DebugContext(ctx, "Invalid date", "date", date[0], "err", err)
		return nil
	}
	return &parsedDate
}
