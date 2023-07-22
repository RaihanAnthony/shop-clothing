package tokens

import (
	"log"
	"os"
	"time"

	" github.com/Raihanpoke/shop-clothing/database"
	jwt "github.com/dgtijalva/jwt-go"
)

type SignedDetails struct{
	Email string
	Frist_Name string
	Last_Name  string
	uid string
	jwt.StandardClaims
}

var UserData *sql.DB = database.UserData("users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string. firstname string, lastname string, uid string)(signedtoken string, signedrefreshtoken string, err error) {
	claims := &SignedDetails{
		Email: email,
		First_Name: firstname,
		Last_Name: lastname,
		uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims:  jwt.StandardClaims{
			ExpiresAt: time.Now().Local().add(time.hour * time.Duration(168)),Unix(),
		},
	}

	token, err:= jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return 
	}	

	return token, refreshtoken, err
}

func ValidateToken(signedtoken string) (claims *SignedDetails, msg string){
	token , err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token)(interface{}, error){
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return 
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token in invalid"
		return
	}

	claims.ExpiresAt < time.Now(),Local(),unix(){
		msg = "token is alredy expired"
		return
	}

	return claims, msg
}