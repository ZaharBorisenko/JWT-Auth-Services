package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userRole := c.GetString("role")
	err = nil
	if userRole != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	return err
}

func UserRoleType(c *gin.Context, userId string) (err error) {
	userRole := c.GetString("role")
	uid := c.GetString("uid")
	err = nil

	if userRole == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userRole)
	return err
}
