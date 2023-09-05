package cardrest

import (
	"net/http"
	"strconv"

	loanbiz "github.com/tronglv92/loans/module/loan/biz"
	loanrepo "github.com/tronglv92/loans/module/loan/repository"
	loanstorage "github.com/tronglv92/loans/module/loan/storage/gorm"
	"github.com/tronglv92/loans/plugin/storage/sdkgorm"

	goservice "github.com/tronglv92/ecommerce_go_common"
	"github.com/tronglv92/loans/common"

	"github.com/gin-gonic/gin"
)

func ListLoanByCustomerId(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customerId := ctx.Param("id")
		id, err := strconv.Atoi(customerId)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := sc.MustGet(common.DBMain).(sdkgorm.GormInterface)
		db.WithContext(ctx.Request.Context())
		dbSession := db.Session()
		store := loanstorage.NewSQLStore(dbSession)

		repo := loanrepo.NewListLoanByCustomerIdRepo(store)
		biz := loanbiz.NewListLoanByCustomerIdBiz(repo)
		result, err := biz.ListLoanByCustomerId(ctx.Request.Context(), id)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask()
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
