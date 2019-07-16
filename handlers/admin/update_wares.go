package admin

import (
	"bytes"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
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
	itoa := strconv.Itoa(len(bufferBytes))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Length", itoa)
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

	xls_file, header, err := r.FormFile("update_wares")

	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}

	if header.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		common.ErrorResponse(w, r, "Неверный формат файла")
		return
	}

	buffer := new(bytes.Buffer)
	bufferBytes := make([]byte, 0)
	if _, err := io.Copy(buffer, xls_file); err != nil {
		panic(err)
	}
	bufferBytes = buffer.Bytes()
	xlFile, err := xlsx.OpenBinary(bufferBytes)

	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}
	wareCodes := make([]string, 0)
	monitoringTypeNames := make([]string, 0)
	monitoringGroupNames := make([]string, 0)
	segmentsCodes := make([]string, 0)

	pattern := regexp.MustCompile("^(?P<Code>\\d+)\\s*(?P<Name>.+)")
	_wares := make(map[string]map[string]interface{}, 0)
	for _, sheet := range xlFile.Sheets {
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
					ware["code"] = text
					wareCodes = append(wareCodes, text)
					break
				case 2:
					ware["name"] = text
					break
				case 3:
					seg := pattern.FindStringSubmatch(text)
					if len(seg) > 2 {
						ware["segment"] = seg[2]
						ware["segment_code"] = seg[1]
						segmentsCodes = append(segmentsCodes, seg[1])
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

	db := context.Get(r, "DB").(*gorm.DB)
	var wares []models.Ware
	var monitoringTypes []models.MonitoringType
	var segments []models.Segment

	db.Find(&segments, "code IN (?)", segmentsCodes)
	segmentsStream := koazee.StreamOf(segments)

	type MS struct {
		ID 			uint
		GroupName 	string
	}

	monitoringShopsIDGroupName := []MS{}

	db.Model(models.MonitoringShop{}).Select("DISTINCT monitoring_shops.id id, mg.name group_name").Joins(
		"INNER JOIN work_groups_monitoring_shops wgms ON wgms.monitoring_shop_id = monitoring_shops.id",
	).Joins(
		"INNER JOIN work_groups_monitoring_groups wgmg ON wgmg.work_group_id = wgms.work_group_id",
	).Joins(
		"INNER JOIN monitoring_groups mg ON mg.id = wgmg.monitoring_groups_id",
	).Where("mg.name IN (?)", monitoringGroupNames).Scan(&monitoringShopsIDGroupName)

	db.Find(&wares, "code IN (?)", wareCodes)
	db.Find(&monitoringTypes, "name IN (?)", monitoringTypeNames)

	monitoringTypesStream := koazee.StreamOf(monitoringTypes)
	monitoringShopsIDGroupNameStream := koazee.StreamOf(monitoringShopsIDGroupName)
	monitoringShopsIDX := monitoringShopsIDGroupNameStream.Map(func(ms MS) uint {return ms.ID}).RemoveDuplicates().Out().Val().([]uint)

	var monitoringShops []models.MonitoringShop
	db.Find(&monitoringShops, "id IN (?)", monitoringShopsIDX)

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
	tx := db.Begin()
	for _, ware := range wares {
		_ware := _wares[ware.Code]
		_ware["exist"] = true
		segment := segmentsStream.Filter(func(s models.Segment) bool {return s.Code == _ware["segment_code"]}).Out().Val()
		if segment != nil && len(segment.([]models.Segment)) > 0 {
			segment = segment.([]models.Segment)[0]
 		} else {
 			segment = models.Segment{
 				Name: _ware["segment"].(string),
 				Code: _ware["segment_code"].(string),
			}
		}
		tx.Save(&segment)
		ware.Segment = segment.(models.Segment)

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

		tx.Save(&ware)
	}

	for _, _ware := range _wares {
		if _ware["exist"] == nil {
			ware := models.Ware{
				Code: _ware["code"].(string),
				Name: _ware["name"].(string),
			}

			segment := segmentsStream.Filter(func(s models.Segment) bool {return s.Code == _ware["segment_code"]}).Out().Val()
			if segment != nil && len(segment.([]models.Segment)) > 0 {
				segment = segment.([]models.Segment)[0]
			} else {
				segment = models.Segment{
					Name: _ware["segment"].(string),
					Code: _ware["segment_code"].(string),
				}
			}
			tx.Save(&segment)
			ware.Segment = segment.(models.Segment)

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

			tx.Save(&ware)
		}
	}

	tx.Commit()

	common.JSONResponse(w, types.H{
		"success": true,
	})
}