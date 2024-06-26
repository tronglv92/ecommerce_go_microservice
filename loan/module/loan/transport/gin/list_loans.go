package accountgin

import (
	"net/http"

	"github.com/tronglv92/loans/common"
	loanbiz "github.com/tronglv92/loans/module/loan/biz"
	loanmodel "github.com/tronglv92/loans/module/loan/model"
	loanrepo "github.com/tronglv92/loans/module/loan/repository"
	"github.com/tronglv92/loans/plugin/storage/sdkgorm"

	accountstorage "github.com/tronglv92/loans/module/loan/storage/gorm"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func ListLoan(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pagingData common.Paging
		if err := ctx.ShouldBind(&pagingData); err != nil {
			panic(common.ErrInvalidRequest(err))

		}

		pagingData.Fulfill()

		var filter loanmodel.Filter
		if err := ctx.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		db := sc.MustGet(common.DBMain).(sdkgorm.GormInterface)
		db.WithContext(ctx.Request.Context())
		dbSession := db.Session()

		loanSqlStore := accountstorage.NewSQLStore(dbSession)
		repo := loanrepo.NewListLoanRepo(loanSqlStore)
		biz := loanbiz.NewListLoanBiz(repo)
		result, err := biz.ListLoan(ctx.Request.Context(), &filter, &pagingData)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask()
		}
		ctx.JSON(http.StatusOK, common.NewSuccessResponse(result, pagingData, filter))

	}
}
