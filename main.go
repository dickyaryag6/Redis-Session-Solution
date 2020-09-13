package main

import (
	"Golangecho/authentication"
	"Golangecho/sessions"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func hello(c echo.Context) error {
	cookie, err := c.Cookie("JWTCookie")
	//if cookie doesnt exist
	if err != nil {
		c.String(http.StatusUnauthorized, "Not Authorized")
		return err
	}
	//	get jwt token
	tokenString := cookie.Value
	//get session data from redis
	session, err := sessions.SessionsStore.Get(tokenString)

	return c.String(http.StatusOK, fmt.Sprintf("Hello %s, your are authorized", session.Name))
}


func main() {
	e := echo.New()
	//sessionsStore = NewMemoryStore()
	sessions.SessionsStore = sessions.NewRedisStore()

	e.GET("/login", authentication.Login)

	g := e.Group("")

	//Use Middleware to verify JWT Token
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey: []byte("secret key"),
		TokenLookup: "cookie:JWTCookie",
	}))


	g.GET("/hello", hello)


	e.Start(":8000")
}
