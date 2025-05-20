package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/ZaharBorisenko/JWT-Auth-Services/helpers"
	"github.com/ZaharBorisenko/JWT-Auth-Services/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
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

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Подготовка контекста и структур
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		var inputUser models.User
		var foundUser models.User

		// 2. Парсинг входящих данных
		if err := c.BindJSON(&inputUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// 3. Поиск пользователя в БД
		if err := db.WithContext(ctx).
			Where("email = ?", strings.ToLower(inputUser.Email)).
			First(&foundUser).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// 4. Проверка пароля
		passwordIsValid, msg := VerifyPassword(inputUser.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		// 5. Генерация токенов
		token, refreshToken, err := helpers.GenerateAllTokens(
			foundUser.Email,
			foundUser.FirstName,
			foundUser.LastName,
			foundUser.Role,
			foundUser.UserId,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// 6. Обновление токенов в БД
		if err := db.WithContext(ctx).
			Model(&foundUser).
			Updates(map[string]interface{}{
				"token":         token,
				"refresh_token": refreshToken,
			}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		// 7. Успешный ответ
		c.JSON(http.StatusCreated, gin.H{
			"message": "User login successfully",
			"user":    foundUser,
		})
	}
}

func VerifyPassword(userPassword, hashedPassword string) (bool, string) {
	fmt.Printf("Comparing:\nUserPassword: '%s' (% x)\nHashedPassword: '%s'\n",
		userPassword, []byte(userPassword), hashedPassword)

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
	if err != nil {
		fmt.Printf("Comparison failed: %v\n", err)
		return false, "Invalid email or password"
	}
	return true, ""
}
