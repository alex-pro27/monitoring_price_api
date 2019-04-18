package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var ConfPath = os.Getenv("MONITORING_PRICE_CONF")

type TypeConfig struct {
	Databases struct {
		Default struct {
			User     string `yaml:"USER"`
			Database string `yaml:"DATABASE"`
			Port     string `yaml:"PORT"`
			Password string `yaml:"PASSWORD"`
			Host     string `yaml:"HOST"`
		}
	}

	System struct {
		SecretKey string `yaml:"SECRET_KEY"`
		Debug     bool   `yaml:"DEBUG"`
		Server    string `yaml:"SERVER"`
		LogPath   string `yaml:"LOG_PATH"`
	}

	Session struct {
		Key    string `yaml:"KEY"`
		MaxAge int    `yaml:"MAX_AGE"`
	}

	Admin struct {
		Name  string `yaml:"NAME"`
		Email string `yaml:"EMAIL"`
	}

	Email struct {
		Name     string `yaml:"NAME"`
		From     string `yaml:"FROM"`
		Host     string `yaml:"HOST"`
		Port     int    `yaml:"PORT"`
		User     string `yaml:"USER"`
		Password string `yaml:"PASSWORD"`
	}

	Static struct {
		StaticRoot string `yaml:"STATIC_ROOT"`
		MediaRoot  string `yaml:"MEDIA_ROOT"`
	}
}

var Config *TypeConfig

func Load() {
	data, _ := ioutil.ReadFile(ConfPath)
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatal(err)
	}

}
