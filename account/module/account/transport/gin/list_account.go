package accountgin

import (
	"net/http"

	"github.com/tronglv92/accounts/common"
	restaurantbiz "github.com/tronglv92/accounts/module/account/biz"
	restaurantmodel "github.com/tronglv92/accounts/module/account/model"
	restaurantrepo "github.com/tronglv92/accounts/module/account/repository"

	accountstorage "github.com/tronglv92/accounts/module/account/storage/gorm"
	customerstorage "github.com/tronglv92/accounts/module/customer/storage/gorm"
	sdkgorm "github.com/tronglv92/accounts/plugin/storage/sdkgorm"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func ListAccount(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var pagingData common.Paging
		if err := c.ShouldBind(&pagingData); err != nil {
			panic(common.ErrInvalidRequest(err))

		}

		pagingData.Fulfill()

		var filter restaurantmodel.Filter
		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := sc.MustGet(common.DBMain).(sdkgorm.GormInterface)
		db.WithContext(ctx)
		dbSession := db.Session()
		accountSqlStore := accountstorage.NewSQLStore(dbSession)
		customerSqlStore := customerstorage.NewSQLStore(dbSession)

		repo := restaurantrepo.NewListAccountRepo(accountSqlStore, customerSqlStore)
		biz := restaurantbiz.NewListAccountBiz(repo)
		result, err := biz.ListAccount(ctx, &filter, &pagingData)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))

	}
}
