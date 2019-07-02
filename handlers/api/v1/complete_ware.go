package v1

import (
	"encoding/json"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
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
		common.ErrorResponse(w, r, "Нет данных для выгрузки")
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
		region := models.MonitoringGroups{}
		if len(user.WorkGroup) > 0 && len(user.WorkGroup[0].MonitoringGroups) > 0 {
			region = user.WorkGroup[0].MonitoringGroups[0]
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
		completeWare.MonitoringGroup = region
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

func GetCompletedWares(w http.ResponseWriter, r *http.Request) {
	from, _ := time.Parse("2006-01-02", r.FormValue("from"))
	to, _ := time.Parse("2006-01-02", r.FormValue("to"))
	if from.IsZero() {
		from = time.Now()
	}
	if to.IsZero() {
		to = from
	}
	names := []string{"regions", "shops", "monitoring_types"}
	params := make(map[string][]int)
	for _, name := range names {
		params[name] = koazee.StreamOf(
			strings.Split(r.FormValue(name), ","),
		).Map(func(x string) int {
			n, _ := strconv.Atoi(x)
			return n
		}).Filter(func(x int) bool {
			return x > 0
		}).RemoveDuplicates().Out().Val().([]int)
	}
	limit := 250
	page, _ := strconv.Atoi(r.FormValue("page"))
	start := page*limit - limit
	db := context.Get(r, "DB").(*gorm.DB)
	qs := db.Model(
		&models.CompletedWare{},
	).Select(
		"DISTINCT "+
			"completed_wares.*,"+
			"CONCAT(u.user_name, ' ', u.last_name, ' ', u.first_name) as user,"+
			"s.name as segment,"+
			"w.name as ware, w.code as code,"+
			"ms.name as rival,"+
			"mt.name as monitoring_type,"+
			"mg.name as region",
	).Joins(
		"LEFT JOIN monitoring_groups mg ON mg.id = completed_wares.region_id",
	).Joins(
		"LEFT JOIN monitoring_shops ms ON ms.id = completed_wares.monitoring_shop_id",
	).Joins(
		"LEFT JOIN users u ON u.id = completed_wares.user_id",
	).Joins(
		"LEFT JOIN wares w ON w.id = completed_wares.ware_id",
	).Joins(
		"LEFT JOIN segments s ON s.id = w.segment_id",
	).Joins(
		"LEFT JOIN monitoring_types mt ON mt.id = completed_wares.monitoring_type_id",
	).Where(
		"completed_wares.date_upload BETWEEN date(?) AND (date(?) + '1 day'::interval)", from, to,
	)
	if len(params["regions"]) > 0 {
		qs = qs.Where("mg.id IN (?)", params["regions"])
	}
	if len(params["shops"]) > 0 {
		qs = qs.Where("ms.id IN (?)", params["shops"])
	}
	if len(params["monitoring_types"]) > 0 {
		qs = qs.Where("mt.id IN (?)", params["monitoring_types"])
	}
	type CompleteWare struct {
		ID             uint      `json:"-"`
		User           string    `json:"user"`
		DateUpload     time.Time `json:"date_upload"`
		Segment        string    `json:"segment"`
		Ware           string    `json:"ware"`
		Code           string    `json:"code"`
		Price          float64   `json:"price"`
		MaxPrice       float64   `json:"max_price"`
		MinPrice       float64   `json:"min_price"`
		Comment        string    `json:"comment"`
		Rival          string    `json:"rival"`
		Photos         []string  `json:"photos"`
		Region         string    `json:"region"`
		MonitoringType string    `json:"monitoring_type"`
	}
	completeWares := make([]*CompleteWare, 0)

	count := 0
	qs.Count(&count)
	qs.Offset(start).Limit(limit).Order("date_upload DESC").Scan(&completeWares)
	photos := make([]models.Photos, 0)

	if len(completeWares) > 0 {
		idx := koazee.StreamOf(completeWares).Map(func(x *CompleteWare) uint { return x.ID }).Out().Val()
		db.Find(&photos, "completed_ware_id IN (?)", idx)
	}

	var length int
	if length = limit; limit != len(completeWares) {
		length = len(completeWares)
	}
	for _, ware := range completeWares {
		for _, photo := range photos {
			if photo.CompletedWareId == ware.ID {
				ware.Photos = append(
					ware.Photos,
					fmt.Sprintf("%s/api/monitoring/media/%s", config.Config.System.ServerUrl, photo.Path),
				)
			}
		}
	}

	data := types.H{
		"paginate": common.PaginateInfo{
			CurrentPage: page,
			Count:       count,
			Length:      length,
			CountPage:   int(math.Ceil(float64(count) / float64(limit))),
		},
		"result": completeWares,
	}

	common.JSONResponse(w, data)
}
