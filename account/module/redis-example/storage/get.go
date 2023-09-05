package userstore

import (
	"context"
	"fmt"

	"github.com/tronglv92/accounts/common"
	usermodel "github.com/tronglv92/accounts/module/redis-example/model"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

func (c *userCached) GetUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error) {
	logger := logger.GetCurrent().GetLogger("GetUser")
	var user usermodel.User
	userId := conditions["id"]
	key := fmt.Sprintf(common.CacheKey, userId)

	logger.Debugf("key %v", key)
	err := c.cacheStore.Get(ctx, key, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
