package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CalistaRSaw/uauth-vbasic/database"
	helpers "github.com/CalistaRSaw/uauth-vbasic/helper"
	"github.com/CalistaRSaw/uauth-vbasic/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) string {
	//hash pass
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}
	return string(hash)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := ""

	if err != nil {
		msg = "email or password is incorrect"
		check = false
	}

	return check, msg
}

func Signup() gin.HandlerFunc {
	var db = database.DBConn
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var test models.User

		result := db.Find(&test, "email = ?", user.Email)

		if result.Error == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Email already exist",
			})
			return
		}

		password := HashPassword(user.Password)

		user.Password = password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// create user
		res := database.DBConn.WithContext(ctx).Create(&user)

		if res.Error != nil {
			msg := "User not created"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   res.Error,
				"message": msg})
			return
		}
		defer cancel()

		userID := strconv.FormatUint(uint64(user.ID), 10)

		result = db.Model(&user).Updates(models.User{UserID: userID})

		if result.Error != nil {
			log.Panic(result.Error)
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.Name, user.Category, userID)
		user.Token = token
		user.Refresh_token = refreshToken

		result = db.Model(&user).Updates(models.User{Token: token, Refresh_token: refreshToken})

		if result.Error != nil {
			log.Panic(result.Error)
			return
		}

		// respond
		c.JSON(http.StatusOK, gin.H{"message": "User created"})
	}
}

func Login() gin.HandlerFunc {
	var db = database.DBConn
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		result := db.WithContext(ctx).Find(&foundUser, "email = ?", user.Email)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "email or password is incorrect",
			})
			return
		}

		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.Email, foundUser.Name, foundUser.Category, foundUser.UserID)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		result = db.WithContext(ctx).Find(&foundUser, "user_id = ?", foundUser.UserID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.Writer.Header().Set("token", user.Token)
		// respond
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helpers.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))
		var db = database.DBConn
		var users []models.User
		result := db.WithContext(ctx).Find(&users)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occured"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	var db = database.DBConn
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} // check if admin or not
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		result := db.WithContext(ctx).First(&user, userId)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
