package workoutlog

import (
	"log"
	"time"

	"gesitr/internal/user/workoutlog/models"

	"gorm.io/gorm"
)

// StartCommitmentTicker starts a background goroutine that periodically
// transitions committed workout logs to broken when their due_end has elapsed.
func StartCommitmentTicker(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			now := time.Now()
			result := db.Model(&models.WorkoutLogEntity{}).
				Where("status = ? AND due_end < ?", models.WorkoutLogStatusCommitted, now).
				Updates(map[string]any{
					"status":            models.WorkoutLogStatusBroken,
					"status_changed_at": now,
				})
			if result.Error != nil {
				log.Printf("commitment ticker error: %v", result.Error)
			}
			if result.RowsAffected > 0 {
				log.Printf("commitment ticker: marked %d logs as broken", result.RowsAffected)
			}
		}
	}()
}
