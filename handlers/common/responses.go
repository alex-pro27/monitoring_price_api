package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/alex-pro27/monitoring_price_api/utils"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
)

func JSONResponse(w http.ResponseWriter, data interface{}) {
	var (
		body []byte
		err  error
	)
	if config.Config.System.Debug {
		body, err = json.MarshalIndent(data, "", "	")
	} else {
		body, err = json.Marshal(data)
	}

	if err != nil {
		log.Printf("Failed to encode a JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write the response body: %v", err)
		return
	}
}

func FileResponse(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	media := config.Config.Static.MediaRoot
	isThumb, _ := regexp.MatchString(".+_thumb\\.(jpe?g|png|gif)", name)
	var f *os.File
	f, err := os.Open(path.Join(media, name))
	buffer := new(bytes.Buffer)
	bufferBytes := make([]byte, 0)
	defer func() {
		buffer.Truncate(0)
	}()
	if err != nil && !isThumb {
		Error404(w, r)
		return
	}
	if isThumb {
		f, err = os.Open(path.Join(media, name))
		if err != nil {
			pattern := regexp.MustCompile("(.*)_thumb\\.(jpe?g|png|gif)")
			fname := pattern.ReplaceAllString(name, "${1}.${2}")
			f, err = os.Open(path.Join(media, fname))
			if err != nil {
				Error404(w, r)
				return
			}
			img, _, err := image.Decode(f)
			if err != nil {
				panic(err)
			}
			newImage := resize.Resize(160, 0, img, resize.Lanczos3)
			if err = jpeg.Encode(buffer, newImage, nil); err != nil {
				panic(err)
			}
			f, _ = os.OpenFile(path.Join(media, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
			bufferBytes = buffer.Bytes()
			if _, err := io.Copy(f, buffer); err != nil {
				panic(err)
			}
		}
	}

	if len(bufferBytes) == 0 {
		if _, err := io.Copy(buffer, f); err != nil {
			panic(err)
		}
		bufferBytes = buffer.Bytes()
	}

	itoa := strconv.Itoa(len(bufferBytes))
	w.Header().Set("Content-Type", http.DetectContentType(bufferBytes))
	w.Header().Set("Content-Length", itoa)
	_, err = w.Write(bufferBytes)
	logger.HandleError(err)
	logger.HandleError(f.Close())
}

func InternalServerError(w http.ResponseWriter, r *http.Request, rec interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write([]byte("500 Internal Server Error"))
	logger.Logger.Errorf("500 - IP:%s - %s: %s%s - %v", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path, rec)
	logger.HandleError(e)
}

func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Page not found"))
	logger.Logger.Warningf("404 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	logger.HandleError(err)
}

func Error405(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte("Method not allowed"))
	logger.Logger.Warningf("405 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	logger.HandleError(err)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	logger.Logger.Warningf("403 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	_, err := w.Write([]byte("Forbidden"))
	logger.HandleError(err)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	logger.Logger.Warningf("401 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	if message == "" {
		message = "Unauthorized"
	}
	_, err := w.Write([]byte(message))
	logger.HandleError(err)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.Logger.Warningf("IP:%s - %s: %s%s - %s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path, message)
	JSONResponse(w, types.H{
		"error":   true,
		"message": message,
	})
}
