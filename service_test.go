package tasks

import "testing"

func TestTaskServiceCreate(t *testing.T) {
	t1 := &Task{
		ID:        TaskID(1),
		Label:     "foo",
		Completed: false,
	}
	t2 := &Task{
		ID:        TaskID(2),
		Label:     "bar",
		Completed: false,
	}
	t3 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: false,
	}

	tests := map[string]struct {
		path       []TaskID
		fields     CreateFields
		storage    map[TaskID]*Task
		lastTaskID TaskID
		res        Task
		err        error
	}{
		"empty path": {
			path: []TaskID{},
			fields: CreateFields{
				Label: "foo",
			},
			storage: map[TaskID]*Task{},
			res:     *t1,
		},
		"first level": {
			path: []TaskID{
				TaskID(2),
			},
			fields: CreateFields{
				Label: "baz",
			},
			storage: map[TaskID]*Task{
				TaskID(2): t2,
			},
			lastTaskID: TaskID(2),
			res:        *t3,
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		storage := NewTaskMemoryStorage()
		storage.storage = tc.storage
		storage.lastTaskID = tc.lastTaskID
		service := NewTaskStorageService(storage)

		res, err := service.Create(tc.path, tc.fields)
		if err != tc.err {
			t.Fatalf("expected err %s got %s", tc.err, err)
		}

		if err != nil {
			continue
		}

		if tc.res.ID != res.ID {
			t.Fatalf("expected id %d got %d", tc.res.ID, res.ID)
		}

		if tc.res.Label != res.Label {
			t.Fatalf("expected label %s got %s", tc.res.Label, res.Label)
		}

		if tc.res.Completed != res.Completed {
			t.Fatalf("expected completed %t got %t", tc.res.Completed, res.Completed)
		}
	}
}

func TestTaskServiceFind(t *testing.T) {
	t.Skip("No business logic")
}

func TestTaskServiceFindAll(t *testing.T) {
	t.Skip("No business logic")
}

func TestTaskServiceUpdate(t *testing.T) {
	t1 := &Task{
		ID:        TaskID(1),
		Label:     "foo",
		Completed: true,
	}
	t2 := &Task{
		ID:        TaskID(2),
		Label:     "bar",
		Completed: false,
	}

	foo := "foo"
	bar := "bar"
	completed := true
	notCompleted := false

	tests := map[string]struct {
		path    []TaskID
		fields  UpdateFields
		storage map[TaskID]*Task
		res     Task
		err     error
	}{
		"empty path": {
			path:    []TaskID{},
			storage: map[TaskID]*Task{},
			err:     ErrTaskPathNotValid,
		},
		"first level": {
			path: []TaskID{
				TaskID(1),
			},
			fields: UpdateFields{
				Label:     &foo,
				Completed: &completed,
			},
			storage: map[TaskID]*Task{
				TaskID(1): &Task{
					ID:        TaskID(1),
					Label:     "foo_old",
					Completed: false,
				},
			},
			res: *t1,
		},
		"second level": {
			path: []TaskID{
				TaskID(3),
				TaskID(2),
			},
			fields: UpdateFields{
				Label:     &bar,
				Completed: &notCompleted,
			},
			storage: map[TaskID]*Task{
				TaskID(3): &Task{
					ID:        TaskID(3),
					Label:     "baz",
					Completed: true,
					Children: map[TaskID]*Task{
						t2.ID: &Task{
							ID:        TaskID(2),
							Label:     "bar_old",
							Completed: true,
						},
					},
				},
			},
			res: *t2,
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		storage := NewTaskMemoryStorage()
		storage.storage = tc.storage
		service := NewTaskStorageService(storage)

		res, err := service.Update(tc.path, tc.fields)
		if err != tc.err {
			t.Fatalf("expected err %s got %s", tc.err, err)
		}

		if err != nil {
			continue
		}

		if tc.res.ID != res.ID {
			t.Fatalf("expected id %d got %d", tc.res.ID, res.ID)
		}

		if tc.res.Label != res.Label {
			t.Fatalf("expected label %s got %s", tc.res.Label, res.Label)
		}

		if tc.res.Completed != res.Completed {
			t.Fatalf("expected completed %t got %t", tc.res.Completed, res.Completed)
		}
	}
}

func TestTaskServiceDelete(t *testing.T) {
	t.Skip("No business logic")
}
