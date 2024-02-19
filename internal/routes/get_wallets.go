package routes

import (
	"challenge/internal/durable"
	"challenge/internal/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// GetWallets
// @Summary		Returns a list of wallets.
// @Produce		json
// @Param		limit		query		int		false	"default:10"
// @Param		cursor		query		int		false	"default:0"
// @Success		200		{object}	models.WalletStruct
// @Failure		500		{object}	models.APIReturn
// @Router		/		[get]
func GetWallets(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))
	cursor = cursor * limit

	var wallets []models.Wallet
	result := durable.Connection().Preload("Balance").Offset(cursor).Limit(limit).Find(&wallets)
	if result.Error != nil {
		ctx.JSON(500, models.APIReturn{
			StatusCode:   500,
			Response:     "Failed to retrieve wallets",
			ResponseTime: time.Now().Unix(),
		})
		return
	}

	var responseWallets []models.WalletStruct
	for _, wallet := range wallets {
		var balances []models.Balance
		for _, balance := range wallet.Balance {
			balances = append(balances, models.Balance{
				Currency:   balance.Currency,
				Amount:     balance.Amount,
				LastUpdate: balance.LastUpdate,
			})
		}
		responseWallets = append(responseWallets, models.WalletStruct{
			Id:       wallet.Id,
			UserId:   wallet.UserId,
			Balances: balances,
		})
	}

	ctx.JSON(200, gin.H{
		"wallets": responseWallets,
	})
}
