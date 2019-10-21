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
	Path		string
}

//Init gets value of each Config variables from `yaml` configuration file
func (c *Config) Init(envPath string) error {
	fmt.Printf("%s",envPath)
	if data, err := ioutil.ReadFile(envPath); err == nil {
		if err := yaml.UnmarshalStrict(data, c); err != nil {
			fmt.Println("Unmarshal error:" + err.Error())
		}
		return err
	}
	
	err := fmt.Errorf("Error reading configuration file: %s\n", envPath)
	fmt.Println(err.Error())
	
	return err
}
