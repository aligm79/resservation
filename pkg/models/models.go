package models

import (
	"time"
	"github.com/aligm79/reservation/pkg/config"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

type Ticket struct {
	ID			uuid.UUID	`gorm:"type:uuid;primaryKey"`
	Type		string		`gorm:"type:varchar(100);not null"`
	Remaining	int			`gorm:"not null"`
	StartsAt	time.Time	`gorm:"not null"`
	EndsAt		time.Time	`gorm:"not null"`
	CreatedDate	time.Time
}

type User struct {
	ID			uuid.UUID	`gorm:"type:uuid;primaryKey"`
	UserName	string		`gorm:"not null;unique"`
	Password	string		`gorm:"not null"`
	JoinedDate  time.Time	`gorm:"autoCreateTime"`
}

type Reserved struct {
	ID			uuid.UUID	`gorm:"type:uuid;primaryKey"`
	UserId		uuid.UUID	`gorm:"type:uuid;not null"`
	TicketId 	uuid.UUID	`gorm:"type:uuid;not null"`
	CreatedDate uuid.UUID	`gorm:"autoCreateTime"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return 
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return 
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Ticket{}, &User{}, &Reserved{})
}

func GetTickets() []Ticket {
	var tickets []Ticket
	db.Find(&tickets)
	return tickets
}

func GetTicket(id uuid.UUID) (*Ticket, error) {
	var ticket Ticket
	result := db.First(&ticket, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ticket, nil
}

func GetUserForLogin(usernmame , password string) (*User, error) {
	var user User
	result := db.Where("user_name = ?", usernmame).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return &user, nil
}