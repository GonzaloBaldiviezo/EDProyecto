package configuration

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

//esta estructura cumple la funcion de configurar los datos de nuestra conexion con la db
//convirtiendo los datos para que puedan ser interpretados por gorm
type configuration struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
}

//funcion para obtener la configuracion
func getConfiguration() configuration {
	var c configuration
	//obtenemos el archivo del cual vamos a obtener los datos de conexion
	file, err := os.Open("./config.json")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//lo decodificamos para poder interpretarlo
	err = json.NewDecoder(file).Decode(&c)

	if err != nil {
		log.Fatal(err)
	}

	return c
}

//funcion para realizar la conexion con la db
func GetConection() *gorm.DB {
	//ya obtenida la configuracion, la guardamos en una variable
	c := getConfiguration()
	//dsn = data source name. En esta variable definimos como queda la url
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Password, c.Server, c.Port, c.Database)
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}
	return db
}
