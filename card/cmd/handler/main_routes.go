package handlers

import (
	cardgin "github.com/tronglv92/cards/module/card/transport/gin"

	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func MainRoute(router *gin.Engine, sc goservice.ServiceContext) {
	v1 := router.Group("/v1")
	{
		v1.GET("/cards", cardgin.ListCard(sc))
	}
}
