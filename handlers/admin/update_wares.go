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
		"Группы мониторинга",
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
		"МАСЛО, МАРГАРИН, МАЙОНЕЗЫ",
		"Хабаровск, Комсомольск",
		"KVI, Первые цены",
	}
	for _, title := range example {
		cell = row.AddCell()
		cell.Value = title
	}
	file_path := path.Join(os.TempDir(), "Products_blank.xlsx")
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
	w.Header().Set("Content-Disposition", "attachment; filename=Products_blank.xlsx")
	_, err = w.Write(bufferBytes)
	logger.HandleError(err)
	defer logger.HandleError(f.Close())
}

func UpdateWares(w http.ResponseWriter, r *http.Request) {

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
	user := context.Get(r, "user").(*models.User)
	q := queue.NewQueue(taskUpdateWare, 0)

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

	q.Push([]interface{}{
		xlFile, user.Token.Key,
	})
	common.JSONResponse(w, types.H{
		"success": true,
	})
}


func taskUpdateWare(args interface{}) {

	xlFile := args.([]interface{})[0].(*xlsx.File)
	token := args.([]interface{})[1].(string)
	db := databases.ConnectDefaultDB()
	tx := db.Begin()

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			logger.Logger.Errorf("Error parse product list xls file: %v", rec)
			if AdminWebSocket != nil {
				AdminWebSocket.Emit(token, "on_update_products", types.H{
					"error": true,
					"message": fmt.Sprintf("Ошибка обновления товаров: %v", rec),
				})
			}
		} else {
			if AdminWebSocket != nil {
				AdminWebSocket.Emit(token, "on_update_products", types.H{
					"message": "Товары успешно обновлены!",
				})
			}
		}
		logger.HandleError(db.Close())
	}()

	wareCodes := make([]string, 0)
	wareNames := make([]string, 0)
	monitoringTypeNames := make([]string, 0)
	monitoringGroupNames := make([]string, 0)
	segmentsCodes := make([]string, 0)

	pattern := regexp.MustCompile("^(?P<Code>\\d+)\\s*(?P<Name>.+)")
	_wares := make(map[string]map[string]interface{}, 0)
	for _, sheet := range xlFile.Sheets {
	INNER:
		for i, row := range sheet.Rows {
			if i == 0 {
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
						continue INNER
					}
					ware["code"] = text
					wareCodes = append(wareCodes, text)
					break
				case 2:
					if text == "" {
						continue INNER
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
						continue INNER
					}
					break
				case 4, 5:
					names := koazee.StreamOf(strings.Split(text, ",")).Map(func(w string) string {
						return strings.Trim(w, " ")
					}).Out().Val().([]string)
					switch index {
					case 4:
						ware["monitoring_groups"] = names
						monitoringGroupNames = append(monitoringGroupNames, names...)
						break
					case 5:
						ware["monitoring_types"] = names
						monitoringTypeNames = append(monitoringTypeNames, names...)
						break
					}
				}
			}
			_wares[ware["code"].(string)] = ware
		}
	}
	monitoringTypeNames = koazee.StreamOf(monitoringTypeNames).RemoveDuplicates().Out().Val().([]string)
	monitoringGroupNames = koazee.StreamOf(monitoringGroupNames).RemoveDuplicates().Out().Val().([]string)
	segmentsCodes = koazee.StreamOf(segmentsCodes).RemoveDuplicates().Out().Val().([]string)
	wareCodes = koazee.StreamOf(wareCodes).RemoveDuplicates().Out().Val().([]string)
	wareNames = koazee.StreamOf(wareNames).RemoveDuplicates().Out().Val().([]string)

	var wares []models.Ware
	var monitoringTypes []models.MonitoringType
	segments := make([]models.Segment, 0)

	tx.Find(&segments, "code IN (?)", segmentsCodes)
	segmentsStream := koazee.StreamOf(segments)

	type MS struct {
		ID 			uint
		GroupName 	string
	}

	monitoringShopsIDGroupName := []MS{}

	tx.Model(models.MonitoringShop{}).Select("DISTINCT monitoring_shops.id id, mg.name group_name").Joins(
		"INNER JOIN work_groups_monitoring_shops wgms ON wgms.monitoring_shop_id = monitoring_shops.id",
	).Joins(
		"INNER JOIN work_groups_monitoring_groups wgmg ON wgmg.work_group_id = wgms.work_group_id",
	).Joins(
		"INNER JOIN monitoring_groups mg ON mg.id = wgmg.monitoring_groups_id",
	).Where("mg.name IN (?)", monitoringGroupNames).Scan(&monitoringShopsIDGroupName)

	tx.Find(&wares, "code IN (?) AND name IN (?)", wareCodes, wareNames)
	tx.Find(&monitoringTypes, "name IN (?)", monitoringTypeNames)

	monitoringTypesStream := koazee.StreamOf(monitoringTypes)
	monitoringShopsIDX := make([]uint, 0)
	var monitoringShops []models.MonitoringShop
	if len(monitoringShopsIDGroupName) > 0 {
		monitoringShopsIDX = koazee.StreamOf(
			monitoringShopsIDGroupName,
		).Map(func(ms MS) uint { return ms.ID }).RemoveDuplicates().Out().Val().([]uint)
		tx.Find(&monitoringShops, "id IN (?)", monitoringShopsIDX)
	}

	_monitoringShopsByGroup := make(map[string][]models.MonitoringShop)

	for _, mg := range monitoringShopsIDGroupName {
		for _, monitoringShop := range monitoringShops {
			if monitoringShop.ID == mg.ID {
				if _monitoringShopsByGroup[mg.GroupName] == nil {
					_monitoringShopsByGroup[mg.GroupName] = make([]models.MonitoringShop, 0)
				}
				_monitoringShopsByGroup[mg.GroupName] = append(_monitoringShopsByGroup[mg.GroupName], monitoringShop)
			}
		}
	}
	for _, ware := range wares {
		_ware := _wares[ware.Code]
		_ware["exist"] = true
		var segment models.Segment
		_segment := segmentsStream.Filter(func(s models.Segment) bool {return s.Code == _ware["segment_code"].(string)}).Out().Val()
		if len(_segment.([]models.Segment)) > 0 {
			segment = _segment.([]models.Segment)[0]
		} else {
			segment = models.Segment{
				Name: _ware["segment"].(string),
				Code: _ware["segment_code"].(string),
			}
		}
		if ans := tx.Save(&segment); ans.Error != nil {
			tx.Rollback()
			logger.HandleError(ans.Error)
			return
		}
		ware.Segment = segment

		mt := monitoringTypesStream.Filter(func(mt models.MonitoringType) bool {
			res, _ := koazee.StreamOf(_ware["monitoring_types"].([]string)).IndexOf(mt.Name)
			return res > -1
		}).Out().Val()

		ware.Name = _ware["name"].(string)
		if mt != nil {
			ware.MonitoringType = mt.([]models.MonitoringType)
			db.Model(&ware).Association("MonitoringType").Replace(mt.([]models.MonitoringType))
		}
		for _, mg := range _ware["monitoring_groups"].([]string) {
			db.Model(&ware).Association("MonitoringShops").Replace(_monitoringShopsByGroup[mg])
		}

		if ans := tx.Save(&ware); ans.Error != nil {
			logger.HandleError(ans.Error)
			return
		}
	}

	for _, _ware := range _wares {
		if _ware["exist"] == nil {
			ware := models.Ware{
				Code: _ware["code"].(string),
				Name: _ware["name"].(string),
			}
			var segment models.Segment
			var flagSegmentIsCreated bool = false
			for _, s := range segments {
				if s.Code == _ware["segment_code"].(string) {
					segment = s
				}
			}
			if segment.ID == 0 {
				segment = models.Segment{
					Name: _ware["segment"].(string),
					Code: _ware["segment_code"].(string),
					Active: true,
				}
				flagSegmentIsCreated  = true
			}

			if ans := tx.Save(&segment); ans.Error != nil {
				tx.Rollback()
				logger.HandleError(ans.Error)
				return
			}
			if flagSegmentIsCreated {
				segments = append(segments, segment)
			}
			ware.Segment = segment

			mt := monitoringTypesStream.Filter(func(mt models.MonitoringType) bool {
				res, _ := koazee.StreamOf(_ware["monitoring_types"].([]string)).IndexOf(mt.Name)
				return res > -1
			}).Out().Val()
			ware.Name = _ware["name"].(string)
			if mt != nil {
				ware.MonitoringType = mt.([]models.MonitoringType)
			}
			for _, mg := range _ware["monitoring_groups"].([]string) {
				ware.MonitoringShops = _monitoringShopsByGroup[mg]
			}

			if ans := tx.Save(&ware); ans.Error != nil {
				tx.Rollback()
				logger.HandleError(ans.Error)
				return
			}
		}
	}
	tx.Commit()
}