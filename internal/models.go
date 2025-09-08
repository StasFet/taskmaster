package internal

import "time"

const TaskTableName = "tasks"

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    int64     `json:"priority"`
	Points      int64     `json:"points"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

const UserTableName = "users"

type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	UUID        string    `json:"uuid"`
	TotalPoints int64     `json:"total_points"`
	CreatedAt   time.Time `json:"created_at"`
}
