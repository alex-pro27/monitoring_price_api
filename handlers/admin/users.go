package admin

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

/**
Список пользователей
*/
func AllUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	regionID, _ := strconv.Atoi(r.FormValue("region"))
	workGroupsID, _ := strconv.Atoi(r.FormValue("work_groups"))

	db := context.Get(r, "DB").(*gorm.DB)
	qs := db
	var users []models.User

	if workGroupsID != 0 && regionID == 0 {
		qs = qs.Joins(
			"INNER JOIN users_work_groups uw ON users.id = uw.user_id",
		).Where(
			"uw.work_group_id = ?", workGroupsID,
		)
	}

	if regionID != 0 {
		qs = qs.Joins(
			"INNER JOIN users_work_groups uw ON users.id = uw.user_id",
		).Joins(
			"INNER JOIN work_groups_monitoring_groups wgmg ON uw.work_group_id = wgmg.work_group_id",
		).Where(
			"wgmg.monitoring_groups_id = ?", regionID,
		)
		if workGroupsID != 0 {
			qs = qs.Where("uw.work_group_id = ?", workGroupsID)
		}
	}
	qs = qs.Order("last_name")

	data := common.Paginate(&users, qs, page, 100, []string{}, true)

	if len(data.Result) == 0 {
		common.Error404(w, r)
		return
	}
	common.JSONResponse(w, data)
}
