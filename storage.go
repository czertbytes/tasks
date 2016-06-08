package tasks

import (
	"fmt"
	"sync"
)

// TaskStorage is interface which defines task storage operations.
type TaskStorage interface {
	// Insert stores new Task in storage under given TaskID path.
	Insert([]TaskID, *Task) error
	// Find returns Task from given TaskID path.
	Find([]TaskID) (Task, error)
	// FindAll returns all root Tasks (with their children).
	FindAll() ([]Task, error)
	// Update updates Task at given TaskID path.
	Update([]TaskID, *Task) error
	// Delete removes Task at given TaskID path.
	Delete([]TaskID) error
	// NextTaskID returns next available TaskID.
	NextTaskID() TaskID
}

// TaskMemoryStorage is simple implementation of TaskStorage as hashmap tree.
// This structure is not persisted so it will disappear after shuting down the
// program.
type TaskMemoryStorage struct {
	// storage is top level tree hashmap.
	storage map[TaskID]*Task

	// LastTaskID is the value of next inserted TaskID.
	lastTaskID TaskID
	// LastTaskIDMu is mutex which garantees that next TaskID will not be used
	// twice or cause some threading issues.
	lastTaskIDmu *sync.Mutex
}

// NewTaskMemoryStorage returns a new instance of TaskMemoryStorage
func NewTaskMemoryStorage() *TaskMemoryStorage {
	return &TaskMemoryStorage{
		storage:      map[TaskID]*Task{},
		lastTaskIDmu: &sync.Mutex{},
	}
}

// Insert stores new Task in storage. Path is the key where new task
// will be stored WITHOUT TaskID of the new Task.
// Insert implements TaskStorage interface.
func (s *TaskMemoryStorage) Insert(path []TaskID, task *Task) error {
	if len(path) == 0 {
		s.storage[task.ID] = task
		return nil
	}

	// We need to find last Task in the path and add Task in his children.
	lastPathTask, err := s.search(path)
	if err != nil {
		fmt.Println("(DEBUG) storage: Insert Task by TaskID path failed. Task not found.")
		return err
	}

	if lastPathTask.Children == nil {
		lastPathTask.Children = map[TaskID]*Task{}
	}
	lastPathTask.Children[task.ID] = task

	return nil
}

// Find returns Task under given taskID path.
// Find implements TaskStorage interface.
func (s *TaskMemoryStorage) Find(path []TaskID) (Task, error) {
	if len(path) == 0 {
		fmt.Println("(DEBUG) storage: Find Task by TaskID path failed. TaskID path is empty.")
		return Task{}, ErrTaskPathNotValid
	}

	lastPathTask, err := s.search(path)
	if err != nil {
		fmt.Println("(DEBUG) storage: Find Task by TaskID path failed. Task not found.")
		return Task{}, err
	}

	return *lastPathTask, nil
}

// FindAll returns all root (top level) Tasks with their children.
// FindAll implements TaskStorage interface.
func (s *TaskMemoryStorage) FindAll() ([]Task, error) {
	tasks := []Task{}
	for _, task := range s.storage {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

// Update updates Task under given TaskID path.
// Update implements TaskStorage interface.
func (s *TaskMemoryStorage) Update(path []TaskID, task *Task) error {
	if len(path) == 0 {
		fmt.Println("(DEBUG) storage: Update Task by TaskID path failed. TaskID path is empty.")
		return ErrTaskPathNotValid
	}

	// Look in top level tasks
	if len(path) == 1 {
		if _, found := s.storage[task.ID]; !found {
			fmt.Println("(DEBUG) storage: Update Task by TaskID path failed. Root Task not found.")
			return ErrTaskNotFound
		}

		s.storage[task.ID] = task

		return nil
	}

	// Find the Task before last in the path. The last one is the Task we want
	// to update.
	lastPathTask, err := s.search(path[:len(path)-1])
	if err != nil {
		fmt.Println("(DEBUG) storage: Update Task by TaskID path failed. Child Task not found.")
		return err
	}

	if lastPathTask.Children == nil {
		lastPathTask.Children = map[TaskID]*Task{}
	}
	lastPathTask.Children[task.ID] = task

	return nil
}

// Delete removes Task at given TaskID path.
// Delete implements TaskStorage interface.
func (s *TaskMemoryStorage) Delete(path []TaskID) error {
	if len(path) == 0 {
		fmt.Println("(DEBUG) storage: Search Task by TaskID path failed. TaskID path is empty.")
		return ErrTaskPathNotValid
	}

	// Look in top level tasks
	if len(path) == 1 {
		if _, found := s.storage[path[0]]; !found {
			fmt.Println("(DEBUG) storage: Delete Task by TaskID path failed. Root Task not found.")
			return ErrTaskNotFound
		}

		delete(s.storage, path[0])

		return nil
	}

	// Find the Task before last in the path. The last one is the Task we want
	// to delete.
	taskID := path[len(path)-1]
	lastPathTask, err := s.search(path[:len(path)-1])
	if err != nil {
		fmt.Println("(DEBUG) storage: Delete Task by TaskID path failed. Child Task not found.")
		return err
	}

	delete(lastPathTask.Children, taskID)

	return nil
}

// NextTaskID returns next available TaskID value.
// NextTaskID implements TaskStorage interface.
func (s *TaskMemoryStorage) NextTaskID() TaskID {
	s.lastTaskIDmu.Lock()
	defer s.lastTaskIDmu.Unlock()

	s.lastTaskID += 1

	return s.lastTaskID
}

// search returns Task at given TaskID path.
func (s *TaskMemoryStorage) search(path []TaskID) (*Task, error) {
	if len(path) == 0 {
		fmt.Println("(DEBUG) storage: Search Task by TaskID path failed. TaskID path is empty.")
		return nil, ErrTaskPathNotValid
	}

	// Look in top level Tasks.
	task, found := s.storage[path[0]]
	if !found {
		fmt.Println("(DEBUG) storage: Search Task by TaskID path failed. Root Task not found.")
		return nil, ErrTaskNotFound
	}

	// Loop in Task children to find the match.
	for i := 1; i < len(path); i++ {
		found = false
		for _, childTask := range task.Children {
			if childTask.ID == path[i] {
				found = true
				task = childTask
				continue
			}
		}

		if !found {
			fmt.Println("(DEBUG) storage: Search Task by TaskID path failed. Child Task not found.")
			return nil, ErrTaskNotFound
		}
	}

	return task, nil
}
