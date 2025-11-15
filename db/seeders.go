package db

import (
	"log"
	models "task/src/model"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	SeedRoles()
	SeedUsers()
}

func SeedRoles() {
	roles := []models.Role{
		{Name: "admin", Status: true},
		{Name: "user", Status: true},
	}

	for _, role := range roles {
		err := DB.Where("name = ?", role.Name).FirstOrCreate(&role).Error
		if err != nil {
			log.Println("Failed seeding role:", err)
		}
	}
}

func SeedUsers() {
	var adminRole models.Role
	if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Println("Admin role not found:", err)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)

	user := models.Users{
		Fullname: "Super Admin",
		Username: "superadmin",
		Email:    "superadmin@gmail.com",
		Password: string(hashedPassword),
		RoleID:   adminRole.ID,
	}

	err := DB.Where("email = ?", user.Email).FirstOrCreate(&user).Error
	if err != nil {
		log.Println("Failed seeding user:", err)
	}
}
