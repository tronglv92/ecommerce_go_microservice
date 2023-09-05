package userbiz

import (
	"context"

	"github.com/tronglv92/accounts/common"
	usermodel "github.com/tronglv92/accounts/module/redis-example/model"
)

type CreateUserStorage interface {
	SetUser(ctx context.Context,
		data *usermodel.User) error
}

type createUserBiz struct {
	createUserStorage CreateUserStorage
}

func NewCreateUserBiz(createUserStorage CreateUserStorage) *createUserBiz {
	return &createUserBiz{
		createUserStorage: createUserStorage,
	}
}
func (business *createUserBiz) CreateUser(ctx context.Context, data *usermodel.User) (*usermodel.User, error) {

	data.Fullfill()

	//4. create new comment
	err := business.createUserStorage.SetUser(ctx, data)

	if err != nil {
		return nil, common.ErrCannotCreateEntity(usermodel.EntityName, err)
	}

	return data, nil
}
