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

//func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		c.Response().Header().Set(echo.HeaderServer, "Server-1.0")
//		return next(c)
//	}
//}


//func Middleware() echo.MiddlewareFunc {
//	return func(handler echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			//	get token from cookie
//			cookie, err := c.Cookie("JWTCookie")
//			//if cookie doesnt exist
//			if err != nil {
//				c.String(http.StatusUnauthorized, "Not Authorized")
//				return handler(c)
//			}
//			//	get jwt token
//			tokenString := cookie.Value
//
//			claims := &JwtClaims{}
//
//			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//				return []byte("secret key"), nil
//			})
//
//			if err != nil {
//				if err == jwt.ErrSignatureInvalid {
//					c.String(http.StatusUnauthorized, "Not Authorized")
//					return handler(c)
//				}
//				c.String(http.StatusBadRequest, "Something is wrong")
//				return handler(c)
//			}
//			if !token.Valid {
//				c.String(http.StatusUnauthorized, "Not Authorized")
//				return handler(c)
//			}
//
//			return handler(c)
//		}
//	}
//}

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
