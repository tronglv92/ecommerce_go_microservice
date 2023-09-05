package userstore

import (
	"github.com/tronglv92/accounts/plugin/storage/sdkredis"
)

type userCached struct {
	cacheStore sdkredis.Cache
}

func NewUserCache(cacheStore sdkredis.Cache) *userCached {
	return &userCached{

		cacheStore: cacheStore,
	}
}
