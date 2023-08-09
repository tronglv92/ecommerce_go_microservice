package handlers

import (
	goservice "github.com/tronglv92/ecommerce_go_common"
	loanrest "github.com/tronglv92/loans/module/loan/transport/rest"

	"github.com/gin-gonic/gin"
)

func InternalRoute(router *gin.Engine, sc goservice.ServiceContext) {
	internal := router.Group("/internal")
	{
		internal.GET("/loans/:id", loanrest.ListLoanByCustomerId(sc))
	}
}
