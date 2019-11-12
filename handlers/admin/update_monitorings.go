package admin

import (
	"bytes"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/alex-pro27/monitoring_price_api/utils"
	"github.com/gorilla/context"
	"github.com/otium/queue"
	"github.com/tealeg/xlsx"
	"github.com/wesovilabs/koazee"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetTemplateBlank(w http.ResponseWriter, r *http.Request) {
	header := []string{
		"ШК товара",
		"Код товара",
		"Наименование",
		"Сегмент",
		"Типы мониторинга",
	}
	example := []string{
		"",
		"33982",
		"Майонез Провансаль 67% 500г п/ст МЖК Хабаровск",
		"130 МАСЛО, МАРГАРИН, МАЙОНЕЗЫ",
		"KVI, Первые цены",
	}

	filePath, err := utils.CreateXLSX(header, [][]string{example}, "")
	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	f, err := os.Open(filePath)
	defer func() {
		defer logger.HandleError(f.Close())
	}()
	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	buffer := new(bytes.Buffer)
	bufferBytes := make([]byte, 0)
	if _, err := io.Copy(buffer, f); err != nil {
		panic(err)
	}
	bufferBytes = buffer.Bytes()
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Length", strconv.Itoa(len(bufferBytes)))
	w.Header().Set("Content-Disposition", "attachment; filename=products_blank.xlsx")
	_, err = w.Write(bufferBytes)
	logger.HandleError(err)
}

func UpdateMonitorings(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	mFile, header, err := r.FormFile("update_wares")
	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	if header.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		common.ErrorResponse(w, r, "Неверный формат файла")
		return
	}
	isUpdate, _ := strconv.ParseBool(r.FormValue("is_update"))
	monitoringIDX := koazee.StreamOf(strings.Split(r.FormValue("monitoring_idx"), ",")).Map(func(x string) uint {
		res, _ := strconv.ParseUint(x, 10, 64)
		return uint(res)
	}).Out().Val().([]uint)

	if len(monitoringIDX) == 0 {
		common.ErrorResponse(w, r, "Не выбранны мониториги")
		return
	}
	user := context.Get(r, "user").(*models.User)

	buffer := new(bytes.Buffer)
	bufferBytes := make([]byte, 0)
	if _, err := io.Copy(buffer, mFile); err != nil {
		logger.HandleError(err)
		return
	}
	bufferBytes = buffer.Bytes()
	xlFile, err := xlsx.OpenBinary(bufferBytes)

	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	q := queue.NewQueue(taskUpdateMonitoring, 0)
	q.Push([]interface{}{
		xlFile, user, isUpdate, monitoringIDX,
	})
	common.JSONResponse(w, types.H{
		"success": true,
	})
}

func taskUpdateMonitoring(args interface{}) {
	_args := args.([]interface{})
	xlFile := _args[0].(*xlsx.File)
	user := _args[1].(*models.User)
	isUpdate := _args[2].(bool)
	monitoringIDX := _args[3].([]uint)

	db := databases.ConnectDefaultDB()
	tx := db.Begin()

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			logger.Logger.Errorf("Error parse product list xls file: %v", rec)
			AdminWebSocket.Emit(user.Token.Key, "update_products", types.H{
				"error":   true,
				"message": fmt.Sprintf("Ошибка обновления мониторинга: %v", rec),
			})
		} else {
			AdminWebSocket.Emit(user.Token.Key, "update_products", types.H{
				"message": "Мониториг успешно обновлен!",
			})
		}
		logger.HandleError(db.Close())
	}()

	wareCodes := make([]string, 0)
	wareNames := make([]string, 0)
	monitoringTypeNames := make([]string, 0)
	segmentsCodes := make([]string, 0)
	pattern := regexp.MustCompile("^(?P<Code>\\d+)\\s*(?P<Name>.+)")
	_wares := make(map[string]map[string]interface{}, 0)

	for _, sheet := range xlFile.Sheets {
	INNER_CYCLE:
		for i, row := range sheet.Rows {
			if i == 0 || len(row.Cells) < 4 {
				continue
			}
			ware := make(map[string]interface{})

			for index, cell := range row.Cells {
				text := cell.String()
				switch index {
				case 0:
					ware["barcode"] = text
					break
				case 1:
					if text == "" {
						continue INNER_CYCLE
					}
					ware["code"] = text
					wareCodes = append(wareCodes, text)
					break
				case 2:
					if text == "" {
						continue INNER_CYCLE
					}
					ware["name"] = text
					wareNames = append(wareNames, text)
					break
				case 3:
					seg := pattern.FindStringSubmatch(text)
					if len(seg) > 2 {
						ware["segment"] = seg[2]
						ware["segment_code"] = seg[1]
						segmentsCodes = append(segmentsCodes, seg[1])
					} else {
						continue INNER_CYCLE
					}
					break
				case 4:
					names := koazee.StreamOf(strings.Split(text, ",")).Map(func(w string) string {
						return strings.Trim(w, " ")
					}).Out().Val().([]string)
					ware["monitoring_types"] = names
					monitoringTypeNames = append(monitoringTypeNames, names...)
					break
				}
			}
			_wares[ware["code"].(string)] = ware
		}
	}
	monitoringTypeNames = koazee.StreamOf(monitoringTypeNames).RemoveDuplicates().Out().Val().([]string)
	segmentsCodes = koazee.StreamOf(segmentsCodes).RemoveDuplicates().Out().Val().([]string)
	wareCodes = koazee.StreamOf(wareCodes).RemoveDuplicates().Out().Val().([]string)
	wareNames = koazee.StreamOf(wareNames).RemoveDuplicates().Out().Val().([]string)

	monitoringTypes := make([]models.MonitoringType, 0)
	tx.Find(&monitoringTypes, "name IN (?)", monitoringTypeNames)

	__wares := make([]models.Ware, 0)
	tx.Find(&__wares, "code IN (?)", wareCodes)
	wares := make(map[string]models.Ware)
	for _, w := range __wares {
		wares[w.Code] = w
	}

	__segments := make([]models.Segment, 0)
	segments := make(map[string]models.Segment, 0)
	tx.Find(&__segments, "code IN (?)", segmentsCodes)

	for _, s := range __segments {
		segments[s.Code] = s
	}

	waresByMonitoringTypes := make(map[uint][]models.Ware)
	userMonitoringIDX := make([]uint, 0)
	isAdmin := false

	for _, r := range user.Roles {
		if r.RoleType == models.IS_ADMIN {
			isAdmin = true
			break
		}
	}

	allMonitorings := make([]models.Monitoring, 0)

	if isAdmin {
		tx = tx.Find(&allMonitorings)
	} else {
		for _, wg := range user.WorkGroups {
			for _, m := range wg.Monitorings {
				allMonitorings = append(allMonitorings, m)
			}
		}
	}

	for _, m := range allMonitorings {
		for i, mID := range monitoringIDX {
			if mID == m.ID {
				userMonitoringIDX = append(userMonitoringIDX, m.ID)
				monitoringIDX = append(monitoringIDX[:i], monitoringIDX[i+1:]...)
			}
		}
	}

	for _code, _ware := range _wares {
		ware := wares[_code]
		notWare := ware.ID == 0
		segment := segments[_ware["segment_code"].(string)]
		notSegment := segment.ID == 0
		segment.Name = _ware["segment"].(string)
		segment.Code = _ware["segment_code"].(string)
		if err := tx.Save(segment).Error; err != nil {
			panic(fmt.Sprintf("Не удалось добавить сегмент: %s %s, %s", segment.Name, segment.Code, err))
		} else if notSegment {
			segments[_ware["segment_code"].(string)] = segment
		}

		ware.Code = _code
		ware.Name = _ware["name"].(string)
		ware.Barcode = _ware["barcode"].(string)
		ware.SegmentId = segment.ID
		if err := tx.Save(&ware).Error; err != nil {
			panic(fmt.Sprintf("Не удалось добавить товар: %s %s, %v", _code, _ware["name"], err))
		} else if notWare {
			wares[_code] = ware
		}

		for _, mt := range monitoringTypes {
			for _, mt_name := range _ware["monitoring_types"].([]string) {
				if mt_name == mt.Name {
					if waresByMonitoringTypes[mt.ID] == nil {
						waresByMonitoringTypes[mt.ID] = make([]models.Ware, 0)
					}
					waresByMonitoringTypes[mt.ID] = append(waresByMonitoringTypes[mt.ID], ware)
				}
			}
		}
	}
	var monitorings []models.Monitoring
	tx.Preload("Wares").Select("DISTINCT monitorings.*").Joins(
		"INNER JOIN monitoring_types mt ON mt.id = monitoring_type_id",
	).Find(&monitorings, "mt.name IN (?) AND monitorings.id IN (?)", monitoringTypeNames, userMonitoringIDX)

	for _, monitoring := range monitorings {
		wares := waresByMonitoringTypes[monitoring.MonitoringTypeId]
		if wares == nil {
			continue
		}
		if !isUpdate {
			tx.Model(&monitoring).Association("Wares").Replace(wares)
		} else {
			monitoring.Wares = append(monitoring.Wares, wares...)
			if err := tx.Save(&monitoring).Error; err != nil {
				panic(fmt.Sprintf("Не удалось обновить мониторинг %s: %v", monitoring.Name, err))
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		panic(fmt.Sprintf("Не удалось обновить мониторинг: %v", err))
	}
}
