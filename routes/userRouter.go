package routes

import (
	"github.com/ZaharBorisenko/JWT-Auth-Services/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func UserRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authenticate())
	routes.GET("users", controller.GetUsers())
	routes.GET("users/:user_id", controller.GetUserId(DB))
}
