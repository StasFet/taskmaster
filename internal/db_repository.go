package internal

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Returns a slice of all the users in the db. Only for use in-development
func (s *SupabaseClient) GetAllUsers() (*[]User, error) {
	client := s.GetClient()
	data, _, err := client.From(UserTableName).Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}
	results := []User{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// Returns a User object with the provided id
func (s *SupabaseClient) GetUserById(id int) (*User, error) {
	client := s.GetClient()
	data, count, err := client.From(UserTableName).Select("*", "exact", false).Eq("id", strconv.Itoa(id)).Execute()
	if err != nil {
		return nil, err
	} else if count == 0 {
		return nil, nil
	}
	result := []User{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result[0], nil
}

// Returns all the tasks in the db. only for testing and development
func (s *SupabaseClient) GetAllTasks() (*[]Task, error) {
	client := s.GetClient()
	data, _, err := client.From(TaskTableName).Select("*", "exact", false).Execute()
	if err != nil {
		return nil, err
	}
	result := []Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create new user
func (s *SupabaseClient) CreateNewUser(u *User) (*User, error) {
	client := s.GetClient()
	data, count, err := client.From(UserTableName).Insert(u, false, "", "", "exact").Execute()
	if err != nil {
		return nil, err
	} else if count != 1 {
		return nil, errors.New("create new user failed: supabase returned count > 1")
	}
	result := []User{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result[0], nil
}

// Returns a Task object with the provided id
func (s *SupabaseClient) GetTaskById(id int) (*Task, error) {
	client := s.GetClient()
	data, count, err := client.From(TaskTableName).Select("*", "exact", false).Eq("id", strconv.Itoa(id)).Execute()
	if err != nil {
		return nil, err
	} else if count > 1 {
		return nil, errors.New("more than one task with provided ID found")
	}
	results := []Task{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results[0], nil
}

// Returns a slice of tasks associated with the user with the given id
func (s *SupabaseClient) GetTasksByUUID(uuid string) (*[]Task, error) {
	client := s.GetClient()
	data, _, err := client.From(TaskTableName).Select("*", "exact", false).Eq("owner_uuid", uuid).Execute()
	if err != nil {
		return nil, err
	}
	result := []Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Inserts the provided task into the database. Returns the created task
func (s *SupabaseClient) CreateNewTask(newTask *Task) (*Task, error) {
	client := s.GetClient()
	data, count, err := client.From(TaskTableName).Insert(newTask, false, "", "", "exact").Execute()
	if err != nil {
		return nil, err
	} else if count != 1 {
		return nil, errors.New("insertion failed: count was not 1")
	}
	result := []Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result[0], nil
}

// Updates the task with the given id to match the given task
func (s *SupabaseClient) UpdateTask(id int, targetTask *Task) (*Task, error) {
	// important to note that the user can only change the following values: title, description, due date, priority and points
	client := s.GetClient()
	id_str := strconv.Itoa(id)
	data, count, err := client.From(TaskTableName).Update(map[string]any{
		"title": targetTask.Title,
		"description": targetTask.Description,
		"due_date": targetTask.DueDate,
		"priority": targetTask.Priority,
		"points": targetTask.Points,
	}, "", "exact").Eq("id", id_str).Execute()

	if err != nil {
		return nil, err
	} else if count == 0 {
		return nil, errors.New("no rows were updated")
	}
	
	result := []Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New("unmarshalled data from updated task had length 0")
	}
	return &result[0], nil
}