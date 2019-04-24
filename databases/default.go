package databases

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DefaultModels = []interface{}{
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
	models.Views{},
	models.ContentType{},
}

func ConnectDefaultDB() *gorm.DB {
	dbConf := config.Config.Databases.Default
	params := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		dbConf.Host,
		dbConf.Port,
		dbConf.User,
		dbConf.Database,
		dbConf.Password,
	)
	db, err := gorm.Open("postgres", params)
	logger.HandleError(err)
	if config.Config.System.Debug {
		db.LogMode(true)
	}

	return db
}

func FindModelByContentType(db *gorm.DB, contentType string) interface{} {
	for _, model := range DefaultModels {
		tableName := db.NewScope(model).GetModelStruct().TableName(db)
		if tableName == contentType {
			return model
		}
	}
	return nil
}

func MigrateDefaultDB() {
	db := ConnectDefaultDB()

	db.AutoMigrate(DefaultModels...)

	type UsersWorkGroups struct {
		UserID      uint
		WorkGroupID uint
	}

	type MonitoringShopsSegments struct {
		MonitoringShopID uint
		SegmentID        uint
	}

	type MonitoringTypesPeriods struct {
		MonitoringTypeID uint
		PeriodID         uint
	}

	type RolesPermissions struct {
		RoleID       uint
		PermissionID uint
	}

	type WorkGroupsMonitoringShops struct {
		WorkGroupID      uint
		MonitoringShopID uint
	}

	type WorkGroupsRegions struct {
		WorkGroupID uint
		RegionsID   uint
	}

	type UsersRoles struct {
		UserID uint
		RoleID uint
	}

	db.Model(&models.User{}).AddForeignKey("token_id", "tokens(id)", "CASCADE", "CASCADE")
	db.Model(&models.Ware{}).AddForeignKey("segment_id", "segments(id)", "CASCADE", "CASCADE")
	db.Model(&models.CompletedWare{}).AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "RESTRICT", "CASCADE")
	db.Model(&models.CompletedWare{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE")
	db.Model(&models.CompletedWare{}).AddForeignKey("monitoring_type_id", "monitoring_types(id)", "RESTRICT", "CASCADE")
	db.Model(&models.CompletedWare{}).AddForeignKey("region_id", "regions(id)", "RESTRICT", "CASCADE")
	db.Model(&models.CompletedWare{}).AddForeignKey("ware_id", "wares(id)", "RESTRICT", "CASCADE")
	db.Model(&models.Photos{}).AddForeignKey("completed_ware_id", "completed_wares(id)", "CASCADE", "CASCADE")
	db.Model(&models.Views{}).AddForeignKey("parent_id", "views(id)", "RESTRICT", "CASCADE")
	db.Model(&models.Views{}).AddForeignKey("content_type_id", "content_types(id)", "RESTRICT", "CASCADE")
	db.Model(&models.Permission{}).AddForeignKey("view_id", "views(id)", "CASCADE", "CASCADE")

	db.Model(&UsersWorkGroups{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&UsersWorkGroups{}).AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")

	db.Model(&MonitoringShopsSegments{}).AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "CASCADE", "CASCADE")
	db.Model(&MonitoringShopsSegments{}).AddForeignKey("segment_id", "segments(id)", "CASCADE", "CASCADE")

	db.Model(&MonitoringTypesPeriods{}).AddForeignKey("monitoring_type_id", "monitoring_types(id)", "CASCADE", "CASCADE")
	db.Model(&MonitoringTypesPeriods{}).AddForeignKey("period_id", "periods(id)", "CASCADE", "CASCADE")

	db.Model(&RolesPermissions{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")
	db.Model(&RolesPermissions{}).AddForeignKey("permission_id", "permissions(id)", "CASCADE", "CASCADE")

	db.Model(&WorkGroupsMonitoringShops{}).AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")
	db.Model(&WorkGroupsMonitoringShops{}).AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "CASCADE", "CASCADE")

	db.Model(&WorkGroupsRegions{}).AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")
	db.Model(&WorkGroupsRegions{}).AddForeignKey("regions_id", "regions(id)", "CASCADE", "CASCADE")

	db.Model(&UsersRoles{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&UsersRoles{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

	for _, model := range DefaultModels {
		tableName := db.NewScope(model).GetModelStruct().TableName(db)
		db.FirstOrCreate(&models.ContentType{}, models.ContentType{
			Table: tableName,
		})
	}

	logger.HandleError(db.Close())
}
