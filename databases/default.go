package databases

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/jinzhu/gorm"
)

func ConnectDefaultDB() *gorm.DB {
	defaultDb := config.Config.Databases.Default
	params := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		defaultDb.Host,
		defaultDb.Port,
		defaultDb.User,
		defaultDb.Database,
		defaultDb.Password,
	)
	db, err := gorm.Open("postgres", params)
	helpers.HandlerError(err)
	if config.Config.System.Debug {
		db.LogMode(true)
	}

	return db
}

func MigrateDefaultDB() {
	db := ConnectDefaultDB()
	db.AutoMigrate(
		models.Role{},
		models.Permission{},
		models.User{},
		models.Segment{},
		models.Ware{},
		models.WorkGroup{},
		models.Regions{},
		models.MonitoringShop{},
		models.Token{},
		models.Period{},
		models.CompletedWare{},
		models.MonitoringType{},
		models.Photos{},
	)
	helpers.HandlerError(db.Close())
}
