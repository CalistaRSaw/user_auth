package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/CalistaRSaw/uauth-vbasic/database"
	"github.com/CalistaRSaw/uauth-vbasic/models"
	"github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	Email    string
	Name     string
	Uid      string
	Category string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET")

func GenerateAllTokens(email string, name string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {

	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		Uid:      uid,
		Category: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			Issuer:    "bookite.auth.service",
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			Issuer: "bookite.auth.service",
		},
	}
	// generate jwt
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// Sign and get the complete encoded token as a string using the secret
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err

}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var db = database.DBConn
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	defer cancel()

	var user models.User
	db.WithContext(ctx).Find(&user, "user_id = ?", userId)

	result := db.Model(&user).Updates(models.User{Token: signedToken, Refresh_token: signedRefreshToken, Updated_at: Updated_at})

	if result.Error != nil {
		log.Panic(result.Error)
		return
	}
}

func RefreshToken(email string, name string, userType string, uid string) {
	var db = database.DBConn
	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		Uid:      uid,
		Category: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			Issuer:    "bookite.auth.service",
		},
	}
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	var user models.User
	db.Find(&user, "user_id = ?", uid)

	result := db.Model(&user).Updates(models.User{Token: token, Updated_at: Updated_at})

	if result.Error != nil {
		log.Panic(result.Error)
		return
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	// decode & validate
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = "the token is invalid"
		msg = err.Error()
		return
	}

	// check expiration
	if time.Now().Local().Unix() > claims.ExpiresAt {
		msg = "token is expired"
	}
	return claims, msg
}
