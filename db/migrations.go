package db

import models "task/src/model"

func AutoMigrateModels() []interface{} {
	return []interface{}{
		&models.Users{},
		&models.Role{},
		&models.Field{},
		&models.Booking{},
	}
}
