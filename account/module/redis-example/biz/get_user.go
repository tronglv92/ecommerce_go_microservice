package userbiz

import (
	"context"

	"github.com/tronglv92/accounts/common"
	usermodel "github.com/tronglv92/accounts/module/redis-example/model"
)

type GetUserStorage interface {
	GetUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type getUserBiz struct {
	getUserStorage GetUserStorage
}

func NewGetUserBiz(getUserStorage GetUserStorage) *getUserBiz {
	return &getUserBiz{
		getUserStorage: getUserStorage,
	}
}
func (business *getUserBiz) GetUser(ctx context.Context, id string) (*usermodel.User, error) {

	//4. create new comment
	user, err := business.getUserStorage.GetUser(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return nil, common.ErrCannotCreateEntity(usermodel.EntityName, err)
	}

	return user, nil
}
