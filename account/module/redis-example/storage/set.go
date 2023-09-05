package userstore

import (
	"context"
	"fmt"

	"github.com/tronglv92/accounts/common"
	usermodel "github.com/tronglv92/accounts/module/redis-example/model"

	"time"
)

func (c *userCached) SetUser(ctx context.Context,
	data *usermodel.User) error {

	key := fmt.Sprintf(common.CacheKey, data.ID)

	err := c.cacheStore.Set(ctx, key, data, time.Hour*2)

	return err
}
