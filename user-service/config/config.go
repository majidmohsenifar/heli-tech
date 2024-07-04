package config

import (
	"flag"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const FileName = "config"
const FilePath = "./config/"

// since viper does not accept absolute path  we have this relative path for our functional test
const TestConfigFilePath = "./../../config/"

func NewViper() *viper.Viper {
	v := viper.New()

	v.SetConfigType("yml")
	v.SetConfigName(FileName)

	if !isTestEnv() {
		v.AddConfigPath(FilePath)
	} else {
		v.AddConfigPath(TestConfigFilePath)
	}

	err := v.ReadInConfig()
	if err != nil {
		panic("can not read config file" + err.Error())
	}

	//reading from .env if exists
	_, err = os.Stat("./.env")
	if err == nil {
		v.SetConfigType("env")
		v.SetConfigName("./../.env")
		err = v.MergeInConfig()
		if err != nil {
			panic("can not merge config file" + err.Error())
		}
	}
	v.AutomaticEnv()
	v.SetEnvPrefix("heli")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//rewrite configs for test env
	//if isTestEnv() {
	//v.SetConfigName(TestConfigFileName)
	//err = v.MergeInConfig()
	//if err != nil {
	//panic("can not merge config file" + err.Error())
	//}
	//}

	_ = v.AllSettings()
	return v
}

func isTestEnv() bool {
	return flag.Lookup("test.v") != nil
}
