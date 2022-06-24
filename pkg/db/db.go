package db

import (
	"log"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dbURL := "postgres://postgres:123456@localhost:5432/melon_service"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&entity.Transaction{})

	return db
}
