package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Config is this app yaml configuration
type Config struct {
	Host        string	
	Port        string	
	Type        string
	Country 	string
	Name        string
}

//Init gets value of each Config variables from `yaml` configuration file
func (c *Config) Init(env string) error {
	if data, err := ioutil.ReadFile(LoadConfigPath()); err == nil {
		if err := yaml.UnmarshalStrict(data, c); err != nil {
			fmt.Println("Unmarshal error:" + err.Error())
		}
		return err
	}
	return nil
}
