package main

import (
	"fmt"

	routes "github.com/CalistaRSaw/uauth-vbasic/routes"

	"github.com/CalistaRSaw/uauth-vbasic/database"
	"github.com/CalistaRSaw/uauth-vbasic/initializers"
	"github.com/CalistaRSaw/uauth-vbasic/models"
	"github.com/gin-gonic/gin"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	initializers.LoadEnvVariables()
	ConnectToDb()
}

func ConnectToDb() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{}) //"sqlite3", "users.db")
	if err != nil {
		panic("failed to connect database")
	}

	database.DBConnOrg, err = gorm.Open(sqlite.Open("organisations.db"), &gorm.Config{}) //"sqlite3", "users.db")
	if err != nil {
		panic("failed to connect database")
	}

	dsn := "host=localhost user=postgres password=090901Cl* dbname=problem port=5432 sslmode=disable"
	database.DBConnProb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection opened to database")
	database.DBConn.AutoMigrate(&models.User{})
	database.DBConnOrg.AutoMigrate(&models.Organisation{})
	database.DBConnProb.AutoMigrate(&models.Problem{})
	fmt.Println("Database Migrated")
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())

	routes.AuthRoutes(r)
	routes.UserRoutes(r)
	routes.OrgRoutes(r)

	r.GET("api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	r.GET("api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	r.GET("api-3", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-3"})
	})

	r.Run() // listen and serve on 0.0.0.0:PORT

}
