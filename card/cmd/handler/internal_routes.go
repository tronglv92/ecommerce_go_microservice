package handlers

import (
	cardrest "github.com/tronglv92/cards/module/card/transport/rest"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func InternalRoute(router *gin.Engine, sc goservice.ServiceContext) {
	internal := router.Group("/internal")
	{
		internal.GET("/cards/:id", cardrest.ListCardByCustomerID(sc))
	}
}
