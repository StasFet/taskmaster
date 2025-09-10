package internal

import "time"

const TaskTableName = "tasks"

type Task struct {
	ID          int       `json:"-"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	Points      int       `json:"points,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

const UserTableName = "users"

type User struct {
	Name        string    `json:"name,omitempty"`
	Email       string    `json:"email,omitempty"`
	UUID        string    `json:"uuid,omitempty"`
	TotalPoints int       `json:"total_points,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// Returns a boolean that is True if the due date of a task is past
func (t *Task) IsOverdue() bool {
	return t.DueDate.Before(time.Now())
}
