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
	ParentId  uint    `gorm:"null;index" json:"parent_id" db:"parent_id" faker:"-"`
	Verified  bool    `gorm:"default:false" json:"verified" db:"verified" faker:"-"`
	Picture   *string `gorm:"null" json:"picture" db:"picture" faker:"-"`
	//Parent   User
	//Messages []Message `gorm:"foreignKey:sender_id;references:id"`
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

type Prize struct {
	GormModel
	ParentId uint `gorm:"not null;index" json:"parent_id"`
	//Parent   User   `gorm:"foreignKey:parent_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Points uint   `gorm:"not null;default:1" json:"points"`
	Icon   string `gorm:"null" json:"icon" db:"icon" faker:"-"`
}

type KidPrize struct {
	GormModel
	KidId   uint  `gorm:"not null;index" json:"kid_id"`
	Kid     User  `gorm:"foreignKey:kid_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	PrizeId uint  `gorm:"not null;index" json:"prize_id"`
	Prize   Prize `gorm:"foreignKey:prize_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Points  uint  `gorm:"not null;default:1" json:"points"`
}

type TaskCategory struct {
	GormModel
	ParentId uint `gorm:"not null;index" json:"parent_id"`
	//Parent           User   `gorm:"foreignKey:parent_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Icon             string `gorm:"null" json:"icon" db:"icon" faker:"-"`
	NeedConfirmation bool   `gorm:"default:false" json:"need_confirmation"`
	NeedReply        bool   `gorm:"default:false" json:"need_reply"`
}

type Task struct {
	GormModel
	KidId     uint   `gorm:"not null;index" json:"kid_id"`
	Kid       User   `gorm:"foreignKey:kid_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Confirmed bool   `gorm:"default:false" json:"confirmed"`
	Reply     string `gorm:"not null" json:"reply" db:"reply" faker:"-"`
}
