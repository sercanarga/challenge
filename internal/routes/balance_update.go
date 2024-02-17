package routes

import (
	"challenge/internal/models"
	"github.com/gin-gonic/gin"
	"time"
)

// BalanceUpdate
// @Summary		Increases or decrease the balance of the respective currency in a wallet.
// @Accept		json
// @Produce		json
// @Param		_		body		models.EventList	true "raw"
// @Success		200		{object}	models.APIReturn
// @Failure		400		{object}	models.APIReturn
// @Router		/		[post]
func BalanceUpdate(ctx *gin.Context) {
	var req models.EventList
	if ctx.BindJSON(&req) != nil {
		ctx.JSON(400, models.APIReturn{
			StatusCode:   400,
			Response:     "Invalid JSON format",
			ResponseTime: time.Now().Unix(),
		})
		return
	}

	// Validate the request
	for _, event := range req.Events {
		if event.App == "" || event.Type == "" || event.Time == "" || event.Meta.User == "" || event.Wallet == "" || event.Attributes.Amount == "" || event.Attributes.Currency == "" {
			ctx.JSON(400, models.APIReturn{
				StatusCode:   400,
				Response:     "Missing required fields",
				ResponseTime: time.Now().Unix(),
			})
			return
		}
	}

	ctx.JSON(200, models.APIReturn{
		StatusCode:   200,
		ResponseTime: time.Now().Unix(),
	})
}
