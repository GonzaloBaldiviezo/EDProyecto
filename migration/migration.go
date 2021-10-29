package migration

import (
	"EDProyecto/configuration"
	"EDProyecto/models"
)

//Este archivo tiene el objetivo de crear las tablas que vamos a usar y la idea es que solo
// se ejecute una sola vez

//funcion que crea las tablas
func Migrate() {
	//conectamos a la db
	db := configuration.GetConection()
	defer db.Close()

	//creamos las tablas
	db.CreateTable(&models.User{})
	db.CreateTable(&models.Comment{})
	db.CreateTable(&models.Vote{})
	//en este caso, en la estructura vote solo un user puede votar una sola vez, y para hacer eso
	//debemos hacer que el user y el comment sean unicos combinados
	db.Model(&models.Vote{}).AddUniqueIndex("comment_id_user_id_unique", "comment_id", "user_id")
}
