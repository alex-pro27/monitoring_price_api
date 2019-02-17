package config

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)


var BaseDir, _ = os.Getwd()
var ConfPath = "/home/alex/projects/go_monitoring_price_api/config/config.yml"

type TypeConfig struct {
	Databases struct {
		Default struct{
			User string 	`yaml:"USER"`
			Database string `yaml:"DATABASE"`
			Port string		`yaml:"PORT"`
			Password string `yaml:"PASSWORD"`
			Host string		`yaml:"HOST"`
		}
	}

	System struct {
		Debug bool		`yaml:"DEBUG"`
		Server string	`yaml:"SERVER"`
	}
}

var Config TypeConfig

func Load() {
	fmt.Println(BaseDir)
	data, err := ioutil.ReadFile(ConfPath)
	common.HandlerError(err)
	err = yaml.Unmarshal(data, &Config)
	common.HandlerError(err)
}