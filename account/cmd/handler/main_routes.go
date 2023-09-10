package handlers

import (
	"github.com/gin-gonic/gin"
	accountgin "github.com/tronglv92/accounts/module/account/transport/gin"
	commentgin "github.com/tronglv92/accounts/module/comment/transport"
	customergin "github.com/tronglv92/accounts/module/customer/transport/gin"
	kafkagin "github.com/tronglv92/accounts/module/kafka/topic/transport"
	pubsubgin "github.com/tronglv92/accounts/module/rabbitmq/direct/transport"
	messagedlx "github.com/tronglv92/accounts/module/rabbitmq/dle-retry/transport"
	fanoutgin "github.com/tronglv92/accounts/module/rabbitmq/fanout/transport"
	topicgin "github.com/tronglv92/accounts/module/rabbitmq/topic/transport"
	redisgin "github.com/tronglv92/accounts/module/redis-example/transport"
	goservice "github.com/tronglv92/ecommerce_go_common"
)

func MainRoute(router *gin.Engine, sc goservice.ServiceContext) {
	v1 := router.Group("/v1")
	{
		v1.GET("/accounts", accountgin.ListAccount(sc))
		v1.GET("/customer/:id", customergin.GetCustomerByID(sc))

		comments := v1.Group("/comments")
		{
			comments.POST("/create-comment", commentgin.CreateComment(sc))
			comments.GET("", commentgin.ListComments(sc))
		}
		redis := v1.Group("/redis")
		{
			redis.POST("/user", redisgin.CreateUser(sc))
			redis.GET("/user/:id", redisgin.GetUser(sc))
		}
		pubsub := v1.Group("/rabbitmq")
		{
			pubsub.GET("/direct/sendmessage", pubsubgin.SendMessage(sc))
			pubsub.GET("/fanout/sendnotification", fanoutgin.SendNotification(sc))
			pubsub.GET("/topic/sendemail", topicgin.SendEmail(sc))
			pubsub.GET("/direct/sendmessagedlx", messagedlx.SendMessageDLX(sc))
		}
		kafka := v1.Group("/kafka")
		{
			kafka.GET("/topic/sendmessage", kafkagin.SendMessage(sc))

		}

	}
}
