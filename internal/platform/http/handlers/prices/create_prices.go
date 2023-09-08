package prices

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pvpc-backend/internal/platform/http/responses"
	"pvpc-backend/internal/services"
)

type response struct {
	IDs []string `json:"IDs"`
}

// CreatePricesV1 returns a gin.HandlerFunc to fetch and store PVPC prices.
func CreatePricesV1(pricesService services.PricesService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ids, err := pricesService.FetchAndStorePricesFromREE(ctx)
		if err != nil {
			statusCode, response := responses.NewAPIErrorResponse(err)
			ctx.JSON(statusCode, response)
			return
		}

		response := response{
			IDs: make([]string, len(ids)),
		}

		for i, id := range ids {
			response.IDs[i] = id.String()
		}
		ctx.JSON(http.StatusCreated, response)
	}
}