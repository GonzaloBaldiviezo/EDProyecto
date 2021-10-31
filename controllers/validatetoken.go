package controllers

import (
	"EDProyecto/commons"
	"EDProyecto/models"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
)

//funcionalidad para validar el token que le permite al usuario votar y comentar
func ValidateToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var m models.Message

	//recibimos el token y lo validamos. Este metodo puede devolver varios errores
	token, err := request.ParseFromRequest(
		r,
		request.OAuth2Extractor,
		//&models.Claim{},
		func(t *jwt.Token) (interface{}, error) {
			return commons.PublicKey, nil
		},
		request.WithClaims(&models.Claim{}),
	)

	//controlamos los errores
	if err != nil {
		m.Code = http.StatusUnauthorized
		switch err.(type) {
		case *jwt.ValidationError:
			vError := err.(*jwt.ValidationError)
			switch vError.Errors {
			case jwt.ValidationErrorExpired:
				m.Message = "Su token ha expirado"
				commons.DisplayMessage(w, m)
				return
			case jwt.ValidationErrorSignatureInvalid:
				m.Message = "La firma del token no coincide"
				commons.DisplayMessage(w, m)
				return
			default:
				m.Message = "Su token no es válido"
				commons.DisplayMessage(w, m)
				return
			}
		}
	}

	if token.Valid {
		ctx := context.WithValue(r.Context(), "user", token.Claims.(*models.Claim).User)
		next(w, r.WithContext(ctx))
	} else {
		m.Code = http.StatusUnauthorized
		m.Message = "Su token no es válido"
		commons.DisplayMessage(w, m)
	}
}
