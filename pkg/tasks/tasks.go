package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/aligm79/reservation/pkg/config"
	"github.com/aligm79/reservation/pkg/models"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

var db *gorm.DB

const TenMinuteCheck = "task:10minuteCheck"

func init() {
	config.Connect()
	db = config.GetDB()
}

func HandleTenMinuteCheck(ctx context.Context, t *asynq.Task) error {
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)

	err := db.Transaction(func(tx *gorm.DB) error {
		subQuery := tx.Model(&models.Reserved{}).
			Select("id").
			Where("created_date <= ? AND status = ?", tenMinutesAgo, 0)

		if err := tx.Model(&models.Reserved{}).
			Where("id IN (?)", subQuery).
			Update("status", -1).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Ticket{}).
			Where("id IN (SELECT ticket_id FROM reserveds WHERE created_date <= ? AND status = -1)", tenMinutesAgo).
			Update("remaining", gorm.Expr("remaining + (SELECT COUNT(*) FROM reserveds WHERE ticket_id = tickets.id AND status = -1)")).
			Error; err != nil {
			return err
		}

		fmt.Println("Reservations canceled and tickets updated using subqueries.")
		return nil
	})

	return err
}