package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrTaskNotFound
	ErrTaskNotFound error = errors.New("Task not found")
	// ErrTaskPathNotValid
	ErrTaskPathNotValid error = errors.New("Task path is not valid")
	// ErrTaskLabelIsRequired
	ErrTaskLabelIsRequired error = errors.New("Task field Label is required")
	// ErrTaskLabelIsNotValid
	ErrTaskLabelIsNotValid error = errors.New("Task field Label is not valid")
	// ErrTaskLabelOrCompletedRequired
	ErrTaskLabelOrCompletedRequired error = errors.New("Task field Label or Completed is required")
)

// TaskID is alias for int type.
type TaskID int

// Task is struct which holds data for Task.
type Task struct {
	// ID is identifier for given Task
	ID TaskID `json:"id,string"`
	// Label is name of the Task.
	Label string `json:"label"`
	// Completed identifies if given Task is completed.
	Completed bool `json:"completed"`
	// Children contains tasks which have given Task as parent.
	Children SubTasks `json:"sub_tasks,omitempty"`
}

// ByTaskID is alias type for slice of Task. Used for sorting only.
type ByTaskID []Task

// Len implements sort.Sort interface.
func (a ByTaskID) Len() int { return len(a) }

// Swap implements sort.Sort interface.
func (a ByTaskID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less implements sort.Sort interface.
func (a ByTaskID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// SubTasks is alias type for map of Tasks.
type SubTasks map[TaskID]*Task

// MarshalJSON marshals Task's subtask. We must use our own marshaler because
// default Go marshaler don't know how to properly serialize map[TaskID]*Task.
// We also don't want to serialize as map but as array. This function causes
// recursive json marshaling for tasks in the tree.
// MarshalJSON implements json.Marshaler interface.
func (sb SubTasks) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")

	l := len(sb)
	for _, task := range sb {
		jsonValue, err := json.Marshal(task)
		if err != nil {
			fmt.Printf("(WARN) task: Marshaling task struct failed: %s\n", err)
			return nil, err
		}
		buffer.Write(jsonValue)

		l--
		if l > 0 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

// JSONTask represents Task is JSON request and response. This struct uses
// pointers because Go uses default values for structs so we can't distiguish
// if the value was set or not. With pointers we know that value was set (has
// value) or was not set (is nil). JSONTask also support only fields which are
// used in create and update flow.
type JSONTask struct {
	Label     *string `json:"label"`
	Completed *bool   `json:"completed"`
}

// Valid returns if current Task is valid for given action.
func (t *JSONTask) Validate(validator TaskActionValidator) error {
	return validator.Validate(t)
}

// TaskActionValidator is interface with method which validates if given task
// is valid for given operation.
type TaskActionValidator interface {
	Validate(*JSONTask) error
}

// CreateValidator implements TaskActionValidator for create operation on Task.
type CreateValidator struct{}

// NewCreateValidator returns new instance of CreateValidator.
func NewCreateValidator() *CreateValidator {
	return &CreateValidator{}
}

// Validate returns error if given task is not valid and should not be stored
// in storage.
// Validate implements TaskActionValidator.
func (v *CreateValidator) Validate(t *JSONTask) error {
	if t.Label == nil {
		fmt.Println("(DEBUG) task: Create task validation failed. Missing field Label.")
		return ErrTaskLabelIsRequired
	}

	if len(*t.Label) < 1 || len(*t.Label) > 100 {
		fmt.Println("(DEBUG) task: Create task validation failed. Field Label is not valid.")
		return ErrTaskLabelIsNotValid
	}

	return nil
}

// UpdateValidator implements TaskActionValidator for update operation on Task.
type UpdateValidator struct{}

// NewUpdateValidator returns new instance of UpdateValidator.
func NewUpdateValidator() *UpdateValidator {
	return &UpdateValidator{}
}

// Validate returns error if fiven Task is not valid and should not be updated
// in storage.
// Validate implements TaskActionValidator.
func (v *UpdateValidator) Validate(t *JSONTask) error {
	// At least one of the value should be set.
	if t.Label == nil && t.Completed == nil {
		fmt.Println("(DEBUG) task: Update task validation failed. Neither Label or Completed fields are set.")
		return ErrTaskLabelOrCompletedRequired
	}

	if t.Label != nil {
		if len(*t.Label) < 1 || len(*t.Label) > 100 {
			fmt.Println("(DEBUG) task: Update task validation failed. Field Label is not valid.")
			return ErrTaskLabelIsNotValid
		}
	}

	return nil
}
