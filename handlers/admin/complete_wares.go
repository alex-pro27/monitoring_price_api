package admin

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetCompletedWares(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*models.User)
	allowAllData := -1
	if !user.IsSuperUser {
		for _, it := range user.Roles {
			if it.RoleType == models.IS_ADMIN {
				allowAllData = 1
				break
			} else if it.RoleType == models.IS_MANAGER {
				allowAllData = 0
			}
		}
	} else {
		allowAllData = 1
	}
	from, _ := time.Parse("2006-01-02", r.FormValue("datefrom"))
	to, _ := time.Parse("2006-01-02", r.FormValue("dateto"))
	names := []string{"regions", "monitoring_shops", "monitoring_types", "work_groups"}
	orderBy := r.FormValue("order_by")
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
		"completed_wares.*," +
			"u.user_name as user_barcode," +
			"CONCAT(u.last_name, ' ', u.first_name) as user_name," +
			"s.name as segment," +
			"s.code as segment_code," +
			"w.name as ware, " +
			"w.code as code," +
			"ms.name as rival," +
			"ms.code as rival_code," +
			"mt.name as monitoring_type," +
			"r.name as region," +
			"wg.name as work_group",
	).Joins(
		"LEFT JOIN regions r ON r.id = completed_wares.region_id",
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
	).Joins(
		"INNER JOIN work_groups_users wgu ON wgu.user_id = u.id",
	).Joins(
		"INNER JOIN work_groups wg ON wg.id = wgu.work_group_id",
	)
	if !(from.IsZero() && to.IsZero()) {
		qs = qs.Where(
			"completed_wares.date_upload BETWEEN date(?) AND (date(?) + '1 day'::interval)", from, to,
		)
	}

	if len(params["regions"]) > 0 {
		qs = qs.Where("r.id IN (?)", params["regions"])
	}
	if len(params["monitoring_shops"]) > 0 {
		qs = qs.Where("ms.id IN (?)", params["monitoring_shops"])
	}
	if len(params["monitoring_types"]) > 0 {
		qs = qs.Where("mt.id IN (?)", params["monitoring_types"])
	}

	if allowAllData == 0 {
		workGroupIDX := make([]uint, 0)
		for _, wg := range user.WorkGroups {
			workGroupIDX = append(workGroupIDX, wg.ID)
		}
		qs = qs.Where("wg.id IN (?)", workGroupIDX)
	} else if allowAllData == 1 && len(params["work_groups"]) > 0 {
		qs = qs.Where("wg.id IN (?)", params["work_groups"])
	} else if allowAllData == -1 {
		qs = qs.Where("u.user_name = ?", user.UserName)
	}
	if orderBy != "" {
		orderMap := map[string]string{
			"ware":         "w.name",
			"code":         "w.code",
			"user_name":    "u.last_name",
			"segment":      "s.name",
			"segment_code": "s.code",
			"region":       "r.name",
			"rival":        "ms.name",
			"rival_code":   "ms.code",
			"price":        "completed_wares.price",
			"date_upload":  "completed_wares.date_upload",
		}
		order := ""
		for _, ordName := range strings.Split(orderBy, ",") {
			ord := " ASC,"
			if match, _ := regexp.MatchString("^-", ordName); match {
				ord = " DESC,"
				ordName = ordName[1:]
			}
			if orderMap[ordName] != "" {
				order += orderMap[ordName] + ord
			}
		}
		if order != "" {
			qs = qs.Order(order[:len(order)-1])
		}
	}

	searchKeyWords := r.FormValue("keywords")
	if searchKeyWords != "" {
		keywords := strings.Split(searchKeyWords, " ")
		for _, fieldName := range []string{"w.name", "w.code", "s.name", "CONCAT(u.last_name, ' ', u.first_name)"} {
			isOr := len(keywords) < 2
			for _, keyword := range keywords {
				if keyword == "" {
					continue
				}
				if !isOr {
					qs = qs.Where(fmt.Sprintf("%s ilike ?", fieldName), "%"+keyword+"%")
					isOr = true
				} else {
					qs = qs.Or(fmt.Sprintf("%s ilike ?", fieldName), "%"+keyword+"%")
				}
			}
		}
	}

	type CompleteWare struct {
		ID             uint      `json:"id"`
		UserBarcode    string    `json:"user_barcode"`
		UserName       string    `json:"user_name"`
		DateUpload     time.Time `json:"date_upload"`
		Segment        string    `json:"segment"`
		SegmentCode    string    `json:"segment_code"`
		Ware           string    `json:"ware"`
		Code           string    `json:"code"`
		Price          float64   `json:"price"`
		MaxPrice       float64   `json:"max_price"`
		MinPrice       float64   `json:"min_price"`
		Discount       bool      `json:"discount"`
		Missing        bool      `json:"missing"`
		Comment        string    `json:"comment"`
		Rival          string    `json:"rival"`
		RivalCode      string    `json:"rival_code"`
		Photos         []string  `json:"photos"`
		Region         string    `json:"region"`
		WorkGroup      string    `json:"work_group"`
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
				ware.Photos = append(ware.Photos, "/api/admin/media/"+photo.Path)
			}
		}
	}

	if page == 0 {
		page = 1
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
