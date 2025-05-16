package controller

import (
	"context"
	"errors"
	"github.com/ZaharBorisenko/JWT-Auth-Services/helpers"
	"github.com/ZaharBorisenko/JWT-Auth-Services/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUserId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		userID := c.Param("user_id")

		if err := helpers.UserRoleType(c, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		if err := db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found!"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword() {

}

func VerifyPassword() {

}
