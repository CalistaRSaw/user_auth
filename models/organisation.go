package models

import (
	//"github.com/CalistaRSaw/uauth-vbasic/database"
	//"github.com/gofiber/fiber"

	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Organisation struct {
	gorm.Model

	Name       string    `gorm:"unique" json:"name" validate:"required"`
	Manager    string    `json:"manager" gorm:"many2many:manager;"`
	Member     string    `json:"member" gorm:"many2many:manager;type:text;"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Problem struct {
	gorm.Model
	Title        string         `json:"title"`
	Organisation string         `json:"organisation"`
	Description  string         `json:"description"`
	Rating       int            `json:"rating" validate:"min=1,max=5"`
	Category     string         `json:"category"`
	Comments     pq.StringArray `json:"comments" gorm:"type:text[];foreignkey:ProblemRef"`
	Created_at   time.Time      `json:"created_at"`
	Updated_at   time.Time      `json:"updated_at"`
}
