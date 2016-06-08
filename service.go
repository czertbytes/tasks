package tasks

import "fmt"

// TaskService is interface which defines business logic with Task entity.
// Current TaskService is very simple and offers only CRUD operations. In real
// usage it would have real business logic methods whic interacts with other
// data providers and storages.
type TaskService interface {
	// Create creates and stores new Task in storage under given TaskID path.
	Create([]TaskID, CreateFields) (Task, error)
	// Find returns Task from given TaskID path.
	Find([]TaskID) (Task, error)
	// FindAll returns all root Tasks(with their children).
	FindAll() ([]Task, error)
	// Update updates Task at given TaskID path.
	Update([]TaskID, UpdateFields) (Task, error)
	// Delete removes Tasks at given TaskID path.
	Delete([]TaskID) (Task, error)
}

// TaskStorageService is simple implementation of TaskService working with
// given TaskStorage implementation. Current TaskStorageService has only
// CRUD operations - in real case it would have more business logic related
// operations.
type TaskStorageService struct {
	storage TaskStorage
}

// NewTaskStorageService returns new instance of TaskStorageService
func NewTaskStorageService(storage TaskStorage) *TaskStorageService {
	return &TaskStorageService{
		storage: storage,
	}
}

// CreateFields is struct which contains only allowed fields for Task in Create
// flow.
type CreateFields struct {
	Label string
}

// Create creates and stores new Task in storage under given TaskID path. Task
// is created from CreateFields provided in parameter. TaskID is received from
// TaskStorage service which guarantees unique TaskID.
// Create implements TaskService interface.
func (s *TaskStorageService) Create(path []TaskID, fields CreateFields) (Task, error) {
	// Create a new Task: copy allowed (whitelisted) fields from CreateFields
	newTask := &Task{
		ID:        TaskID(s.storage.NextTaskID()),
		Label:     fields.Label,
		Completed: false,
		Children:  SubTasks{},
	}

	if err := s.storage.Insert(path, newTask); err != nil {
		fmt.Printf("(DEBUG) service: Inserting a new Task failed: %s\n", err)
		return Task{}, err
	}

	return *newTask, nil
}

// Find returns Task from given TaskID path or error if Task is not found.
// Find implements TaskService interface.
func (s *TaskStorageService) Find(path []TaskID) (Task, error) {
	return s.storage.Find(path)
}

// FindAll returns complete Task tree in storage. Every root Task with its
// all chidren and subchildren. This can be quite verbose and huge.
// FindAll implements TaskService interface.
func (s *TaskStorageService) FindAll() ([]Task, error) {
	return s.storage.FindAll()
}

// UpdateFields is struct which contains only allowed fields for Task in
// Update flow. Notice that it contains pointers: if value field is not nil
// then it will set the value.
type UpdateFields struct {
	Label     *string
	Completed *bool
}

// Update updates Task at given TaskID path with UpdateFields provided in
// parameter. It updates only "set" fields (fields which are not nil).
// Update implements TaskService interface.
func (s *TaskStorageService) Update(path []TaskID, fields UpdateFields) (Task, error) {
	oldVersionTask, err := s.storage.Find(path)
	if err != nil {
		fmt.Printf("(DEBUG) service: Updating existing Task failed: %s\n", err)
		return Task{}, err
	}

	newVersionTask := oldVersionTask

	if fields.Label != nil {
		newVersionTask.Label = *fields.Label
	}

	if fields.Completed != nil {
		newVersionTask.Completed = *fields.Completed
	}

	if err := s.storage.Update(path, &newVersionTask); err != nil {
		fmt.Printf("(DEBUG) service: Updating existing Task failed: %s\n", err)
		return oldVersionTask, err
	}

	return newVersionTask, nil
}

// Delete removes Tasks at given TaskID path or error if Task is not found.
// Delete implements TaskService interface.
func (s *TaskStorageService) Delete(path []TaskID) (Task, error) {
	task, err := s.storage.Find(path)
	if err != nil {
		fmt.Printf("(DEBUG) service: Deleting Task failed: %s\n", err)
		return Task{}, err
	}

	if err := s.storage.Delete(path); err != nil {
		fmt.Printf("(DEBUG) service: Deleting Task failed: %s\n", err)
		return Task{}, err
	}

	return task, nil
}
