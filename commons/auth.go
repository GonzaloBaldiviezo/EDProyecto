package commons

//este archivo se ejecuta cada vez que queramos loguearnos en la api
import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	jwt "github.com/dgrijalva/jwt-go"
	//"EDProyecto\EDProyecto\models"
)

var(
	privateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
)

func init()  {
	privateBytes, err := ioutil.ReadFile("./keys/private.rsa")

	if err != nil {
		log.Fatal("No se pudo leer el archivo privado")
	}

	publicBytes, err := ioutil.ReadFile("./keys/public.rsa")

	if err != nil {
		log.Fatal("No se pudo leer el archivo publico")
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)

	if err != nil {
		log.Fatal("No se pudo hacer el parse a privateKey")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)

	if err != nil {
		log.Fatal("No se pudo hacer el parse a publicKey")
	}

	func GenerateJWT(user models.User) string {
		claims := models.Claim{
			User: user,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
				Issuer: "Escuela Digital",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		result, err := token.Signedtring(privateKey)

		if err != nil {
			log.Fatal("No se pudo firmar el token")
		}

		return result
	}
}