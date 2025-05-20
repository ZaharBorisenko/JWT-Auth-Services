package routes

import (
	"github.com/ZaharBorisenko/JWT-Auth-Services/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(routes *gin.Engine, db *gorm.DB) {
	//routes.Use(middleware.Authenticate())
	routes.GET("/users", controller.GetUsers())
	routes.GET("/users/:user_id", controller.GetUserId(db))
}
