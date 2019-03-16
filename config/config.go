package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

var BaseDir, _ = os.Getwd()
var ConfPath = path.Join(BaseDir, "./config/config.yml")

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
		Email    string `yaml:"EMAIL"`
		Host     string `yaml:"HOST"`
		Port     string `yaml:"PORT"`
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
	yaml.Unmarshal(data, &Config)
}
