package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

type AdderPayload struct {
	X int
	Y int	
}

const AdderTask = "task:adder"
const PeriodicHello = "task:periodicHello"

func Adder(x, y int) (*asynq.Task, error) {
	payload, err := json.Marshal(AdderPayload{X : x, Y : y})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(AdderTask, payload), nil
}