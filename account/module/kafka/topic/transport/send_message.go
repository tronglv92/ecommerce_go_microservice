package sendmessagegin

import (
	"net/http"

	"github.com/tronglv92/accounts/common"
	"github.com/tronglv92/accounts/plugin/kafka"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	sendmessagebiz "github.com/tronglv92/accounts/module/kafka/topic/biz"
)

func SendMessage(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		ps := sc.MustGet(common.PluginKafka).(kafka.Pubsub)
		biz := sendmessagebiz.NewSendMessageBiz(ps)

		err := biz.SendMessage(ctx.Request.Context())
		if err != nil {
			panic(err)

		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse("success"))
	}
}
