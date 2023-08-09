package accountgin

import (
	"net/http"

	"github.com/tronglv92/cards/common"
	cardbiz "github.com/tronglv92/cards/module/card/biz"
	cardmodel "github.com/tronglv92/cards/module/card/model"
	cardrepo "github.com/tronglv92/cards/module/card/repository"

	accountstorage "github.com/tronglv92/cards/module/card/storage/gorm"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func ListCard(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pagingData common.Paging
		if err := ctx.ShouldBind(&pagingData); err != nil {
			panic(common.ErrInvalidRequest(err))

		}

		pagingData.Fulfill()

		var filter cardmodel.Filter
		if err := ctx.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		db := sc.MustGet(common.DBMain).(*gorm.DB)
		accountSqlStore := accountstorage.NewSQLStore(db)
		repo := cardrepo.NewListCardRepo(accountSqlStore)
		biz := cardbiz.NewListCardBiz(repo)
		result, err := biz.ListCard(ctx.Request.Context(), &filter, &pagingData)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask()
		}
		ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))

	}
}
