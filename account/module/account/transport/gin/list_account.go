package accountgin

import (
	"net/http"

	"github.com/tronglv92/accounts/common"
	restaurantbiz "github.com/tronglv92/accounts/module/account/biz"
	restaurantmodel "github.com/tronglv92/accounts/module/account/model"
	restaurantrepo "github.com/tronglv92/accounts/module/account/repository"

	accountstorage "github.com/tronglv92/accounts/module/account/storage/gorm"
	customerstorage "github.com/tronglv92/accounts/module/customer/storage/gorm"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func ListAccount(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pagingData common.Paging
		if err := ctx.ShouldBind(&pagingData); err != nil {
			panic(common.ErrInvalidRequest(err))

		}

		pagingData.Fulfill()

		var filter restaurantmodel.Filter
		if err := ctx.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		db := sc.MustGet(common.DBMain).(*gorm.DB)

		accountSqlStore := accountstorage.NewSQLStore(db)
		customerSqlStore := customerstorage.NewSQLStore(db)

		repo := restaurantrepo.NewListAccountRepo(accountSqlStore, customerSqlStore)
		biz := restaurantbiz.NewListAccountBiz(repo)
		result, err := biz.ListAccount(ctx.Request.Context(), &filter, &pagingData)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}
		ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))

	}
}
