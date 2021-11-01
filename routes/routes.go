package routes

import "github.com/gorilla/mux"

//este archivo va a ser el centro de todas nuestras rutas

//inicia las rutas
func InitRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	SetLoginRouter(router)
	SetUserRouter(router)
	SetCommentRouter(router)
	SetVoteRouter(router)
	SetRealtimeRouter(router)
	SetPublicRouter(router)

	return router
}
