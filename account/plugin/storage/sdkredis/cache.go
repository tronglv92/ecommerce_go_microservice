package sdkredis

import (
	"context"
	"time"

	"github.com/tronglv92/accounts/common"
	goservice "github.com/tronglv92/ecommerce_go_common"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
}
type redisCache struct {
	store  *cache.Cache
	client *redis.ClusterClient
}

func NewRedisCache(sc goservice.ServiceContext) *redisCache {
	rdClient := sc.MustGet(common.PluginRedis).(*redis.ClusterClient)

	c := cache.New(&cache.Options{
		Redis:      rdClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	return &redisCache{store: c, client: rdClient}
}
func (rdc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {

	return rdc.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}
func (rdc *redisCache) Get(ctx context.Context, key string, value interface{}) error {

	return rdc.store.Get(ctx, key, value)
}

func (rdc *redisCache) Delete(ctx context.Context, key string) error {

	return rdc.store.Delete(ctx, key)
}
func (rdc *redisCache) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {

	return rdc.client.Scan(ctx, cursor, match, count)
}
