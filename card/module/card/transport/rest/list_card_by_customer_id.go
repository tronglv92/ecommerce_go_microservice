package cardrest

import (
	"net/http"
	"strconv"

	"github.com/tronglv92/cards/common"
	cardbiz "github.com/tronglv92/cards/module/card/biz"
	cardrepo "github.com/tronglv92/cards/module/card/repository"
	cardstorage "github.com/tronglv92/cards/module/card/storage/gorm"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListCardByCustomerID(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customerId := ctx.Param("id")
		id, err := strconv.Atoi(customerId)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := sc.MustGet(common.DBMain).(*gorm.DB)
		store := cardstorage.NewSQLStore(db)
		repo := cardrepo.NewListCardByCustomerIdRepo(store)
		biz := cardbiz.NewListCardByCustomerIdBiz(repo)
		result, err := biz.ListCardByCustomerId(ctx.Request.Context(), id)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask()
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
