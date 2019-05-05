package v1

import (
	"encoding/json"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func CompleteWare(w http.ResponseWriter, r *http.Request) {
	var data types.H
	if err := json.Unmarshal([]byte(r.PostFormValue("data")), &data); err != nil {
		common.ErrorResponse(w, r, "Ошибка выгрузки")
		return
	}
	user := context.Get(r, "user").(*models.User)
	_wares := data["wares"]
	if _wares == nil {
		common.ErrorResponse(w, r, "Нечего выгружать")
		return
	}
	wares := _wares.([]interface{})
	db := context.Get(r, "DB").(*gorm.DB)
	tx := db.Begin()
	for _, wareData := range wares {
		wareData := wareData.(map[string]interface{})
		completeWare := models.CompletedWare{}
		monitoringType := models.MonitoringType{}
		tx.Select(
			"DISTINCT monitoring_types.*",
		).Joins(
			"INNER JOIN monitoring_types_periods mtp ON mtp.monitoring_type_id = monitoring_types.id",
		).Joins(
			"INNER JOIN periods p ON mtp.period_id = p.id",
		).First(&monitoringType, "p.id = ?", wareData["period"])

		tx.FirstOrCreate(
			&completeWare,
			"ware_id = ? "+
				"AND user_id = ? "+
				"AND monitoring_shop_id = ? "+
				"AND monitoring_type_id = ?"+
				"AND date_upload BETWEEN current_date and (current_date + '1 day'::interval)",
			wareData["id"],
			user.ID,
			wareData["rival_id"],
			monitoringType.ID,
		)
		completeWare.MonitoringType = monitoringType
		tx.Preload("Segment").First(&completeWare.Ware, "id = ?", wareData["id"])
		tx.First(&completeWare.MonitoringShop, "id = ?", wareData["rival_id"])
		region := models.Regions{}
		if len(user.WorkGroup) > 0 && len(user.WorkGroup[0].Regions) > 0 {
			region = user.WorkGroup[0].Regions[0]
		}
		completeWare.DateUpload = time.Now()
		completeWare.User = *user
		if err := helpers.SetFieldsForModel(&completeWare, wareData); err != nil {
			common.ErrorResponse(w, r, err.Error())
			tx.Rollback()
			return
		}
		for _, _path := range wareData["photos"].([]interface{}) {
			photo := models.Photos{Path: _path.(string)}
			tx.FirstOrCreate(&photo, photo)
			completeWare.Photos = append(completeWare.Photos, photo)
		}
		completeWare.Region = region
		tx.Save(&completeWare)
	}
	if res := tx.Commit(); res.Error != nil {
		tx.Rollback()
		common.ErrorResponse(w, r, res.Error.Error())
	} else {
		common.JSONResponse(w, types.H{
			"success": true,
		})
	}
}

func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	photo, header, err := r.FormFile("photo")
	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}

	filePath := path.Join(config.Config.Static.MediaRoot, header.Filename)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	if _, err = io.Copy(f, photo); err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}

	defer func() {
		logger.HandleError(photo.Close())
		logger.HandleError(f.Close())
	}()
	common.JSONResponse(w, types.H{
		"error":     false,
		"url_photo": filePath,
	})

}
