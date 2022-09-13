package config

import (
	"fmt"
	"github.com/DenrianWeiss/avalium/model"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

type Config struct {
	Rpc   model.RPCList `json:"rpc" mapstructure:"rpc"`
	Serve Serve         `json:"serve" mapstructure:"serve"`
	Debug bool          `json:"debug" mapstructure:"debug"`
}

type Serve struct {
	ServerAddr   string `json:"server_addr" mapstructure:"server_addr"`
	ControlPlane string `json:"control_plane" mapstructure:"control_plane"`
}

var config Config

func Init() {
	conf := Config{}
	dir, _ := filepath.Abs(`.`)
	// load system
	fmt.Println(dir + `/config`)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(dir + `/config`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
	log.Println("config:", conf)
	config = conf
}

func IsDebug() bool {
	return config.Debug
}

func GetConfig() Config {
	return config
}
