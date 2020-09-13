package authentication

import (
	"Golangecho/sessions"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"time"
)

// struct that will be encoded to a JWT.
type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	name := c.FormValue("name")
	pass := c.FormValue("pass")
	userID := 1

	//Verify user exist in database
	if name == "jack" && pass == "1234" {

		token, err := createJWTToken(name)
		if err != nil {
			log.Println("failed creating token", err)
			return c.String(http.StatusInternalServerError, "error")
		}

		c.SetCookie(&http.Cookie{
			Name:  "JWTCookie",
			Value: token,
			Expires: time.Now().Add(time.Hour * 10),
		})

		//set token in redis
		sessions.SessionsStore.Set(token, sessions.Session{
			Name:       name,
			UserID:     userID,
		})

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You are logged in",
			"token": token,
		})
	}
	return c.String(http.StatusUnauthorized, "Not registered")
}


func createJWTToken(name string) (string, error) {
	claims := JwtClaims{
		Name:           name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(),
		},
	}
	//create token
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	//hash token
	token, err := rawToken.SignedString([]byte("secret key"))
	if err != nil {
		return "", err
	}
	return token, nil
}
