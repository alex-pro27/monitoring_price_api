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
	models.WorkGroup{},
	models.MonitoringGroups{},
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

	db.Model(models.User{}).AddForeignKey("token_id", "tokens(id)", "CASCADE", "CASCADE")
	db.Model(models.Ware{}).AddForeignKey("segment_id", "segments(id)", "CASCADE", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "RESTRICT", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("monitoring_type_id", "monitoring_types(id)", "RESTRICT", "CASCADE")
	db.Model(models.CompletedWare{}).AddForeignKey("ware_id", "wares(id)", "RESTRICT", "CASCADE")
	db.Model(models.Photos{}).AddForeignKey("completed_ware_id", "completed_wares(id)", "CASCADE", "CASCADE")
	db.Model(models.Views{}).AddForeignKey("parent_id", "views(id)", "RESTRICT", "CASCADE")
	db.Model(models.Views{}).AddForeignKey("content_type_id", "content_types(id)", "RESTRICT", "CASCADE")
	db.Model(models.Permission{}).AddForeignKey("view_id", "views(id)", "CASCADE", "CASCADE")
	db.Model(models.Monitoring{}).AddForeignKey("monitoring_type_id", "monitoring_types(id)", "RESTRICT", "CASCADE")
	db.Model(models.Monitoring{}).AddForeignKey("monitoring_group_id", "monitoring_groups(id)", "RESTRICT", "CASCADE")

	db.Table("users_work_groups").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Table("users_work_groups").AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")

	db.Table("monitoring_types_periods").AddForeignKey("monitoring_type_id", "monitoring_types(id)", "CASCADE", "CASCADE")
	db.Table("monitoring_types_periods").AddForeignKey("period_id", "periods(id)", "CASCADE", "CASCADE")

	db.Table("roles_permissions").AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")
	db.Table("roles_permissions").AddForeignKey("permission_id", "permissions(id)", "CASCADE", "CASCADE")

	db.Table("work_groups_monitorings").AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")
	db.Table("work_groups_monitorings").AddForeignKey("monitoring_id", "monitorings(id)", "CASCADE", "CASCADE")

	db.Table("work_groups_monitoring_shops").AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")
	db.Table("work_groups_monitoring_shops").AddForeignKey("monitoring_shop_id", "monitoring_shops(id)", "CASCADE", "CASCADE")

	db.Table("work_groups_monitoring_groups").AddForeignKey("work_group_id", "work_groups(id)", "CASCADE", "CASCADE")
	db.Table("work_groups_monitoring_groups").AddForeignKey("monitoring_groups_id", "monitoring_groups(id)", "CASCADE", "CASCADE")

	db.Table("users_roles").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Table("users_roles").AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

	db.Table("monitorings_wares").AddForeignKey("ware_id", "wares(id)", "CASCADE", "CASCADE")
	db.Table("monitorings_wares").AddForeignKey("monitoring_id", "monitorings(id)", "CASCADE", "CASCADE")

	for _, model := range DefaultModels {
		tableName := db.NewScope(model).GetModelStruct().TableName(db)
		db.FirstOrCreate(new(models.ContentType), models.ContentType{
			Table: tableName,
		})
	}

	logger.HandleError(db.Close())
}
