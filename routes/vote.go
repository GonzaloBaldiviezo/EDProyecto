package routes

import (
	"EDProyecto/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//rutas para el controlador de votos

//ruta para el registro de un voto
func SetVoteRouter(router *mux.Router) {
	prefix := "/api/votes"
	subRouter := mux.NewRouter().PathPrefix(prefix).Subrouter().StrictSlash(true)
	subRouter.HandleFunc("/", controllers.VoteRegister).Methods("POST")

	router.PathPrefix(prefix).Handler(
		negroni.New(
			negroni.HandlerFunc(controllers.ValidateToken),
			negroni.Wrap(subRouter),
		),
	)
}
