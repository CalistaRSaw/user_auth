package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/CalistaRSaw/uauth-vbasic/database"
	helpers "github.com/CalistaRSaw/uauth-vbasic/helper"
	"github.com/CalistaRSaw/uauth-vbasic/models"
	"github.com/gin-gonic/gin"
)

type Comment struct {
	Comment    string `json:"comment"`
	ProblemRef int64
}

func CreateOrg() gin.HandlerFunc {
	// var db = database.DBConnOrg
	return func(c *gin.Context) {
		err := helpers.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var organisation models.Organisation

		if err := c.BindJSON(&organisation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		validationErr := validate.Struct(organisation)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var test models.Organisation

		database.DBConnOrg.Find(&test, "name = ?", organisation.Name)

		if test.Name != "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Organisation already exist",
			})
			return
		}

		organisation.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		organisation.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		// organisation.Member = make([]string, 0)

		// create user
		res := database.DBConnOrg.WithContext(ctx).Create(&organisation)

		if res.Error != nil {
			msg := "Organisation not created"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   res.Error,
				"message": msg,
				"what":    test})
			return
		}
		defer cancel()

		// respond
		c.JSON(http.StatusOK, gin.H{"message": "Organisation created"})
	}
}

func GetAllOrgs() gin.HandlerFunc {
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
		var db = database.DBConnOrg
		var organisations []models.Organisation
		result := db.WithContext(ctx).Find(&organisations)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occured"})
			return
		}
		c.JSON(http.StatusOK, organisations)
	}
}

func AssignManager() gin.HandlerFunc {
	var dbUser = database.DBConn
	var dbOrg = database.DBConnOrg
	return func(c *gin.Context) {
		err := helpers.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId := c.Param("user_id") // manager user id
		org := c.Param("orgname")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		result := dbUser.WithContext(ctx).First(&user, userId)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "user not found",
			})
			return
		}

		var organisation models.Organisation
		result = dbOrg.WithContext(ctx).First(&organisation, "name =?", org)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "organisation not found"})
			return
		}

		Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		organisation.Manager = organisation.Manager + ";" + user.UserID

		result = dbUser.Model(&user).Updates(models.User{Category: "MANAGER", Organisation: org, Updated_at: Updated_at})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "Update failed, user role not changed",
			})
			return
		}

		result = dbOrg.Model(&organisation).Updates(models.Organisation{Manager: organisation.Manager, Updated_at: Updated_at})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "Update failed, manager not assigned",
			})
			return
		}

		c.JSON(http.StatusOK, organisation)

	}
}

func AddUser() gin.HandlerFunc {
	var dbUser = database.DBConn
	var dbOrg = database.DBConnOrg
	return func(c *gin.Context) {
		err1 := helpers.CheckUserType(c, "MANAGER")
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
			return
		}

		userId := c.Param("user_id") // manager user id
		org := c.Param("orgname")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		result := dbUser.WithContext(ctx).First(&user, userId)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "user not found",
			})
			return
		}

		if user.Organisation != "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user already in an organisation",
			})
			return
		}

		var organisation models.Organisation
		result = dbOrg.WithContext(ctx).First(&organisation, "name =?", org)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "organisation not found"})
			return
		}

		Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		organisation.Member = organisation.Member + ";" + user.UserID

		result = dbOrg.Model(&organisation).Updates(models.Organisation{Member: organisation.Member, Updated_at: Updated_at})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "Update failed, user not assigned",
			})
			return
		}

		result = dbUser.Model(&user).Updates(models.User{Category: "ORGUSER", Organisation: org, Updated_at: Updated_at})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "Update failed, user role not changed",
			})
			return
		}

		c.JSON(http.StatusOK, organisation)

	}
}

func AddProblem() gin.HandlerFunc {
	var dbUser = database.DBConn
	// var dbOrg = database.DBConnOrg
	var dbProb = database.DBConnProb
	return func(c *gin.Context) {
		uid := c.GetString("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		result := dbUser.WithContext(ctx).First(&user, uid)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "user not found",
			})
			return
		}

		org := c.Param("orgname")

		if user.Organisation != org {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user does not belong to this organisation",
			})
			return
		}

		var problem models.Problem

		if err := c.BindJSON(&problem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		problem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		problem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		problem.Organisation = user.Organisation
		problem.Comments = make([]string, 0)

		res := dbProb.WithContext(ctx).Create(&problem)

		if res.Error != nil {
			msg := "Problem not created"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   res.Error,
				"message": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, problem)

	}

}

func GetProblemByID() gin.HandlerFunc {
	dbUser := database.DBConn
	dbProb := database.DBConnProb
	return func(c *gin.Context) {
		problemID := c.Param("problem_id")
		orgname := c.Param("orgname")

		uid := c.GetString("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		result := dbUser.WithContext(ctx).First(&user, uid)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "user not found",
			})
			return
		}

		if user.Organisation != orgname {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user does not belong to this organisation",
			})
			return
		}

		var problem models.Problem
		result = dbProb.WithContext(ctx).First(&problem, problemID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "requested problem not found",
			})
			return
		}

		c.JSON(http.StatusOK, problem)
	}
}

func AddComment() gin.HandlerFunc {
	dbUser := database.DBConn
	dbProb := database.DBConnProb
	return func(c *gin.Context) {
		problemID := c.Param("problem_id")
		orgname := c.Param("orgname")

		uid := c.GetString("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		result := dbUser.WithContext(ctx).First(&user, uid)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "user not found",
			})
			return
		}

		if user.Organisation != orgname {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user does not belong to this organisation",
			})
			return
		}

		var problem models.Problem
		result = dbProb.WithContext(ctx).First(&problem, problemID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "requested problem not found",
			})
			return
		}

		var comment Comment
		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		problem.Comments = append(problem.Comments, comment.Comment)

		Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result = dbProb.Model(&problem).Updates(models.Problem{Comments: problem.Comments, Updated_at: Updated_at})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   result.Error,
				"message": "Update failed",
			})
			return
		}

		c.JSON(http.StatusOK, problem)
	}
}
