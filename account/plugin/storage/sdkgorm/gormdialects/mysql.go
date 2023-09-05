package gormdialects

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// MySqlDB Get MySQL DB connection
// URI string
// Ex: user:password@/db_name?charset=utf8&parseTime=True&loc=Local
func MySqlDB(uri string, uriReadOnly []string) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(uri), &gorm.Config{})

	var replicas []gorm.Dialector
	for _, uriRead := range uriReadOnly {
		item := mysql.Open(uriRead)
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
