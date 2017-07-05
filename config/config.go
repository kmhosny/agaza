package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"gopkg.in/yaml.v2"
)

var config *Configuration
var once sync.Once

const (
	configurationPath        = "config.yml"
	exampleConfigurationPath = "config.yml.example"
)

type Configuration struct {
	Dbnumber                    int           `yaml:"dbnumber"`
	Redisport                   int           `yaml:"redisport"`
	Redishost                   string        `yaml:"redishost"`
	Redisprotocol               string        `yaml:"redisprotocol"`
	Redispoolsize               int           `yaml:"redispoolsize"`
	APIserverport               string        `yaml:"APIserverport"`
	Environment                 string        `yaml:"enviroment"`
}

func LoadConfiguration() *Configuration {
	once.Do(func() {
		configFile, err := ioutil.ReadFile(configurationPath)
		if err != nil {
			panic("Cannot load configration from config.yml, please run `cp config.yml.example config.yml` to load default configration")
		} else {
			env := GetENV()
			var allConfigs map[string]*Configuration
			if err := yaml.Unmarshal([]byte(configFile), &allConfigs); err != nil {
				panic("Cannot load configration with error" + err.Error())
			}
			if val, ok := allConfigs[env]; ok {
				config = val
				validateConfigrations()
			} else {
				panic(fmt.Sprintf("Cannot find configration for %s environment, the available envirnoments: %s", env, reflect.ValueOf(allConfigs).MapKeys()))
			}
		}
	})
	return config
}

func GetENV() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	return env
}

func validateConfigrations() {
	exConfigFile, exErr := ioutil.ReadFile(exampleConfigurationPath)
	configFile, err := ioutil.ReadFile(configurationPath)
	if err == nil && exErr == nil {
		env := GetENV()
		var exampleConfigs, allConfigs map[string]map[string]interface{}

		exErr2 := yaml.Unmarshal([]byte(exConfigFile), &exampleConfigs)
		err2 := yaml.Unmarshal([]byte(configFile), &allConfigs)

		if err2 == nil && exErr2 == nil {
			if exampleConfigration, ok := exampleConfigs[env]; ok {
				currentConfigration := allConfigs[env]

				diffkeys := difference(currentConfigration, exampleConfigration)
				if len(diffkeys) != 0 {
					error_str :=
						`Boot halted! Application configuration appears to be invalid.
       Reason: config.yml and config.yml.example are out of sync.
       Check the following keys: %s on config.yml.example`
					panic(fmt.Sprintf(error_str, diffkeys))
				}
			}
		}
	}
}

func difference(currentConfigration map[string]interface{}, exampleConfigration map[string]interface{}) []string {
	diffStr := []string{}
	set := map[string]bool{}

	for k, _ := range currentConfigration {
		set[k] = true
	}

	for k2, _ := range exampleConfigration {
		if _, ok := set[k2]; !ok {
			diffStr = append(diffStr, k2)
		}
	}

	return diffStr
}
