package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	User userHandler
	Auth authHandler
}

func Routes(route Route) http.Handler {
	r := gin.New()
	r.HandleMethodNotAllowed = true

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.MaxMultipartMemory = 8 << 20

	v1 := r.Group("/api/v1")
	v1.POST("/user", route.User.PostUserHandler)

	auth := v1.Group("/auth")
	{
		auth.POST("/", route.Auth.PostLoginHandler)
		auth.PUT("/", route.Auth.PutAccessTokenHandler)
		auth.DELETE("/", route.Auth.DeleteLogoutHandler)
	}
	return r
}
