package gormdialects

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// PostgresDB Get Postgres DB connection
// URI string
// Ex: host=myhost port=myport user=gorm dbname=gorm password=mypassword
func PostgresDB(uri string, uriReadOnly []string) (db *gorm.DB, err error) {
	fmt.Printf("PostgresDB uri: %v uriReadOnly: %v", uri, uriReadOnly)
	db, err = gorm.Open(postgres.Open(uri))

	var replicas []gorm.Dialector
	for _, uriRead := range uriReadOnly {
		item := postgres.Open(uriRead)
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
