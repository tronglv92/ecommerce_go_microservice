package sendmessagegin

import (
	"net/http"

	"github.com/tronglv92/accounts/common"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	sendmessagebiz "github.com/tronglv92/accounts/module/rabbitmq/dle-retry/biz"
	apprabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"
)

func SendMessageDLX(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		ps := sc.MustGet(common.PluginRabbitMQ).(apprabbitmq.Pubsub)
		biz := sendmessagebiz.NewSendMessageDLXBiz(ps)

		err := biz.SendMessageDLX(ctx.Request.Context())
		if err != nil {
			panic(err)

		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse("success"))
	}
}
