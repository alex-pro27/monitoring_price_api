package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/nfnt/resize"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ResizeImage(imageFile io.Reader, name, mediaPath string, width, height uint, bytesArray *[]byte) error {
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	newImage := resize.Resize(width, height, img, resize.NearestNeighbor)
	if err = jpeg.Encode(buffer, newImage, nil); err != nil {
		return err
	}
	r := regexp.MustCompile("^(.*)\\.(?:jpe?g|png|gif)")
	name = r.ReplaceAllString(name, "${1}.jpg")
	newFile, _ := os.OpenFile(path.Join(mediaPath, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer func() {
		logger.HandleError(newFile.Close())
	}()
	if bytesArray != nil {
		*bytesArray = append(*bytesArray, buffer.Bytes()...)
	}
	if _, err := io.Copy(newFile, buffer); err != nil {
		return err
	}
	return nil
}

func CreateXLSX(header []string, items [][]string, fileName string) (filePath string, err error) {
	if fileName == "" {
		fileName = string(time.Now().Unix())
	}
	filePath = path.Join(os.TempDir(), fmt.Sprintf("%s.xlsx", fileName))
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}
	row := sheet.AddRow()
	headerStyle := xlsx.Style{
		Font: xlsx.Font{
			Bold: true,
		},
	}
	for _, title := range header {
		cell := row.AddCell()
		cell.SetStyle(&headerStyle)
		cell.Value = title
	}
	for _, rowData := range items {
		row := sheet.AddRow()
		for _, value := range rowData {
			cell := row.AddCell()
			cell.Value = value
		}
	}
	if err = file.Save(filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

func GenerateHash() string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(strconv.FormatInt(time.Now().Unix(), 10)),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal(err)
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetIPAddress(r *http.Request) string {
	var ipAddress string
	ipAddress = strings.Split(r.RemoteAddr, ":")[0]
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			} else {
				ipAddress = ip
				goto Done
			}
		}
	}
Done:
	return ipAddress
}
