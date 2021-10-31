package controllers

import (
	"EDProyecto/commons"
	"EDProyecto/configuration"
	"EDProyecto/models"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Controller de user para controlar las funcionalidades

//Login es el controlador de login. cumple la funcion de autorizar o no al usuario
func Login(w http.ResponseWriter, r *http.Request) {
	//creo un objeto de tipo user y le asigno los valores que vienen en la peticion de login
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		fmt.Fprintf(w, "Error%s\n", err)
		return
	}

	//me conecto a la base de datos
	db := configuration.GetConection()
	defer db.Close()

	//encripto la password
	c := sha256.Sum256([]byte(user.Password))

	pwd := fmt.Sprintf("%x", c)

	//compruebo que los datos en la db coinciden con los datos ingresados
	db.Where("email = ? and password = ?", user.Email, pwd).First(&user)

	//si el usuario existe, genero el token
	if user.ID > 0 {
		user.Password = ""
		token := commons.GenerateJWT(user)

		j, err := json.Marshal(models.Token{Token: token})
		if err != nil {
			log.Fatalf("Error al convertir el token a json: %s", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		//si no existe, arrojo un codigo de no autorizado
		m := models.Message{
			Message: "Usuario o clave no válido",
			Code:    http.StatusUnauthorized,
		}
		commons.DisplayMessage(w, m)
	}

}

//esta funcion sirve para registrar usuarios
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	m := models.Message{}

	//si hay un error en el json de los datos a registrar, que lance un error
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		m.Message = fmt.Sprintf("Error al leer el usuario a registrar: %s", err)
		m.Code = http.StatusBadRequest
		commons.DisplayMessage(w, m)
		return
	}

	//si las contraseñas no coinciden, que lance un error
	if user.Password != user.ConfirmPassword {
		m.Message = "Las contraseñas no coinciden"
		m.Code = http.StatusBadRequest
		commons.DisplayMessage(w, m)
	}

	c := sha256.Sum256([]byte(user.Password))

	pwd := fmt.Sprintf("%x", c)

	user.Password = pwd

	//para la imagen de usuario vamos a usar gravatar. Para ello debemos codificar el email a md5
	picmd5 := md5.Sum([]byte(user.Email))
	picstn := fmt.Sprintf("%x", picmd5)

	//generamos la url a gravatar con el size = 100
	user.Picture = "https://gravatar.com/avatar/" + picstn + "?s=100"

	//nos conectamos a la db para guardar el user
	db := configuration.GetConection()
	defer db.Close()

	err = db.Create(&user).Error
	if err != nil {
		m.Message = fmt.Sprintf("Error al crear el registro: %s", err)
		m.Code = http.StatusBadRequest
		commons.DisplayMessage(w, m)
		return
	}

	m.Message = "Usuario creado con éxito"
	m.Code = http.StatusCreated
	commons.DisplayMessage(w, m)
}
