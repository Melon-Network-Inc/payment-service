package dbcontext

import (
	"context"
	"log"

	"github.com/Melon-Network-Inc/entity-repo/pkg/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB represents a DB connection that can be used to run SQL queries.
type DB struct {
	db *gorm.DB
}

// New returns a new DB connection that wraps the given dbx.DB instance.
func NewDatabase(db *gorm.DB) *DB {
	return &DB{db}
}

// DB returns the dbx.DB wrapped by this object.
func (db *DB) DB() *gorm.DB {
	return db.db
}

// With returns a Builder that can be used to build and execute SQL queries.
// With will return the transaction if it is found in the given context.
// Otherwise it will return a DB connection associated with the context.
func (db *DB) With(c context.Context) *gorm.DB {
	return db.db.WithContext(c)
}

// DB returns the gorm.DB wrapped by this object.
func ConnectToDatabase(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&entity.Transaction{})

	return db, err
}
