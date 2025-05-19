package controller

import (
	"context"
	"errors"
	"github.com/ZaharBorisenko/JWT-Auth-Services/helpers"
	"github.com/ZaharBorisenko/JWT-Auth-Services/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

var validate = validator.New()

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

func Signup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		var emailCount int64
		var phoneCount int64

		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		//парсинг данных
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//валидация данных
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//проверка уникальности email
		if err := db.WithContext(ctx).Model(&models.User{}).
			Where("email = ?", user.Email).
			Count(&emailCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while checking email",
			})
			return
		}
		if emailCount > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": "this email already exists",
			})
			return
		}

		//проверка уникальности телефона
		if err := db.WithContext(ctx).Model(&models.User{}).
			Where("phone = ?", user.Phone).
			Count(&phoneCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while checking phone",
			})
			return
		}
		if phoneCount > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": "this phone already exists",
			})
			return
		}

		//хэширование пароля
		hashedPassword, err := helpers.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to hashing password",
			})
			return
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserId)
		user.Password = hashedPassword
		user.Token = &token
		user.RefreshToken = &refreshToken

		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create user",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user":    user,
		})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func VerifyPassword() {

}
