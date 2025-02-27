package services

import (
	"errors"

	"github.com/aligm79/reservation/pkg/config"
	"github.com/aligm79/reservation/pkg/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

type myTickets struct {
	ID 		uuid.UUID	`json:"ID"`
    models.Ticket		`json:"Ticket"`
	Status 	int			`json:"Status"`
}

func MyTickets(userId uuid.UUID, page int, pageSize int) ([]myTickets, error) {
	var myTickets []myTickets
	offset := (page - 1) * pageSize

	err := db.Table("reserveds").
		Select("tickets.*, reserveds.status as Status, reserveds.id as ID").
		Joins("JOIN tickets ON reserveds.ticket_id = tickets.id").
		Where("reserveds.user_id = ?", userId).
		Find(&myTickets).
		Limit(pageSize).Offset(offset). 
		Find(&myTickets).Error
	if err != nil {
		return nil, err
	}
	return myTickets, nil
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
    return db.Transaction(func(tx *gorm.DB) error {
        var ticket models.Ticket
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("id = ?", r.TicketId).
            First(&ticket).Error; err != nil {
            return err 
        }

        if ticket.Remaining <= 0 {
            return errors.New("no remaining tickets")
        }

        if err := tx.Create(&r).Error; err != nil {
            return err 
        }

        ticket.Remaining -= 1
        if err := tx.Save(&ticket).Error; err != nil {
            return err 
        }

        return nil 
    }) == nil
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
