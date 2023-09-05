package grpcstorage

import (
	"github.com/tronglv92/cards/plugin/storage/sdkgorm"
	"gorm.io/gorm"
)

type sqlStore struct {
	gorm      sdkgorm.GormInterface
	dbSession *gorm.DB
}

func NewSQLStore(gorm sdkgorm.GormInterface, dbSession *gorm.DB) *sqlStore {
	return &sqlStore{gorm: gorm, dbSession: dbSession}
}
