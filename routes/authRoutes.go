package routes

import (
	"github.com/ZaharBorisenko/JWT-Auth-Services/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("user/signup", controller.Signup(DB))
	routes.POST("user/login", controller.Login())
}
