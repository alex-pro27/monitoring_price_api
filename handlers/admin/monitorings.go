package admin

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetAllMonitoringList(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	user := context.Get(r, "user").(*models.User)
	is_admin := false
	for _, r := range user.Roles {
		if is_admin = r.RoleType == models.IS_ADMIN; is_admin {
			break
		}
	}
	monitorings := make([]models.Monitoring, 0)
	if is_admin {
		db.Preload(
			"Region",
		).Order(
			"\"name\", cast(NULLIF(regexp_replace(\"name\", E'\\\\D', '', 'g'), '') AS integer)",
		).Find(
			&monitorings,
		)
	} else {
		for _, wg := range user.WorkGroups {
			monitorings = append(monitorings, wg.Monitorings...)
		}
	}
	data := make([]types.H, 0)
	for _, it := range monitorings {
		data = append(data, it.Serializer())
	}
	common.JSONResponse(w, data)
}
