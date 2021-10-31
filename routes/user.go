package routes

import (
	"EDProyecto/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//este archivo enruta las urls de nuestra api

//ruta para el registro de un usuario
func SetUserRouter(router *mux.Router) {
	prefix := "/api/users"
	subRouter := mux.NewRouter().PathPrefix(prefix).Subrouter().StrictSlash(true)
	subRouter.HandleFunc("/", controllers.CreateUser).Methods("POST")

	router.PathPrefix(prefix).Handler(
		negroni.New(
			negroni.Wrap(subRouter),
		),
	)
}
