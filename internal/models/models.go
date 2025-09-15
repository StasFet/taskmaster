package models

import "time"

const TaskTableName = "tasks"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
	Points      int       `json:"points"`
	OwnerUUID   string    `json:"owner_uuid"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"completion_status"`
}

const UserTableName = "users"

type User struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	UUID        string    `json:"uuid"`
	TotalPoints int64     `json:"total_points"`
	CreatedAt   time.Time `json:"created_at"`
}

// Returns a boolean that is True if the due date of a task is past
func (t *Task) IsOverdue() bool {
	return t.DueDate.Before(time.Now())
}
