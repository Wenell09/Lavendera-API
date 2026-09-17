package app

import (
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	return database.DBConnection()
}
