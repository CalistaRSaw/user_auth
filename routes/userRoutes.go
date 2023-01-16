package routes

import (
	"github.com/CalistaRSaw/uauth-vbasic/controllers"
	"github.com/CalistaRSaw/uauth-vbasic/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate()) // protected route
	incomingRoutes.GET("/users", controllers.GetAllUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	incomingRoutes.POST("/createorg", controllers.CreateOrg())
	incomingRoutes.GET("/orgs", controllers.GetAllOrgs())
	incomingRoutes.POST("/:orgname/assignmanager/:user_id", controllers.AssignManager()) // manager user id
}
