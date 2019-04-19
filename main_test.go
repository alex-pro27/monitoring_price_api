package main

import (
	"context"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"testing"
	"time"
)

var client *http.Client

func setUp() {
	Init()
	go StartServer()
	cookieJar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: cookieJar,
	}
	time.Sleep(2 * time.Second)
}

func tearDown() {
	logger.HandleError(Server.Shutdown(context.TODO()))
}

func TestMain(m *testing.M) {
	setUp()
	m.Run()
	tearDown()
}

func readBody(resp *http.Response) (body string, err error) {
	b, _ := ioutil.ReadAll(resp.Body)
	body = string(b)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("body: %s; status: %d", body, resp.StatusCode)
	}
	return body, err
}

func getUrl(path string) string {
	return "http://" + config.Config.System.Server + path
}

func authAdmin() (err error) {
	resp, err := client.PostForm(getUrl("/api/admin/login"), url.Values{
		"username": {"alex"},
		"password": {"123123"},
	})
	if err != nil {
		return err
	}
	_, err = readBody(resp)
	return err
}

func allUser() (err error) {
	resp, err := client.Get(getUrl("/api/admin/users"))
	if err != nil {
		return err
	}
	_, err = readBody(resp)
	return err
}

func getUser() (err error) {
	resp, err := client.Get(getUrl("/api/admin/user/1"))
	if err != nil {
		return err
	}
	_, err = readBody(resp)
	return err
}

func TestAdmin(t *testing.T) {
	var err error
	if err = authAdmin(); err != nil {
		t.Error(err)
	} else {
		parallel := []func() error{
			allUser,
			getUser,
		}
		for _, f := range parallel {
			name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			name = regexp.MustCompile("\\.([\\w\\d]+)$").FindStringSubmatch(name)[1]
			t.Run(name, func(t *testing.T) {
				if err := f(); err != nil {
					t.Error(err)
				}
			})
		}
	}
}
