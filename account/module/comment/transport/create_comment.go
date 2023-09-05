package commentgin

import (
	"fmt"
	"net/http"

	"github.com/tronglv92/accounts/common"
	commentbiz "github.com/tronglv92/accounts/module/comment/biz"
	commentmodel "github.com/tronglv92/accounts/module/comment/model"
	commentstore "github.com/tronglv92/accounts/module/comment/storage/mongo/comment"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateComment(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data commentmodel.CommentCreate
		if err := ctx.ShouldBindBodyWith(&data, binding.JSON); err != nil {
			fmt.Printf("CreateComment err %v \n", err)
			panic(err)
		}

		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		db := sc.MustGet(common.DBMongo).(*mongo.Client)

		storage := commentstore.NewMongoStore(db)
		biz := commentbiz.NewCreateCommentBiz(storage)

		result, err := biz.CreateComment(ctx.Request.Context(), &data)
		if err != nil {
			panic(err)

		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
