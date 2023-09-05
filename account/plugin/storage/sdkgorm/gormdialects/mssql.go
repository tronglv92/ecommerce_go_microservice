package gormdialects

import (
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// MSSqlDB Get MS SQL DB connection
// URI string
// Ex: sqlserver://username:password@localhost:1433?database=dbname
func MSSqlDB(uri string, uriReadOnly []string) (db *gorm.DB, err error) {

	db, err = gorm.Open(sqlserver.Open(uri), &gorm.Config{})

	var replicas []gorm.Dialector
	for _, uriRead := range uriReadOnly {
		item := sqlserver.Open(uriRead)
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
