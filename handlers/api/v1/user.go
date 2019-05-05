package v1

import (
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CheckPin(w http.ResponseWriter, r *http.Request) {
	pincode := r.PostFormValue("pincode")
	if pincode == config.Config.System.PinCode {
		common.JSONResponse(w, types.H{"auth": true})
	} else {
		common.Forbidden(w, r)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	barcode := vars["barcode"]
	db := context.Get(r, "DB").(*gorm.DB)
	user := models.User{}
	user.Manager(db).GetByUserName(barcode)
	if user.ID == 0 {
		common.Error404(w, r)
		return
	}
	var region, shopName, roleName string
	if len(user.WorkGroup) > 0 {
		shopName = user.WorkGroup[0].Name
		if len(user.WorkGroup[0].Regions) > 0 {
			region = user.WorkGroup[0].Regions[0].Name
		}
	}
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
	}

	data := types.H{
		"id":       user.ID,
		"username": user.UserName,
		"name":     user.GetFullName(),
		"email":    user.Email,
		"region":   region,
		"shop":     shopName,
		"barcode":  user.UserName,
		"type":     roleName,
		"token":    user.Token.Key,
	}
	common.JSONResponse(w, data)
}
