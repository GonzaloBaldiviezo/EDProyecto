package commons

//este archivo se ejecuta cada vez que queramos loguearnos en la api
import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"EDProyecto/models"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	privateKey *rsa.PrivateKey

	//PublicKey es exportable
	PublicKey *rsa.PublicKey
)

func init() {

	//leemos los archivos private y public
	privateBytes, err := ioutil.ReadFile("./keys/private.rsa")

	if err != nil {
		log.Fatal("No se pudo leer el archivo privado")
	}

	publicBytes, err := ioutil.ReadFile("./keys/public.rsa")

	if err != nil {
		log.Fatal("No se pudo leer el archivo publico")

	}

	//parseamos ambos archivos y se los asociamos a las variables creadas mas arriba
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se pudo hacer el parse a privateKey")
	}

	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se pudo hacer el parse a PublicKey")
	}
}

//metodo exportable que firma el token
func GenerateJWT(user models.User) string {
	claims := models.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			//tiempo de expiracion
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "Escuela Digital",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	result, err := token.SignedString(privateKey)

	if err != nil {
		log.Fatal("No se pudo firmar el token")
	}

	return result
}
