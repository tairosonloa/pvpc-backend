package zones

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pvpc "pvpc-backend/internal"
	"pvpc-backend/internal/platform/http/common"
	"pvpc-backend/internal/services"
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

// ListZonesHandlerV1 returns a gin.HandlerFunc to list prices zones.
func ListZonesHandlerV1(listingService services.ZonesService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zones, err := listingService.ListZones(ctx)
		if err != nil {
			statusCode, response := common.NewAPIErrorResponse(err)
			ctx.JSON(statusCode, response)
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
