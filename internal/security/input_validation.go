package security

import (
	model "taskmaster/internal/models"
	"unicode/utf8"
)

/*
	Validate task
	Here are the requirements that a task must meet:
	- Title and description can have any characters
		- 0 < len(title) <= 100
		- 0 < len(description) <= 500
	- Priority is between 0 and 10 (inclusive)
	- owner uuid is valid uuid, len(uuid) > 0
*/

func ValidateTask(t *model.Task) (bool, error) {
	// check name lengths
	lenName := utf8.RuneCountInString(t.Title)
	lenDesc := utf8.RuneCountInString(t.Description)
	if !(lenName > 0 && lenName <= 100) {
		return false, nil
	} else if !(lenDesc > 0 && lenDesc <= 100) {
		return false, nil
	}

	// check priority
	if !(t.Priority >= 0 && t.Priority <= 500) {
		return false, nil
	}

	return true, nil
}
