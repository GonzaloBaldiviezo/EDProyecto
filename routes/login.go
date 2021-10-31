package routes

import (
	"EDProyecto/controllers"
	"github.com/gorilla/mux"
)

//ruta para hacer el login

//funcion que hace el router para login
func SetLoginRouter(router *mux.Router) {
	router.HandleFunc("/api/login", controllers.Login).Methods("POST")
}
