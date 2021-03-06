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
	models.Monitoring{},
	models.Region{},
	models.MonitoringShop{},
	models.WorkGroup{},
	models.Token{},
	models.Period{},
	models.CompletedWare{},
	models.MonitoringType{},
	models.Photos{},
	models.Views{},
	models.ContentType{},
}

var DB *gorm.DB

func ConnectDefaultDB() *gorm.DB {
	if DB == nil {
		var err error
		dbConf := config.Config.Databases.Default
		params := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			dbConf.Host,
			dbConf.Port,
			dbConf.User,
			dbConf.Database,
			dbConf.Password,
		)
		DB, err = gorm.Open("postgres", params)
		logger.HandleError(err)
		if config.Config.System.Debug {
			DB.LogMode(true)
		}
	}
	return DB
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
	db.Model(models.User{}).AddForeignKey("token_id", "tokens(id)", "CASCADE", "CASCADE")
	db.Model(models.Ware{}).AddForeignKey("segment_id", "segments(id)", "CASCADE", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "RESTRICT", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("monitoring_type_id", "monitoring_types(id)", "RESTRICT", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("ware_id", "wares(id)", "RESTRICT", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("region_id", "regions(id)", "RESTRICT", "CASCADE")
	db.Model(models.Photos{}).AddForeignKey("completed_ware_id", "completed_wares(id)", "CASCADE", "CASCADE")
	db.Model(models.Views{}).AddForeignKey("parent_id", "views(id)", "RESTRICT", "CASCADE")
	db.Model(models.Views{}).AddForeignKey("content_type_id", "content_types(id)", "RESTRICT", "CASCADE")
	db.Model(models.Permission{}).AddForeignKey("view_id", "views(id)", "CASCADE", "CASCADE")

	db.Table("work_groups_users").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Table("work_groups_users").AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")

	db.Table("monitoring_types_periods").AddForeignKey("monitoring_type_id", "monitoring_types(id)", "CASCADE", "CASCADE")
	db.Table("monitoring_types_periods").AddForeignKey("period_id", "periods(id)", "CASCADE", "CASCADE")

	db.Table("roles_permissions").AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")
	db.Table("roles_permissions").AddForeignKey("permission_id", "permissions(id)", "CASCADE", "CASCADE")

	db.Table("work_groups_monitoring_shops").AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")
	db.Table("work_groups_monitoring_shops").AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "CASCADE", "CASCADE")

	db.Table("regions_monitorings").AddForeignKey("monitoring_id", "monitorings(id)", "CASCADE", "CASCADE")
	db.Table("regions_monitorings").AddForeignKey("region_id", "regions(id)", "CASCADE", "CASCADE")

	db.Table("monitorings_work_groups").AddForeignKey("monitoring_id", "monitorings(id)", "CASCADE", "CASCADE")
	db.Table("monitorings_work_groups").AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")

	db.Table("users_roles").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Table("users_roles").AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

	db.Table("monitorings_wares").AddForeignKey("ware_id", "wares(id)", "CASCADE", "CASCADE")
	db.Table("monitorings_wares").AddForeignKey("monitoring_id", "monitorings(id)", "CASCADE", "CASCADE")

	db.Table("monitoring_shops_segments").AddForeignKey("segment_id", "segments(id)", "CASCADE", "CASCADE")
	db.Table("monitoring_shops_segments").AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "CASCADE", "CASCADE")

	tx := db.Begin()
	for _, model := range DefaultModels {
		tableName := tx.NewScope(model).GetModelStruct().TableName(tx)
		tx.FirstOrCreate(new(models.ContentType), models.ContentType{
			Table: tableName,
		})
	}
	tx.Commit()
}
