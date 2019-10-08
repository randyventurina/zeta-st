package main

import (
	"os"
)

//LoadEnv returns the current env
func LoadEnv() string {
	if value, exists := os.LookupEnv("ENV"); exists {
		return value
	}
	return ""

}

//LoadConfigPath returns the current env
func LoadConfigPath() string {
	if value, exists := os.LookupEnv("CONFIG_PATH"); exists {
		return value
	}
	return ""

}
