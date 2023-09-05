package usergin

import (
	"errors"
	"net/http"

	"github.com/tronglv92/accounts/common"
	userbiz "github.com/tronglv92/accounts/module/redis-example/biz"
	userstore "github.com/tronglv92/accounts/module/redis-example/storage"
	"github.com/tronglv92/accounts/plugin/storage/sdkredis"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
)

func GetUser(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if len(id) == 0 {
			panic(common.ErrInvalidRequest(errors.New("Missing Id")))
		}
		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		storage := userstore.NewUserCache(sdkredis.NewRedisCache(sc))
		biz := userbiz.NewGetUserBiz(storage)

		result, err := biz.GetUser(c.Request.Context(), id)
		if err != nil {
			panic(err)

		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
