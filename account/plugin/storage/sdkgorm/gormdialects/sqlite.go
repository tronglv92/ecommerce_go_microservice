package gormdialects

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// SQLiteDB Get SQLite DB connection
// URI string
// Ex: /tmp/gorm.db
func SQLiteDB(uri string, uriReadOnly []string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(uri))
	var replicas []gorm.Dialector
	for _, uriRead := range uriReadOnly {
		item := sqlite.Open(uriRead)
		replicas = append(replicas, item)
	}

	db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200))
	return db, err
}
