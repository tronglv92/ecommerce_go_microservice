package sendnotificationgin

import (
	"net/http"

	"github.com/tronglv92/accounts/common"
	apprabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	sendNotifibiz "github.com/tronglv92/accounts/module/rabbitmq/fanout/biz"
)

func SendNotification(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		ps := sc.MustGet(common.PluginRabbitMQ).(apprabbitmq.Pubsub)
		biz := sendNotifibiz.NewSendNotificationBiz(ps)

		err := biz.SendNotification(ctx.Request.Context())
		if err != nil {
			panic(err)

		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse("success"))
	}
}
