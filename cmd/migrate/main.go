package main

import (
	"flag"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"server/models"
)

type config struct {
	prod bool
	dev  bool
}

func main() {
	var cfg config

	defEnv := os.Getenv("ENV_TYPE")
	var err error
	if defEnv == "" {
		defEnv = "development"
		err = godotenv.Load(".env.development")
	} else {
		err = godotenv.Load(".env")
	}

	if err != nil {
		println(err.Error())
	}

	defDSN := os.Getenv("POSTGRES_DSN")

	flag.BoolVar(&cfg.prod, "prod", false, "Populate db for production use")
	flag.BoolVar(&cfg.dev, "dev", false, "Populate db for development use")
	flag.Parse()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  defDSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	//db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "201608301431",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{})
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
				return tx.AutoMigrate(&models.TaskStatus{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("task_statuses")
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
				return tx.AutoMigrate(&models.PrizeStatus{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("prize_statuses")
			},
		},
	})

	_ = m.RollbackTo("201608301431")
	_ = m.RollbackLast()
	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	if cfg.dev {
		Seed(db, &models.User{})
	} else {
		ProdSeed(db, &models.User{})
	}
	log.Printf("Migration did run successfully")
}
