package customergin

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/tronglv92/accounts/common"

	customerbiz "github.com/tronglv92/accounts/module/customer/biz"
	customerrepo "github.com/tronglv92/accounts/module/customer/repo"
	cardrestful "github.com/tronglv92/accounts/module/customer/storage/client_remotecall/grpc/card"
	loanrestful "github.com/tronglv92/accounts/module/customer/storage/client_remotecall/resful/loan"
	customerstore "github.com/tronglv92/accounts/module/customer/storage/gorm"
	sdkgorm "github.com/tronglv92/accounts/plugin/storage/sdkgorm"
	cardgrpc "github.com/tronglv92/accounts/proto/card"
	"github.com/tronglv92/ecommerce_go_common/logger"

	// cardrestful "github.com/tronglv92/accounts/module/customer/storage/client_remotecall/resful/card"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func GetCustomerByID(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.GetCurrent().GetLogger("customer.transport.get_customer_by_id")
		// tracer := sc.MustGet(common.PluginOpenTelemetry).(trace.Tracer)
		// ctx, span := tracer.Start(c.Request.Context(), "GetCustomerByID")
		// span.RecordError(errors.New("test"))

		// status := http.StatusBadRequest
		// // attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		// spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)

		// span.SetAttributes(attrs...)
		// span.SetStatus(spanStatus, spanMessage)
		// defer span.End()

		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}

		db := sc.MustGet(common.DBMain).(sdkgorm.GormInterface)
		db.WithContext(c.Request.Context())
		dbSession := db.Session()
		customerSqlStore := customerstore.NewSQLStore(dbSession)

		// restfulCardUri := viper.GetString(common.RestfulCardUri)
		// logger.Debugf("restful-card-uri=%v", restfulCardUri)
		// cardRestApi := cardrestful.NewCardRestfulStore(clientRest, restfulCardUri)
		cardClientGrpc := sc.MustGet(common.PluginGrpcCardClient).(cardgrpc.CardServiceClient)
		cardGrpcStore := cardrestful.NewCardRestfulStore(cardClientGrpc)

		clientRest := sc.MustGet(common.PluginRestService).(*resty.Client)
		restfulLoanUri := viper.GetString(common.RestfulLoanUri)
		logger.Debugf("restful-loan-uri=%v", restfulLoanUri)

		loanRestApi := loanrestful.NewLoanRestfulStore(clientRest, restfulLoanUri)

		repo := customerrepo.NewCustomerByIdRepo(customerSqlStore, cardGrpcStore, loanRestApi)
		biz := customerbiz.NewGetCustomerBiz(repo)
		customer, err := biz.GetCustomerById(c, int(uid.GetLocalID()))

		if err != nil {
			panic(err)
		}
		customer.Mask()

		for i, _ := range customer.Accounts {

			customer.Accounts[i].Mask()
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(customer))

	}
}
