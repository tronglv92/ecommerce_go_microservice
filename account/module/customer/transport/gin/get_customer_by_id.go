package customergin

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/tronglv92/accounts/common"
	customerbiz "github.com/tronglv92/accounts/module/customer/biz"
	customerrepo "github.com/tronglv92/accounts/module/customer/repo"
	"github.com/tronglv92/ecommerce_go_common/logger"

	customerstorage "github.com/tronglv92/accounts/module/customer/storage/gorm"

	// cardrestful "github.com/tronglv92/accounts/module/customer/storage/client_remotecall/resful/card"

	cardgrpcstore "github.com/tronglv92/accounts/module/customer/storage/client_remotecall/grpc/card"
	loanrestful "github.com/tronglv92/accounts/module/customer/storage/client_remotecall/resful/loan"
	cardgrpc "github.com/tronglv92/accounts/proto/card"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func GetCustomerByID(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logger.GetCurrent().GetLogger("customer.transport.get_customer_by_id")
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		db := sc.MustGet(common.DBMain).(*gorm.DB)
		customerSqlStore := customerstorage.NewSQLStore(db)

		// restfulCardUri := viper.GetString(common.RestfulCardUri)
		// logger.Debugf("restful-card-uri=%v", restfulCardUri)
		// cardRestApi := cardrestful.NewCardRestfulStore(clientRest, restfulCardUri)
		cardClientGrpc := sc.MustGet(common.PluginGrpcCardClient).(cardgrpc.CardServiceClient)
		cardGrpcStore := cardgrpcstore.NewCardRestfulStore(cardClientGrpc)

		clientRest := sc.MustGet(common.PluginRestService).(*resty.Client)
		restfulLoanUri := viper.GetString(common.RestfulLoanUri)
		logger.Debugf("restful-loan-uri=%v", restfulLoanUri)
		loanRestApi := loanrestful.NewLoanRestfulStore(clientRest, restfulLoanUri)

		repo := customerrepo.NewCustomerByIdRepo(customerSqlStore, cardGrpcStore, loanRestApi)
		biz := customerbiz.NewGetCustomerBiz(repo)
		customer, err := biz.GetCustomerById(ctx.Request.Context(), int(uid.GetLocalID()))

		if err != nil {
			panic(err)
		}
		customer.Mask()

		for i, _ := range customer.Accounts {

			customer.Accounts[i].Mask()
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(customer))

	}
}
