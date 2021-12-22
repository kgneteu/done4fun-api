package main

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
	"log"
	"server/models"
)

func Seed(tx *gorm.DB, u *models.User) {
	var users []models.User

	for i := 0; i < 1000; i++ {
		user := models.User{}
		err := faker.FakeData(&user)
		if err != nil {
			fmt.Println(err)
		}
		user.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
		if i < 500 {
			user.Role = "kid"
			user.ParentId = 4
		} else {
			user.Role = "parent"
		}
		users = append(users, user)
		fmt.Println(user)
	}
	//admin
	admin := models.User{}
	admin.Email = "admin@admin.com"
	admin.FirstName = "John"
	admin.LastName = "Admin"
	admin.Verified = true
	admin.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	admin.Role = "admin"
	res := tx.Create(&admin)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	//parent
	parent := models.User{}
	parent.Email = "parent@parent.com"
	parent.FirstName = "Adam"
	parent.LastName = "Parent"
	parent.Verified = true
	parent.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	parent.Role = "parent"
	res = tx.Create(&parent)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	//kid
	kid := models.User{}
	kid.Email = "kid@kid.com"
	kid.FirstName = "Monica"
	kid.LastName = "Clever"
	kid.Verified = true
	kid.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	kid.Role = "kid"
	kid.ParentId = 2
	res = tx.Create(&kid)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	//parent with 500 kids
	parent2 := models.User{}
	parent2.Email = "parent2@parent.com"
	parent2.FirstName = "Adam"
	parent2.LastName = "Full"
	parent2.Verified = true
	parent2.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	parent2.Role = "parent"
	res = tx.Create(&parent2)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	//parent with no kids
	parent3 := models.User{}
	parent3.Email = "parent3@parent.com"
	parent3.FirstName = "Adam"
	parent3.LastName = "Empty"
	parent3.Verified = true
	parent3.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	parent3.Role = "parent"
	res = tx.Create(&parent3)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	res = tx.Create(&users)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	//fmt.Printf("%+v", user)
}
