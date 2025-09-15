package database

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	model "taskmaster/internal/models"
)

// Returns a slice of all the users in the db. Only for use in-development
func (s *SupabaseClient) GetAllUsers() (*[]model.User, error) {
	client := s.GetClient()
	data, _, err := client.From(model.UserTableName).Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}
	results := []model.User{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// Returns a User object with the provided id
func (s *SupabaseClient) GetUserByUUID(uuid string) (*model.User, error) {
	client := s.GetClient()
	data, count, err := client.From(model.UserTableName).Select("*", "exact", false).Eq("uuid", uuid).Execute()
	if err != nil {
		return nil, err
	} else if count == 0 {
		return nil, nil
	}
	result := []model.User{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result[0], nil
}

// Returns all the tasks in the db. only for testing and development
func (s *SupabaseClient) GetAllTasks() (*[]model.Task, error) {
	client := s.GetClient()
	data, _, err := client.From(model.TaskTableName).Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}
	result := []model.Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create new user. Automatically sets CreatedAt
func (s *SupabaseClient) CreateNewUser(u *model.User) (*model.User, error) {
	client := s.GetClient()
	u.CreatedAt = time.Now()
	data, count, err := client.From(model.UserTableName).Insert(u, false, "", "", "exact").Execute()
	if err != nil {
		return nil, err
	} else if count != 1 {
		return nil, errors.New("create new user failed: supabase returned count > 1")
	}
	result := []model.User{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result[0], nil
}

// Returns a Task object with the provided id
func (s *SupabaseClient) GetTaskById(id int) (*model.Task, error) {
	client := s.GetClient()
	data, count, err := client.From(model.TaskTableName).Select("*", "exact", false).Eq("id", strconv.Itoa(id)).Execute()
	if err != nil {
		return nil, err
	} else if count > 1 {
		return nil, errors.New("more than one task with provided ID found")
	}
	results := []model.Task{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results[0], nil
}

// Returns a slice of tasks associated with the user with the given id
func (s *SupabaseClient) GetTasksByUUID(uuid string) (*[]model.Task, error) {
	client := s.GetClient()
	data, _, err := client.From(model.TaskTableName).Select("*", "exact", false).Eq("owner_uuid", uuid).Eq("completed", "FALSE").Execute()
	if err != nil {
		return nil, err
	}
	result := []model.Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Inserts the provided task into the database. Returns the created task
func (s *SupabaseClient) CreateNewTask(newTask *model.Task) (*model.Task, error) {
	client := s.GetClient()
	data, count, err := client.From(model.TaskTableName).Insert(newTask, false, "", "", "exact").Execute()
	if err != nil {
		return nil, err
	} else if count != 1 {
		return nil, errors.New("insertion failed: count was not 1")
	}
	result := []model.Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result[0], nil
}

// Updates the task with the given id to match the given task
func (s *SupabaseClient) UpdateTask(targetTask *model.Task) (*model.Task, error) {
	// important to note that the user can only change the following values: title, description, due date, priority and points
	client := s.GetClient()
	id_str := strconv.Itoa(targetTask.ID)
	data, count, err := client.From(model.TaskTableName).Update(map[string]any{
		"title":             targetTask.Title,
		"description":       targetTask.Description,
		"due_date":          targetTask.DueDate,
		"priority":          targetTask.Priority,
		"points":            targetTask.Points,
		"completion_status": targetTask.Status,
	}, "", "exact").Eq("id", id_str).Execute()

	if err != nil {
		return nil, err
	} else if count == 0 {
		return nil, errors.New("no rows were updated")
	}

	result := []model.Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New("unmarshalled data from updated task had length 0")
	}
	return &result[0], nil
}

// Updates the user with the given id to match the given user
func (s *SupabaseClient) UpdateUser(uuid string, targetUser *model.User) (*model.User, error) {
	// the user can only change the following details: name, totalpoints
	client := s.GetClient()
	data, count, err := client.From(model.UserTableName).Update(map[string]any{
		"name":         targetUser.Name,
		"total_points": targetUser.TotalPoints,
	}, "", "exact").Eq("uuid", uuid).Execute()
	if err != nil {
		return nil, errors.New("error updating user details")
	} else if count == 0 {
		return nil, errors.New("error updating user details: no users have been changed")
	}

	// unmarshal user returned by supabase
	result := []model.User{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New("updated user")
	}
	return &result[0], nil
}

// Deletes the task with the given id
func (s *SupabaseClient) DeleteTask(id string, owner_uuid string) error {
	client := s.GetClient()
	_, count, err := client.From(model.TaskTableName).Delete("", "exact").Eq("id", id).Eq("owner_uuid", owner_uuid).Execute()
	if err != nil {
		return err
	} else if count != 1 {
		return errors.New("the number of matching tasks to delete was is not 1")
	}
	return nil
}
