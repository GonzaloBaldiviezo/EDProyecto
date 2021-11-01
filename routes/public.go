package routes

import (
	"github.com/gorilla/mux"
	"net/http"
)

//expone los archivos estaticos (html, css, js, etc)
func SetPublicRouter(router *mux.Router) {
	router.Handle("/", http.FileServer(http.Dir("./public")))
}
