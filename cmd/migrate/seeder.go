package main

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
	"log"
	"server/models"
)

func ProdSeed(tx *gorm.DB, u *models.User) {
	admin := models.User{}
	admin.Email = "admin@admin.com"
	admin.FirstName = "John"
	admin.LastName = "Admin"
	admin.Verified = true
	admin.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	admin.Role = "admin"
	//admin.Parent = models.User{}
	res := tx.Create(&admin)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
}

func Seed(tx *gorm.DB, u *models.User) {
	//admin ID=1
	admin := models.User{}
	admin.Email = "admin@admin.com"
	admin.FirstName = "John"
	admin.LastName = "Admin"
	admin.Verified = true
	admin.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	admin.Role = "admin"
	//admin.Parent = models.User{}
	res := tx.Create(&admin)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	//parent ID=2
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

	//kid ID=3 - kid for parent 2
	kid := models.User{}
	kid.Email = "kid@kid.com"
	kid.FirstName = "Monica"
	kid.LastName = "Clever"
	kid.Verified = true
	kid.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
	kid.Role = "kid"
	var pid uint = 2
	kid.ParentId = &pid
	res = tx.Create(&kid)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	//parent with 500 kids ID=4
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

	//parent with no kids ID=5
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

	//Kids of parent 4 (ID=6-505) & parents ID: 506-1005
	var users []models.User
	pid = 4
	for i := 0; i < 1000; i++ {
		user := models.User{}
		err := faker.FakeData(&user)
		if err != nil {
			fmt.Println(err)
		}
		user.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
		if i < 500 {
			user.Role = "kid"

			user.ParentId = &pid
		} else {
			user.Role = "parent"
		}
		users = append(users, user)
		fmt.Println(user)
	}
	res = tx.Create(&users)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	var prizes []models.Prize

	for i := 6; i < 26; i++ {
		for k := 0; k < 20; k++ {
			prize := models.Prize{}
			err := faker.FakeData(&prize)
			if err != nil {
				fmt.Println(err)
			}
			prize.KidId = uint(i)
			prizes = append(prizes, prize)
			fmt.Println(prize)
		}
	}

	res = tx.Create(&prizes)
	if res.Error != nil {
		log.Fatal(res.Error)
	}

	var tasks []models.Task

	for i := 6; i < 26; i++ {
		for k := 0; k < 40; k++ {
			task := models.Task{}
			err := faker.FakeData(&task)
			if err != nil {
				fmt.Println(err)
			}
			task.KidId = uint(i)
			tasks = append(tasks, task)
			fmt.Println(task)
		}
	}

	res = tx.Create(&tasks)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	//fmt.Printf("%+v", user)
}
