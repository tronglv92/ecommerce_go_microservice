package handlers

import (
	accountgin "github.com/tronglv92/accounts/module/account/transport/gin"
	customergin "github.com/tronglv92/accounts/module/customer/transport/gin"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func MainRoute(router *gin.Engine, sc goservice.ServiceContext) {
	v1 := router.Group("/v1")
	{
		v1.GET("/accounts", accountgin.ListAccount(sc))
		v1.GET("/customer/:id", customergin.GetCustomerByID(sc))
	}
}
