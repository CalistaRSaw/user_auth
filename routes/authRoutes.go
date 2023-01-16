package routes

import (
	"github.com/CalistaRSaw/uauth-vbasic/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controllers.Signup())
	incomingRoutes.POST("/login", controllers.Login())
}
