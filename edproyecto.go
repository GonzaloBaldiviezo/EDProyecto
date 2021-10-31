package main

import (
	"EDProyecto/migration"
	"EDProyecto/routes"
	"flag"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

//Este archivo es el ejecutable. En este caso, vamos a crear la logica para lograr que
// la funcion migrate se ejectue una sola vez. Para eso vamos a definir que, si en la terminal escribimos
// "edproyecto.exe -- migrate yes", entonces se va a ejecutar dicho metodo. Caso contrario, si solo ponemos
// "edproyecto.exe" no se ejecutara el metodo (suponiendo que las tablas ya estan creadas)
func main() {
	//creo la variable migrate. Esta variable va a alternar entre yes o no
	var migrate string

	//uso un flag para definir el comando de consola migrate y sus parametros
	flag.StringVar(&migrate, "migrate", "no", "Genera la migracion a la DB")
	flag.Parse()

	//si la variable migrate tiene el valor yes, entonces que comience la migracion
	if migrate == "yes" {
		log.Println("Comenzó la migracion")
		migration.Migrate()
		log.Println("Finalizó la migracion")
	}

	//inicia las rutas
	router := routes.InitRoutes()

	//inicia los middlewares
	n := negroni.Classic()
	n.UseHandler(router)

	//inicia el servidor
	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}

	log.Println("Iniciado el servidor en http://localhost:8080")
	log.Println(server.ListenAndServe())
	log.Println("Fin de la ejecucion del programa")
}
