package migrations

import (
	"github.com/NurymGM/hotell/initializers"
	"github.com/NurymGM/hotell/models"
)

func Migrate() {
	initializers.DB.AutoMigrate(&models.Room{})
}