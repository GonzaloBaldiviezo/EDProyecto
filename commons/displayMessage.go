package commons

import (
	"EDProyecto/models"
	"encoding/json"
	"log"
	"net/http"
)

//Este archivo cumplira la funcion de mostrar todos los msjs que se lancen desde la api
func DisplayMessage(w http.ResponseWriter, m models.Message) {
	j, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error al convertir el mensaje %s", err)
	}

	w.WriteHeader(m.Code)
	w.Write(j)
}
