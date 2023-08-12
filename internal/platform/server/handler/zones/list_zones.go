package zones

import (
	pvpc "go-pvpc/internal"
	"go-pvpc/internal/listing"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Zones []zonesResponse `json:"zones"`
	Total int             `json:"total"`
}

type zonesResponse struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

// ListZonesHandler returns a gin.HandlerFunc to list prices zones.
func ListZonesHandler(listingService listing.ListingService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zones, err := listingService.ListPricesZones(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		response := mapZonesResponse(zones)
		ctx.JSON(http.StatusOK, response)
	}
}

func mapZonesResponse(zones []pvpc.PricesZone) response {
	response := response{
		Zones: make([]zonesResponse, len(zones)),
		Total: len(zones),
	}

	for i, zone := range zones {
		response.Zones[i] = zonesResponse{
			ID:         zone.ID().String(),
			ExternalID: zone.ExternalID(),
			Name:       zone.Name(),
		}
	}

	return response
}
