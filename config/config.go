package config

import (
	"github.com/alex-pro27/monitoring_price_api/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)


var BaseDir, _ = os.Getwd()
var ConfPath = path.Join(BaseDir, "./config/config.yml")

type TypeConfig struct {
	Databases struct {
		Default struct{
			User 		string 	`yaml:"USER"`
			Database 	string 	`yaml:"DATABASE"`
			Port 		string	`yaml:"PORT"`
			Password 	string 	`yaml:"PASSWORD"`
			Host 		string	`yaml:"HOST"`
		}
	}

	System struct {
		SecretKey 	string	`yaml:"SECRET_KEY"`
		Debug 		bool	`yaml:"DEBUG"`
		Server 		string	`yaml:"SERVER"`
	}
}

var Config *TypeConfig

func Load() {
	data, err := ioutil.ReadFile(ConfPath)
	common.HandlerError(err)
	err = yaml.Unmarshal(data, &Config)
	common.HandlerError(err)
}