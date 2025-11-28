package config

import (
	"aadhaar-user-service/internals/database"
	"aadhaar-user-service/models/users"
)

func Automigration() {
	database.Client().AutoMigrate(&users.User{})
}
