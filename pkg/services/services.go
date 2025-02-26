package services

import (
	"github.com/aligm79/reservation/pkg/config"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/aligm79/reservation/pkg/models"
)

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
}


func GetTickets() []models.Ticket {
	var tickets []models.Ticket
	db.Find(&tickets)
	return tickets
}

func GetTicket(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	result := db.First(&ticket, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ticket, nil
}

func ReserveTicket(r *models.Reserved) bool {
	if err := db.Create(&r).Error; err != nil {
		return false
	}
	return true
}

func GetUserForLogin(usernmame , password string) (*models.User, error) {
	var user models.User
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
