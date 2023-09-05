package cardtstorage

import (
	"gorm.io/gorm"
)

type sqlStore struct {
	dbSession *gorm.DB
}

func NewSQLStore(dbSession *gorm.DB) *sqlStore {
	return &sqlStore{dbSession: dbSession}
}
