package main

import (
	"flag"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"server/models"
)

type config struct {
	prod bool
	dev  bool
}

func main() {
	var cfg config

	flag.BoolVar(&cfg.prod, "prod", false, "Populate db for production use")
	flag.BoolVar(&cfg.dev, "dev", false, "Populate db for development use")
	flag.Parse()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=ufo password=!!!UfO:-)1234!!! dbname=done4fun port=5432 sslmode=disable TimeZone=Europe/Warsaw",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	//db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create persons table
		// add pets table
		{
			ID: "201608301431",
			Migrate: func(tx *gorm.DB) error {
				res := tx.AutoMigrate(&models.User{})
				if cfg.dev {
					Seed(tx, &models.User{})
				}
				return res
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		{
			ID: "201608301432",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Message{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("messages")
			},
		},
		{
			ID: "201608301433",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.MessageStatus{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("message_statuses")
			},
		},
		{
			ID: "201608301434",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Prize{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("prizes")
			},
		},
		{
			ID: "201608301435",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.TaskCategory{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("task_categories")
			},
		},
		{
			ID: "201608301436",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Task{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("tasks")
			},
		},
		{
			ID: "201608301437",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.KidPrize{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("kid_prizes")
			},
		},
	})

	m.RollbackTo("201608301431")
	m.RollbackLast()
	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}

func Seed(tx *gorm.DB, u *models.User) {
	var users []models.User

	for i := 0; i < 1000; i++ {
		user := models.User{}
		err := faker.FakeData(&user)
		if err != nil {
			fmt.Println(err)
		}
		user.Password = "$2a$10$pYs2rPQYL7vrYVB/i07WfuHVrGVdEbllPLZAr7IUUWzOqKgOnpvmu"
		user.Role = "parent"
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

	res = tx.Create(&users)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	//fmt.Printf("%+v", user)
}
