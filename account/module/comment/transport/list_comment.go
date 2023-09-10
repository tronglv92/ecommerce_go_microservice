package commentgin

import (
	commentbiz "github.com/tronglv92/accounts/module/comment/biz"
	commentstore "github.com/tronglv92/accounts/module/comment/storage/mongo/comment"
	goservice "github.com/tronglv92/ecommerce_go_common"
	"go.mongodb.org/mongo-driver/mongo"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tronglv92/accounts/common"
)

func ListComments(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var pagingData common.Paging
		if err := ctx.ShouldBind(&pagingData); err != nil {

			panic(common.ErrInvalidRequest(err))
		}
		pagingData.Fulfill()
		// requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		//id, err := strconv.Atoi(ctx.Param("id"))

		db := sc.MustGet(common.DBMongo).(*mongo.Client)
		storage := commentstore.NewMongoStore(db)
		biz := commentbiz.NewListCommentBiz(storage)

		results, err := biz.ListComments(ctx.Request.Context(), &pagingData)
		if err != nil {
			panic(err)

		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(results))
	}
}
