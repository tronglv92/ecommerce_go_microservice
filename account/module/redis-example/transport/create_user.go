package usergin

import (
	"fmt"
	"net/http"

	"github.com/tronglv92/accounts/common"
	userbiz "github.com/tronglv92/accounts/module/redis-example/biz"
	usermodel "github.com/tronglv92/accounts/module/redis-example/model"
	userstore "github.com/tronglv92/accounts/module/redis-example/storage"
	"github.com/tronglv92/accounts/plugin/storage/sdkredis"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateUser(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data usermodel.User
		if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
			fmt.Printf("CreateComment err %v \n", err)
			panic(err)
		}

		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		storage := userstore.NewUserCache(sdkredis.NewRedisCache(sc))
		biz := userbiz.NewCreateUserBiz(storage)

		result, err := biz.CreateUser(ctx.Request.Context(), &data)
		if err != nil {
			panic(err)

		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
