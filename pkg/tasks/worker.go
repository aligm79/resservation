package tasks

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/hibiken/asynq"
)

func HandleLogUserTask(ctx context.Context, t *asynq.Task) error {
    var payload AdderPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return fmt.Errorf("failed to unmarshal task payload: %w", err)
    }

    fmt.Print(payload.X, payload.Y)
    return nil
}