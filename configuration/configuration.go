package configuration

import (
	"encoding/json"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Configuration struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

func GetConfiguration() Configuration {
	var c Configuration
	file, err := os.Open("./config.json")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&c)

	if err != nil {
		log.Fatal(err)
	}

	return c
}

func GetConection() *gorm.DB {
	c := GetConfiguration()
	user:password@tcp(server:port)/database?charset=utf8&parseTime=True&loc=Local
}

