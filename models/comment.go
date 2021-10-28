package models

import (
	"github.com/jinzhu/gorm"
)

//estructura comentario para que el usuario pueda comentar e interactuar
type Comment struct {
	gorm.Model
	//ID del usuario que comenta
	UserID uint `json:"userId"`
	//ID del comentario
	ParentID uint `json:"parentId"`
	//Cantidad de votos que tiene ese comentario
	Votes int32 `json:"votes"`
	//contenido del comentario
	Content string `json:"content"`
	//saber si el usuario voto o no
	HasVote int8 `json:"hasVote" gorm:"-"`
	//usuario que hizo el comentario
	User []User `json:"user,omitempty"`
	//Subcomentarios que se crean a partir de el comentario padre
	Children []Comment `json:"children,omitempty"`
}
