package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type DBModel struct {
	Db *sqlx.DB
}

//type DeletedAt sql.NullTime

type GormModel struct {
	ID        uint         `gorm:"primarykey" json:"id" db:"id" faker:"-"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" faker:"-"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at" faker:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at" db:"deleted_at" faker:"-"`
}

type User struct {
	GormModel
	FirstName string  `gorm:"not null" json:"first_name" db:"first_name" faker:"first_name"`
	LastName  string  `gorm:"not null" json:"last_name" db:"last_name" faker:"last_name"`
	Email     string  `gorm:"not null;unique" json:"email" db:"email" faker:"email"`
	Password  string  `gorm:"not null" json:"password" db:"password"`
	Role      string  `gorm:"default:parent" json:"role" db:"role"`
	ParentId  *uint   `gorm:"null;index;default:null" json:"parent_id" db:"parent_id" faker:"-"`
	Verified  bool    `gorm:"default:false" json:"verified" db:"verified" faker:"-"`
	Picture   *string `gorm:"null" json:"picture" db:"picture" faker:"-"`
	Points    uint    `gorm:"default:0" json:"points" db:"points" faker:"boundary_start=0, boundary_end=1000"`
	Parent    *User   `gorm:"foreignKey:parent_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" faker:"-"`
	//Messages []Message `gorm:"foreignKey:sender_id;references:id"`
}

type Task struct {
	GormModel
	KidId        uint      `gorm:"not null;index" json:"kid_id" faker:"-"`
	Kid          User      `gorm:"foreignKey:kid_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Action       string    `gorm:"not null" json:"icon" db:"icon" faker:"-"`
	Icon         uint      `gorm:"default:1" json:"icon" faker:"boundary_start=1, boundary_end=64"`
	StartAt      time.Time `json:"start_at" db:"start_at" faker:"-"`
	Cyclic       uint      `gorm:"default:0" json:"cyclic" faker:"-"`
	SelectedDays uint      `gorm:"default:0" json:"selected_days"`
	Negligible   bool      `gorm:"default:false" json:"negligible" faker:"-"`
	Deferrable   bool      `gorm:"default:false" json:"deferrable" faker:"-"`
	MaxDelay     uint      `gorm:"default:0" json:"max_delay" faker:"-"`
	Completed    bool      `gorm:"default:false" json:"completed" faker:"-"`
}

type TaskStatus struct {
	GormModel
	TaskId      uint      `gorm:"not null;index" json:"task_id"`
	Task        Task      `gorm:"foreignKey:task_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Reply       string    `gorm:"null" json:"reply" db:"reply" faker:"-"`
	CompletedAt time.Time `json:"start_at" db:"start_at" faker:"-"`
}

type Prize struct {
	GormModel
	KidId     uint   `gorm:"not null;index" json:"kid_id" db:"kid_id" faker:"-"`
	Kid       User   `gorm:"foreignKey:kid_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-" faker:"-"`
	Name      string `gorm:"not null" json:"name" db:"name" faker:"sentence, len=25"`
	Points    uint   `gorm:"not null;default:1" json:"points" db:"points" faker:"boundary_start=1, boundary_end=1000"`
	Icon      uint   `gorm:"default:1" json:"icon" db:"icon" faker:"boundary_start=1, boundary_end=64"`
	OneTime   bool   `gorm:"default:false" json:"one_time" db:"one_time" faker:"-"`
	Published bool   `gorm:"default:true" json:"published" db:"published" faker:"-"`
}

type KidPrize struct {
	GormModel
	KidId     uint      `gorm:"not null;index" json:"kid_id"`
	Kid       User      `gorm:"foreignKey:kid_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	PrizeId   uint      `gorm:"not null;index" json:"prize_id"`
	Prize     Prize     `gorm:"foreignKey:prize_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ChosenAt  time.Time `json:"chosen_at" db:"chosen_at" faker:"-"`
	Collected bool      `gorm:"default:false" json:"collected"`
}

type Message struct {
	GormModel
	SenderId         string `gorm:"not null;index" json:"sender_id"`
	Sender           User   `gorm:"foreignKey:sender_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Message          string `gorm:"not null" json:"message"`
	NeedConfirmation bool   `gorm:"default:false" json:"need_confirmation"`
}

type MessageStatus struct {
	GormModel
	RecipientId uint    `gorm:"not null;index" json:"recipient_id"`
	Recipient   User    `gorm:"foreignKey:recipient_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	MessageId   uint    `gorm:"not null;index" json:"message_id"`
	Message     Message `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Viewed      bool    `gorm:"default:false" json:"viewed"`
	Confirmed   bool    `gorm:"default:false" json:"confirmed"`
}

type Article struct {
	GormModel
	AuthorId  uint `gorm:"not null;index" json:"author_id"`
	Author    User `gorm:"foreignKey:author_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Published bool `gorm:"default:false" json:"published"`
}
