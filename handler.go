package tasks

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
)

var (
	// ErrHandlerURLNotValid is returned when URL contains TaskID values which
	// are not numbers.
	ErrHandlerURLNotValid error = errors.New("URL parameters are not valid number(s)")
)

// TaskHandler is simple Handler which handles Task which are not in root level
// in the tree hierarchy. It handles children of top level Tasks and children
// of children. Handler provides CRUD operations.
// TaskHandler implements http.Handler interface.
type TaskHandler struct {
	service TaskService
}

// NewTaskHandler returns new instance of TaskHandler
func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// ServeHTTP is simple function which dispatches requests to proper function
// handlers.
// ServeHTTP implements http.Handler interface
func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPost:
		h.post(w, r)
	case http.MethodPut:
		h.put(w, r)
	case http.MethodDelete:
		h.remove(w, r)
	case http.MethodOptions:
		options(w, r)
	default:
		methodNotAllowed(w)
	}
}

// Get is handler for GET requests for non top level Tasks.
func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	taskIDPath, err := parseTaskIDPath(r)
	if err != nil {
		log.Printf("(DEBUG) handler: getting child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, ErrHandlerURLNotValid)
		return
	}

	task, err := h.service.Find(taskIDPath)
	if err != nil {
		switch err {
		case ErrTaskNotFound:
			log.Printf("(INFO) handler: getting child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusNotFound, err)
			return
		default:
			log.Printf("(WARN) handler: getting child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	ResponseOK(w, task)
}

// Put is handler for PUT requests for non top level Tasks.
func (h *TaskHandler) put(w http.ResponseWriter, r *http.Request) {
	taskIDPath, err := parseTaskIDPath(r)
	if err != nil {
		log.Printf("(DEBUG) handler: updating child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, ErrHandlerURLNotValid)
		return
	}

	var jsonTask JSONTask
	if err := parseBody(r, &jsonTask); err != nil {
		log.Printf("(DEBUG) handler: updating child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := jsonTask.Validate(NewUpdateValidator()); err != nil {
		log.Printf("(DEBUG) handler: updating child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	updateFields := UpdateFields{
		Label:     jsonTask.Label,
		Completed: jsonTask.Completed,
	}

	updatedTask, err := h.service.Update(taskIDPath, updateFields)
	if err != nil {
		switch err {
		case ErrTaskNotFound:
			log.Printf("(INFO) handler: updating child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusNotFound, err)
			return
		default:
			log.Printf("(WARN) handler: updating child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	ResponseOK(w, updatedTask)
}

// Post is handler for POST requests for non top level Tasks
func (h *TaskHandler) post(w http.ResponseWriter, r *http.Request) {
	taskIDPath, err := parseTaskIDPath(r)
	if err != nil {
		log.Printf("(DEBUG) handler: creating child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, ErrHandlerURLNotValid)
		return
	}

	var jsonTask JSONTask
	if err := parseBody(r, &jsonTask); err != nil {
		log.Printf("(DEBUG) handler: creating child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := jsonTask.Validate(NewCreateValidator()); err != nil {
		log.Printf("(DEBUG) handler: creating child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	createFields := CreateFields{
		Label: *jsonTask.Label,
	}

	newTask, err := h.service.Create(taskIDPath, createFields)
	if err != nil {
		switch err {
		case ErrTaskNotFound:
			log.Printf("(INFO) handler: creating child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusNotFound, err)
			return
		default:
			log.Printf("(WARN) handler: creating child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	url := fmt.Sprintf("%s/%d", r.URL.Path, newTask.ID)
	ResponseCreated(w, url, newTask)
}

// Remove is handler for DELETE requests for non top level Tasks
func (h *TaskHandler) remove(w http.ResponseWriter, r *http.Request) {
	taskIDPath, err := parseTaskIDPath(r)
	if err != nil {
		log.Printf("(DEBUG) handler: deleting child task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, ErrHandlerURLNotValid)
		return
	}

	task, err := h.service.Delete(taskIDPath)
	if err != nil {
		switch err {
		case ErrTaskNotFound:
			log.Printf("(INFO) handler: deleting child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusNotFound, err)
			return
		default:
			log.Printf("(INFO) handler: deleting child task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	ResponseOK(w, task)
}

// TasksHandler is simple Handler which handles top level Tasks in tree
// hierarchy. Handler provides only CR operations.
// TasksHandler implements http.Handler interface.
type TasksHandler struct {
	service TaskService
}

// NewTasksHandler returns new instance of TasksHandler
func NewTasksHandler(service TaskService) *TasksHandler {
	return &TasksHandler{
		service: service,
	}
}

// ServeHTTP is simple function which dispatches requests to proper function
// handlers.
// ServeHTTP implements http.Handler interface
func (h *TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPost:
		h.post(w, r)
	case http.MethodOptions:
		options(w, r)
	default:
		methodNotAllowed(w)
	}
}

// Get is handler for GET requests for top level Tasks.
func (h *TasksHandler) get(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.FindAll()
	if err != nil {
		switch err {
		case ErrTaskNotFound:
			log.Printf("(INFO) handler: getting task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusNotFound, err)
			return
		default:
			log.Printf("(WARN) handler: getting task failed: %s\n", err)
			ErrorAsJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	sort.Sort(ByTaskID(tasks))

	// Do not return array in response - it would break future extensions
	// Better to return object which wraps tasks.
	response := map[string]interface{}{
		"tasks": tasks,
	}

	ResponseOK(w, response)
}

// Post is handler for POST requests for top level Tasks.
func (h *TasksHandler) post(w http.ResponseWriter, r *http.Request) {
	var jsonTask JSONTask
	if err := parseBody(r, &jsonTask); err != nil {
		log.Printf("(DEBUG) handler: creating task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := jsonTask.Validate(NewCreateValidator()); err != nil {
		log.Printf("(DEBUG) handler: creating task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	createFields := CreateFields{
		Label: *jsonTask.Label,
	}

	// Creating op level Task - TaskID path will always be empty.
	newTask, err := h.service.Create([]TaskID{}, createFields)
	if err != nil {
		log.Printf("(WARN) handler: creating task failed: %s\n", err)
		ErrorAsJSON(w, http.StatusInternalServerError, err)
		return
	}

	url := fmt.Sprintf("%s/%d", r.URL.Path, newTask.ID)
	ResponseCreated(w, url, newTask)
}
