package tasks

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTaskHandler(t *testing.T) {
	tests := map[string]struct {
		method        string
		path          string
		body          io.Reader
		res           string
		resStatusCode int
	}{
		"GET /tasks/2": {
			method:        "GET",
			path:          "/tasks/2",
			res:           `{"id":"2","label":"bar","completed":true}`,
			resStatusCode: 200,
		},
		"POST /tasks/1": {
			method:        "POST",
			path:          "/tasks/1",
			body:          strings.NewReader(`{"label":"foo"}`),
			res:           `{"id":"1","label":"foo","completed":false}`,
			resStatusCode: 201,
		},
		"PUT/tasks/1": {
			method:        "PUT",
			path:          "/tasks/1",
			body:          strings.NewReader(`{"label":"foo_new"}`),
			res:           `{"id":"1","label":"foo_new","completed":false}`,
			resStatusCode: 200,
		},
		"DELETE /tasks/3": {
			method:        "DELETE",
			path:          "/tasks/3",
			res:           `{"id":"3","label":"baz","completed":false}`,
			resStatusCode: 200,
		},
		"OPTIONS /tasks/1": {
			method:        "OPTIONS",
			path:          "/tasks/1",
			res:           "",
			resStatusCode: 200,
		},
		"PATCH /tasks/1": {
			method:        "PATCH",
			path:          "/tasks/1",
			res:           `{"error":"method not allowed"}`,
			resStatusCode: 405,
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		r, err := http.NewRequest(tc.method, fmt.Sprintf("http://foo.com%s", tc.path), tc.body)
		if err != nil {
			t.Fatal(err)
		}
		if tc.method == "POST" || tc.method == "PUT" {
			r.Header.Add("Content-Type", "application/json")
		}

		service := &mockService{}
		handler := NewTaskHandler(service)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if tc.resStatusCode != w.Code {
			t.Fatalf("expected status code %d got %d", tc.resStatusCode, w.Code)
		}

		if tc.res != w.Body.String() {
			t.Fatalf("expected response \n%s\n got \n%s\n", tc.res, w.Body.String())
		}
	}
}

func TestTasksHandler(t *testing.T) {
	tests := map[string]struct {
		method        string
		body          io.Reader
		res           string
		resStatusCode int
	}{
		"GET /tasks": {
			method:        "GET",
			res:           `{"tasks":[{"id":"1","label":"foo","completed":true},{"id":"2","label":"bar","completed":false}]}`,
			resStatusCode: 200,
		},
		"POST /tasks": {
			method:        "POST",
			body:          strings.NewReader(`{"label":"foo"}`),
			res:           `{"id":"1","label":"foo","completed":false}`,
			resStatusCode: 201,
		},
		"OPTIONS /tasks": {
			method:        "OPTIONS",
			res:           "",
			resStatusCode: 200,
		},
		"DELETE /tasks": {
			method:        "DELETE",
			res:           `{"error":"method not allowed"}`,
			resStatusCode: 405,
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		r, err := http.NewRequest(tc.method, "http://foo.com/tasks", tc.body)
		if err != nil {
			t.Fatal(err)
		}
		if tc.method == "POST" {
			r.Header.Add("Content-Type", "application/json")
		}

		service := &mockService{}
		handler := NewTasksHandler(service)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if tc.resStatusCode != w.Code {
			t.Fatalf("expected status code %d got %d", tc.resStatusCode, w.Code)
		}

		if tc.res != w.Body.String() {
			t.Fatalf("expected response \n%s\n got \n%s\n", tc.res, w.Body.String())
		}
	}
}

type mockService struct{}

func (s *mockService) Create(path []TaskID, cf CreateFields) (Task, error) {
	return Task{
		ID:        TaskID(1),
		Label:     "foo",
		Completed: false,
	}, nil
}

func (s *mockService) Find(path []TaskID) (Task, error) {
	return Task{
		ID:        TaskID(2),
		Label:     "bar",
		Completed: true,
	}, nil
}

func (s *mockService) FindAll() ([]Task, error) {
	return []Task{
		Task{
			ID:        TaskID(1),
			Label:     "foo",
			Completed: true,
		},
		Task{
			ID:        TaskID(2),
			Label:     "bar",
			Completed: false,
		},
	}, nil
}

func (s *mockService) Update(path []TaskID, uf UpdateFields) (Task, error) {
	return Task{
		ID:        TaskID(1),
		Label:     *uf.Label,
		Completed: false,
	}, nil
}

func (s *mockService) Delete(path []TaskID) (Task, error) {
	return Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: false,
	}, nil
}
