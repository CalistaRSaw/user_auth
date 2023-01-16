package routes

import (
	"github.com/CalistaRSaw/uauth-vbasic/controllers"
	"github.com/CalistaRSaw/uauth-vbasic/middleware"
	"github.com/gin-gonic/gin"
)

func OrgRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())                            // protected route
	incomingRoutes.POST("/adduser/:orgname/:user_id", controllers.AddUser()) // manager user id
	incomingRoutes.POST("/:orgname/addproblem", controllers.AddProblem())
	incomingRoutes.GET("/:orgname/getproblem/:problem_id", controllers.GetProblemByID())
	incomingRoutes.POST("/:orgname/:problem_id/addcomment", controllers.AddComment())
}
