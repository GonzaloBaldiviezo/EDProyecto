package controllers

import (
	"EDProyecto/commons"
	"EDProyecto/configuration"
	"EDProyecto/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// funcionalidad para crear comentarios
func CommentCreate(w http.ResponseWriter, r *http.Request) {

	//creo las variables a utilizar
	comment := models.Comment{}
	user := models.User{}
	m := models.Message{}

	//traigo el usuario
	user, _ = r.Context().Value("user").(models.User)

	//verifico que el body de la solicitud este correcta, sino, controlo el error
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		m.Code = http.StatusBadRequest
		m.Message = fmt.Sprintf("Error al leer el comentario: %s", err)
		commons.DisplayMessage(w, m)
		return
	}

	//le asigno el usuario al comentario
	comment.UserID = user.ID

	//me conecto a la db
	db := configuration.GetConection()
	defer db.Close()

	//creo el comentario y controlo el posible error que me genere
	err = db.Create(&comment).Error
	if err != nil {
		m.Code = http.StatusBadRequest
		m.Message = fmt.Sprintf("Error al registrar el comentario: %s", err)
		commons.DisplayMessage(w, m)
		return
	}

	//lanzo el mensaje de OK
	m.Code = http.StatusCreated
	m.Message = "Comentario creado con éxito"
	commons.DisplayMessage(w, m)
}

//funcionalidad para obtener todos los comentarios
func CommentGetAll(w http.ResponseWriter, r *http.Request) {
	//creo las variables a utilizar

	//slice de todos los comentarios que voy a traer
	comments := []models.Comment{}

	//mensaje de respuesta de todos los posibles casos
	m := models.Message{}

	//usuario del comentario
	user := models.User{}

	//votos del comentario
	vote := models.Vote{}

	//asigno el usuario que viene en el request a mi variable user (para poder identificar si voto o no)
	user, _ = r.Context().Value("user").(models.User)

	//cuando enviamos una solicitud para traer una lista de datos (comentarios), podemos dartle cierta
	//forma a dicha lista, como por ej. ordenarlos por ID o user y tambien mostrar una cantidad
	//limitada por pagina. Lo que hacemos a continuacion es interpretar esa query que viene en
	//el request para que podamos hacer la logica en codigo y devolver lo solicitado

	//creamos una variable que obtenga la query de la request
	vars := r.URL.Query()

	//me conecto a la base de datos
	db := configuration.GetConection()
	defer db.Close()

	//preparo la consulta

	//busco los comentarios padre
	cComment := db.Where("parent_id = 0")

	//si en la consulta los quieren ordenar por votos, planteo el sgte codigo
	if order, ok := vars["order"]; ok {
		if order[0] == "votes" {
			cComment = cComment.Order("votes desc, created_at desc")
		}

		//sino, tambien pueden estar limitados por una cant determinada por pagina
	} else {
		if idlimit, ok := vars["idlimit"]; ok {
			registerByPage := 30
			offset, err := strconv.Atoi(idlimit[0])
			if err != nil {
				log.Println("Error:", err)
			}
			cComment = cComment.Where("id BETWEEN ? AND ?", (offset - registerByPage), offset)
		}
		//ademas los ordeno por id
		cComment = cComment.Order("id desc")
	}

	//todos los datos de cComments los almaceno en comments
	cComment.Find(&comments)

	//recorremos el slice de comments para agregarle a cada comentario el usuario que realizo ese comentario
	for i := range comments {
		//con la sgte linea agregamos el usuario pero viene con su clave
		db.Model(&comments[i]).Related(&comments[i].User)

		//tambien traigo los comentarios hijos con sus usuarios
		comments[i].Children = commentGetChildren(comments[i].ID)

		//con esta linea borramos la clave de cada usuario
		comments[i].User[0].Password = ""

		//se busca el voto del usuario en sesion
		vote.CommentID = comments[i].ID
		vote.UserID = user.ID
		//con la variable count verifico si contabilizo el voto en el comentario
		count := db.Where(&vote).Find(&vote).RowsAffected
		if count > 0 {
			//analizo si el voto es positivo o negativo
			if vote.Value {
				comments[i].HasVote = 1
			} else {
				comments[i].HasVote = -1
			}
		}
	}

	//creamos uana variable para convertir todos los comentarios a un formato json
	j, err := json.Marshal(comments)
	if err != nil {
		m.Code = http.StatusInternalServerError
		m.Message = "Error al convertir los comentarios a json"
		commons.DisplayMessage(w, m)
		return
	}

	if len(comments) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		m.Code = http.StatusNoContent
		m.Message = "No se encontraron comentarios"
		commons.DisplayMessage(w, m)
	}
}

//funcionalidad para traer todos los comentarios hijos (subcomentarios)
func commentGetChildren(id uint) (children []models.Comment) {
	//nos conectamos a la db
	db := configuration.GetConection()
	defer db.Close()

	//obtengo de la db todos los comentarios con un valor de parentID que sea igual al valor id
	//que paso por parametro y los deposito en el slice children que voy a retornar
	db.Where("parent_id = ?", id).Find(&children)

	//recorro el slice children para asociar a cada comment su usuario
	for i := range children {
		//con la sgte linea agregamos el usuario pero viene con su clave
		db.Model(&children[i]).Related(&children[i].User)

		//con esta linea borramos la clave de cada usuario
		children[i].User[0].Password = ""
	}
	return
}
