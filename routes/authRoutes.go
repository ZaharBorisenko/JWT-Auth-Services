package routes

import (
	"github.com/ZaharBorisenko/JWT-Auth-Services/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(routes *gin.Engine, db *gorm.DB) {
	routes.POST("/user/signup", controller.Signup(db))
	routes.POST("/user/login", controller.Login(db))
}
