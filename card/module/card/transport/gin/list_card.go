package accountgin

import (
	"net/http"

	"github.com/tronglv92/cards/common"
	cardbiz "github.com/tronglv92/cards/module/card/biz"
	cardmodel "github.com/tronglv92/cards/module/card/model"
	cardrepo "github.com/tronglv92/cards/module/card/repository"
	"github.com/tronglv92/cards/plugin/storage/sdkgorm"

	accountstorage "github.com/tronglv92/cards/module/card/storage/gorm"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func ListCard(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var pagingData common.Paging
		if err := c.ShouldBind(&pagingData); err != nil {
			panic(common.ErrInvalidRequest(err))

		}

		pagingData.Fulfill()

		var filter cardmodel.Filter
		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		db := sc.MustGet(common.DBMain).(sdkgorm.GormInterface)
		db.WithContext(ctx)
		dbSession := db.Session()
		accountSqlStore := accountstorage.NewSQLStore(dbSession)
		repo := cardrepo.NewListCardRepo(accountSqlStore)
		biz := cardbiz.NewListCardBiz(repo)
		result, err := biz.ListCard(ctx, &filter, &pagingData)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask()
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))

	}
}
