package models

import (
	//"github.com/CalistaRSaw/uauth-vbasic/database"
	//"github.com/gofiber/fiber"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name          string    `json:"name" validate:"required"`
	Email         string    `gorm:"unique" json:"email" validate:"email,required"`
	Password      string    `json:"password" validate:"required,min=6"`
	Category      string    `json:"category"` // validate:"required, eq=ADMIN|eq=ORG|ORGUSER"
	Organisation  string    `json:"organisation"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	UserID        string    `json:"user_id"`
}
