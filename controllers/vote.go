package controllers

import (
	"EDProyecto/commons"
	"EDProyecto/configuration"
	"EDProyecto/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//archivo que controla la funcionalidad de la estructura vote

//Controlador para registrar un voto
func VoteRegister(w http.ResponseWriter, r *http.Request) {
	vote := models.Vote{}
	user := models.User{}
	currentVote := models.Vote{}
	m := models.Message{}

	//asignamos el usuario del voto a nuestra variable user
	user, _ = r.Context().Value("user").(models.User)

	//traemos el voto de la request y controlamos si hay errores
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		m.Message = fmt.Sprintf("Error al leer el voto a registrar: %s", err)
		m.Code = http.StatusBadRequest
		commons.DisplayMessage(w, m)
		return
	}

	//asignamos al voto el ID del user
	vote.UserID = user.ID

	//me conecto a la db
	db := configuration.GetConection()
	defer db.Close()

	//cheque en la db si el usuario voto o no
	db.Where("comment_id = ? and user_id = ?", vote.CommentID, vote.UserID).First(&currentVote)

	//si no realizo ningun voto tambien debo registrarlo
	if currentVote.ID == 0 {
		db.Create(&vote)
		err := updateCommentVotes(vote.CommentID, vote.Value, false)
		if err != nil {
			m.Message = err.Error()
			m.Code = http.StatusBadRequest
			commons.DisplayMessage(w, m)
			return
		}
		//en el caso de que no haya saltado ningun error, devuelvo un mensaje de OK
		m.Message = "Voto registrado"
		m.Code = http.StatusCreated
		commons.DisplayMessage(w, m)
		return
	} else if currentVote.Value != vote.Value {
		//si voto positivo y ahora quiere votar negativo (o viceversa), se procede lo sgte
		currentVote.Value = vote.Value
		db.Save(&currentVote)
		err := updateCommentVotes(vote.CommentID, vote.Value, true)
		if err != nil {
			m.Message = err.Error()
			m.Code = http.StatusBadRequest
			commons.DisplayMessage(w, m)
		}
		m.Message = "Voto actualizado"
		m.Code = http.StatusOK
		commons.DisplayMessage(w, m)
		return
	}

	//si el proceso llego a esta instancia quiere decir que el usuario ya tenia un voto registrado
	//y no puede volver a votar
	m.Message = "Este voto ya está registrado"
	m.Code = http.StatusBadRequest
	commons.DisplayMessage(w, m)
}

//esta funcionalidad actualiza el nro de votos que tiene un comentario
//isUpdate indica si es un voto para actualizar o no
func updateCommentVotes(commentID uint, vote bool, isUpdate bool) (err error) {
	comment := models.Comment{}

	//me conecto a la db
	db := configuration.GetConection()
	defer db.Close()

	//corroboro que el user haya realizado algun voto
	rows := db.First(&comment, commentID).RowsAffected
	if rows > 0 {
		if vote {
			comment.Votes++
			if isUpdate {
				comment.Votes++
			}
		} else {
			comment.Votes--
			if isUpdate {
				comment.Votes--
			}
		}
		db.Save(&comment)
	} else {
		err = errors.New("No se encontró un registro de comentario para asignarle el voto")
	}
	return
}
