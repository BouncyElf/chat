package models

import (
	"github.com/aofei/air"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitModel() {
	cfg, ok := air.Config["DB"].(map[string]interface{})
	if !ok {
		panic("no database config found in config file")
	}
	username, ok := cfg["username"].(string)
	if !ok {
		panic("no database username found in config file")
	}
	password, ok := cfg["password"].(string)
	if !ok {
		panic("no database password found in config file")
	}
	addr, ok := cfg["addr"].(string)
	if !ok {
		panic("no database address found in config file")
	}
	DBName, ok := cfg["dbname"].(string)
	if !ok {
		panic("no database name found in config file")
	}
	var err error
	DB, err = gorm.Open("mysql", username+":"+password+
		"@tcp("+addr+")/"+DBName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	initMessage()
}
