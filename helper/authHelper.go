package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.MustGet("category")

	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this resource: authHelper")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("category")
	uid := c.GetString("uid")
	err = nil

	if userType != "ADMIN" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
