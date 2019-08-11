package admin

import (
	"bytes"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/otium/queue"
	"github.com/tealeg/xlsx"
	"github.com/wesovilabs/koazee"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func GetTemplateBlank(w http.ResponseWriter, r *http.Request) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	header := []string{
		"ШК товара",
		"Код товара",
		"Наименование",
		"Сегмент",
		"Типы мониторинга",
	}
	row = sheet.AddRow()
	headerStyle := xlsx.Style{
		Font: xlsx.Font{
			Bold: true,
		},
	}
	for _, title := range header {
		cell = row.AddCell()
		cell.SetStyle(&headerStyle)
		cell.Value = title
	}
	row = sheet.AddRow()
	example := []string{
		"",
		"33982",
		"Майонез Провансаль 67% 500г п/ст МЖК Хабаровск",
		"130 МАСЛО, МАРГАРИН, МАЙОНЕЗЫ",
		"KVI, Первые цены",
	}
	for _, title := range example {
		cell = row.AddCell()
		cell.Value = title
	}
	file_path := path.Join(os.TempDir(), "products_blank.xlsx")
	err = file.Save(file_path)
	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	f, err := os.Open(file_path)
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
	defer logger.HandleError(f.Close())
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
		xlFile, user, isUpdate,
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

	db := databases.ConnectDefaultDB()
	tx := db.Begin()

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			logger.Logger.Errorf("Error parse product list xls file: %v", rec)
			AdminWebSocket.Emit(user.Token.Key, "on_update_products", types.H{
				"error":   true,
				"message": fmt.Sprintf("Ошибка обновления мониторинга: %v", rec),
			})
		} else {
			AdminWebSocket.Emit(user.Token.Key, "on_update_products", types.H{
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

	__monitoringTypes := make([]models.MonitoringType, 0)
	tx.Find(&__monitoringTypes, "name IN (?)", monitoringTypeNames)
	monitoringTypes := make(map[uint]*models.MonitoringType)
	for _, mt := range __monitoringTypes {
		monitoringTypes[mt.ID] = &mt
	}

	__wares := make([]models.Ware, 0)
	tx.Find(&__wares, "code IN (?)", wareCodes)
	wares := make(map[string]models.Ware)
	for _, w := range __wares {
		wares[w.Code] = w
	}

	__segments := make([]models.Segment, 0)
	segments := make(map[string]*models.Segment, 0)
	tx.Find(&__segments, "code IN (?)", segmentsCodes)

	for _, s := range __segments {
		segments[s.Code] = &s
	}

	waresByMonitoringTypes := make(map[uint][]models.Ware)
	userWorkGroupIDX := make([]uint, 0)

	for _, wg := range user.WorkGroup {
		userWorkGroupIDX = append(userWorkGroupIDX, wg.ID)
	}

	for _code, _ware := range _wares {
		ware := wares[_code]
		segment := segments[_ware["segment_code"].(string)]
		if segment == nil {
			segment = new(models.Segment)
			segment.Name = _ware["segment"].(string)
			segment.Code = _ware["segment_code"].(string)
			if err := tx.Save(segment).Error; err != nil {
				panic(fmt.Sprintf("Не удалось добавить сегмент: %s %s, %s", segment.Name, segment.Code, err))
			}
		}
		if ware.ID == 0 {
			ware = models.Ware{}
			ware.Code = _code
			ware.Name = _ware["name"].(string)
			ware.Barcode = _ware["barcode"].(string)
			ware.SegmentId = segment.ID
			if err := tx.Save(&ware).Error; err != nil {
				panic(fmt.Sprintf("Не удалось добавить товар: %s %s, %v", _code, _ware["name"], err))
			}
		}

		for _, mt := range __monitoringTypes {
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
	).Joins(
		"INNER JOIN monitoring_shops_monitorings msm ON msm.monitoring_id = monitorings.id",
	).Joins(
		"INNER JOIN work_groups_monitorings wgm ON wgm.monitoring_id = monitorings.id",
	).Find(&monitorings, "mt.name IN (?) AND wgm.work_group_id IN (?)", monitoringTypeNames, userWorkGroupIDX)

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
