package models

import (
	"github.com/jinzhu/gorm"
)

//clase user
type User struct {
	gorm.Model
	//como vamos a hacer una api, tenemos que definir que vamos a devolver un json
	Username string `json:"username" gorm:"not null;unique"`
	Email    string `json:"email" gorm:"not null;unique"`
	Fullname string `json:"fullname" gorm:"not null"`
	//el comando omitempty sirve para que cuando consultemos datos no nos devuelva el password
	Password        string    `json:"password,omitempty" gorm:"not null;type:varchar(256)"`
	ConfirmPassword string    `json:"confirmPassword,omitempty" gorm:"-"`
	Picture         string    `json:"picture"`
	Comments        []Comment `json:"commnts,omitempty"`
}
