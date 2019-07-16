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
		"Группы мониторинга",
		"Типы мониторинга",
		"Рабочая группа",
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
		"Хабаровск, Комсомольск",
		"KVI, Первые цены",
		"ЦО, Самбери-9",
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
	monitoringGroupsNames := make([]string, 0)
	workGroupNames := make([]string, 0)
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
				case 3, 4, 5:
					names := koazee.StreamOf(strings.Split(text, ",")).Map(func(w string) string {
						return strings.Trim(w, " ")
					}).Out().Val().([]string)
					switch index {
					case 3:
						ware["monitoring_groups"] = names
						monitoringGroupsNames = append(monitoringGroupsNames, names...)
						break
					case 4:
						ware["monitoring_types"] = names
						monitoringTypeNames = append(monitoringTypeNames, names...)
						break
					case 5:
						ware["work_groups"] = names
						workGroupNames = append(workGroupNames, names...)
					}
				}

			}
			_wares[ware["code"].(string)] = ware
		}
	}
	monitoringTypeNames = koazee.StreamOf(monitoringTypeNames).RemoveDuplicates().Out().Val().([]string)
	monitoringGroupsNames = koazee.StreamOf(monitoringGroupsNames).RemoveDuplicates().Out().Val().([]string)
	workGroupNames = koazee.StreamOf(workGroupNames).RemoveDuplicates().Out().Val().([]string)

	db := context.Get(r, "DB").(*gorm.DB)
	var wares []models.Ware
	var monitoringTypes []models.MonitoringType
	var monitoringGroups []models.MonitoringGroups
	var workGroups []models.WorkGroup
	db.Find(&wares, "code IN (?)", wareCodes)
	db.Find(&monitoringTypes, "name IN (?)", monitoringTypeNames)
	db.Find(&monitoringGroups, "name IN (?)", monitoringGroupsNames)
	db.Find(&workGroups, "name IN (?)", workGroupNames)

	monitoringTypesStream := koazee.StreamOf(monitoringTypes)
	monitoringGroupsStream := koazee.StreamOf(monitoringGroups)
	workGroupsStream := koazee.StreamOf(workGroups)

	for _, ware := range wares {
		_ware := _wares[ware.Code]
		_ware["exist"] = true
		mt := monitoringTypesStream.Filter(func(mt models.MonitoringType) bool {
			res, _ := koazee.StreamOf(_ware["monitoring_types"].([]string)).IndexOf(mt.Name)
			return res > -1
		}).Out().Val()
		mg := monitoringGroupsStream.Filter(func(mt models.MonitoringGroups) bool {
			res, _ := koazee.StreamOf(_ware["monitoring_groups"].([]string)).IndexOf(mt.Name)
			return res > -1
		}).Out().Val()
		wg := workGroupsStream.Filter(func(mt models.WorkGroup) bool {
			res, _ := koazee.StreamOf(_ware["work_groups"].([]string)).IndexOf(mt.Name)
			return res > -1
		}).Out().Val()
		ware.Name = _ware["name"].(string)
		if mg != nil {
			ware.MonitoringGroups = mg.([]models.MonitoringGroups)
		}
		if mt != nil {
			ware.MonitoringType = mt.([]models.MonitoringType)
		}
		if wg != nil {
			ware.WorkGroups = wg.([]models.WorkGroup)
		}
		db.Save(ware)
	}
	for _, _ware := range _wares {
		if _ware["exist"] == nil {
			ware := models.Ware{
				Code: _ware["code"].(string),
				Name: _ware["name"].(string),
			}
			mt := monitoringTypesStream.Filter(func(mt models.MonitoringType) bool {
				res, _ := koazee.StreamOf(_ware["monitoring_types"].([]string)).IndexOf(mt.Name)
				return res > -1
			}).Out().Val()
			mg := monitoringGroupsStream.Filter(func(mt models.MonitoringGroups) bool {
				res, _ := koazee.StreamOf(_ware["monitoring_groups"].([]string)).IndexOf(mt.Name)
				return res > -1
			}).Out().Val()
			wg := workGroupsStream.Filter(func(mt models.WorkGroup) bool {
				res, _ := koazee.StreamOf(_ware["work_groups"].([]string)).IndexOf(mt.Name)
				return res > -1
			}).Out().Val()
			ware.Name = _ware["name"].(string)
			if mg != nil {
				ware.MonitoringGroups = mg.([]models.MonitoringGroups)
			}
			if mt != nil {
				ware.MonitoringType = mt.([]models.MonitoringType)
			}
			if wg != nil {
				ware.WorkGroups = wg.([]models.WorkGroup)
			}
			db.Save(ware)
		}
	}
	common.JSONResponse(w, types.H{
		"success": true,
	})
}
