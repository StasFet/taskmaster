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
	} else if count > 1 {
		return nil, errors.New("more than one user with provided ID found")
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
func (s *SupabaseClient) GetTasksByUserId(id int) (*[]Task, error) {
	client := s.GetClient()
	data, _, err := client.From(TaskTableName).Select("*", "exact", false).Eq("user_id", strconv.Itoa(id)).Execute()
	if err != nil {
		return nil, err
	}
	result := []Task{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
