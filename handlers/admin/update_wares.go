package admin

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"net/http"
)

func UpdateWares(w http.ResponseWriter, r *http.Request)  {
	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}

	xls, header, err := r.FormFile("update_wares")

	if err != nil {
		common.ErrorResponse(w, r, err.Error())
		return
	}

	// FIXME
	fmt.Print(xls)
	fmt.Print(header)
}
