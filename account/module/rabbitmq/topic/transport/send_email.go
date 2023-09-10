package sendemailgin

import (
	"fmt"
	"net/http"

	"github.com/tronglv92/accounts/common"
	apprabbitmq "github.com/tronglv92/accounts/plugin/rabbitmq"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	sendEmailbiz "github.com/tronglv92/accounts/module/rabbitmq/topic/biz"
	model "github.com/tronglv92/accounts/module/rabbitmq/topic/model"
)

func SendEmail(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data model.MessageMail
		if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
			fmt.Printf("MessageMail err %v \n", err)
			panic(err)
		}
		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		ps := sc.MustGet(common.PluginRabbitMQ).(apprabbitmq.Pubsub)
		biz := sendEmailbiz.NewSendEmailBiz(ps)

		err := biz.SendEmail(ctx.Request.Context(), &data)
		if err != nil {
			panic(err)

		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse("success"))
	}
}
