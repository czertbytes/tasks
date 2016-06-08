package tasks

import (
	"reflect"
	"testing"
)

func TestTaskMemoryStorageInsert(t *testing.T) {
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
	t3 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: true,
		Children: map[TaskID]*Task{
			t2.ID: t2,
		},
	}

	tests := map[string]struct {
		path       []TaskID
		task       *Task
		err        error
		storage    map[TaskID]*Task
		expStorage map[TaskID]*Task
	}{
		"empty path": {
			path:    []TaskID{},
			task:    t1,
			storage: map[TaskID]*Task{},
			expStorage: map[TaskID]*Task{
				t1.ID: t1,
			},
		},
		"first level not found": {
			path:    []TaskID{t2.ID},
			task:    t1,
			storage: map[TaskID]*Task{},
			err:     ErrTaskNotFound,
		},
		"first level found": {
			path: []TaskID{t3.ID},
			task: t2,
			storage: map[TaskID]*Task{
				t3.ID: t3,
			},
			expStorage: map[TaskID]*Task{
				t3.ID: t3,
			},
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		storage := NewTaskMemoryStorage()
		storage.storage = tc.storage

		err := storage.Insert(tc.path, tc.task)
		if err != tc.err {
			t.Fatalf("expected err %s got %s", tc.err, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(tc.expStorage, storage.storage) {
			t.Fatalf("expected storage %v got %v", tc.expStorage, storage.storage)
		}
	}
}

func TestTaskMemoryStorageFind(t *testing.T) {
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
	t3 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: true,
		Children: map[TaskID]*Task{
			t2.ID: t2,
		},
	}

	tests := map[string]struct {
		path    []TaskID
		res     Task
		err     error
		storage map[TaskID]*Task
	}{
		"empty path": {
			path: []TaskID{},
			err:  ErrTaskPathNotValid,
		},
		"first level": {
			path: []TaskID{
				TaskID(1),
			},
			res: *t1,
			storage: map[TaskID]*Task{
				TaskID(1): t1,
			},
		},
		"first level not found": {
			path: []TaskID{
				TaskID(2),
			},
			err: ErrTaskNotFound,
			storage: map[TaskID]*Task{
				TaskID(1): t1,
			},
		},
		"second level": {
			path: []TaskID{
				TaskID(3),
				TaskID(2),
			},
			res: *t2,
			storage: map[TaskID]*Task{
				t3.ID: t3,
			},
		},
		"second level not found": {
			path: []TaskID{
				TaskID(3),
				TaskID(1),
			},
			err: ErrTaskNotFound,
			storage: map[TaskID]*Task{
				t3.ID: t3,
			},
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		storage := NewTaskMemoryStorage()
		storage.storage = tc.storage

		res, err := storage.Find(tc.path)
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

		if !reflect.DeepEqual(tc.res.Children, res.Children) {
			t.Fatalf("expected children %v got %v", tc.res.Children, res.Children)
		}
	}
}

func TestTaskMemoryStorageUpdate(t *testing.T) {
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
	t3 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: true,
		Children: map[TaskID]*Task{
			t2.ID: t2,
		},
	}
	t4 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: true,
		Children: map[TaskID]*Task{
			t2.ID: t1,
		},
	}

	tests := map[string]struct {
		path       []TaskID
		task       *Task
		err        error
		storage    map[TaskID]*Task
		expStorage map[TaskID]*Task
	}{
		"empty path": {
			path: []TaskID{},
			err:  ErrTaskPathNotValid,
			storage: map[TaskID]*Task{
				t1.ID: t2,
			},
		},
		"first level not found": {
			path: []TaskID{t1.ID},
			task: t1,
			err:  ErrTaskNotFound,
			storage: map[TaskID]*Task{
				t2.ID: t2,
			},
		},
		"first level found": {
			path: []TaskID{t1.ID},
			task: t1,
			storage: map[TaskID]*Task{
				t1.ID: t2,
			},
			expStorage: map[TaskID]*Task{
				t1.ID: t1,
			},
		},
		"second level not found": {
			path:    []TaskID{t2.ID, t1.ID},
			task:    t1,
			err:     ErrTaskNotFound,
			storage: map[TaskID]*Task{},
		},
		"second level found": {
			path: []TaskID{t3.ID, t2.ID},
			task: t2,
			storage: map[TaskID]*Task{
				t3.ID: t4,
			},
			expStorage: map[TaskID]*Task{
				t3.ID: t3,
			},
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		storage := NewTaskMemoryStorage()
		storage.storage = tc.storage

		err := storage.Update(tc.path, tc.task)
		if err != tc.err {
			t.Fatalf("expected err %s got %s", tc.err, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(tc.expStorage, storage.storage) {
			t.Fatalf("expected storage %v got %v", tc.expStorage, storage.storage)
		}
	}
}

func TestTaskMemoryStorageDelete(t *testing.T) {
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
	t3 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: true,
		Children: map[TaskID]*Task{
			t2.ID: t2,
		},
	}
	t4 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: true,
		Children:  map[TaskID]*Task{},
	}

	tests := map[string]struct {
		path       []TaskID
		err        error
		storage    map[TaskID]*Task
		expStorage map[TaskID]*Task
	}{
		"empty path": {
			path: []TaskID{},
			err:  ErrTaskPathNotValid,
			storage: map[TaskID]*Task{
				t1.ID: t2,
			},
		},
		"first level not found": {
			path: []TaskID{t1.ID},
			err:  ErrTaskNotFound,
			storage: map[TaskID]*Task{
				t2.ID: t1,
			},
		},
		"first level found": {
			path: []TaskID{t1.ID},
			storage: map[TaskID]*Task{
				t1.ID: t1,
			},
			expStorage: map[TaskID]*Task{},
		},
		"second level not found": {
			path:    []TaskID{t2.ID, t1.ID},
			storage: map[TaskID]*Task{},
			err:     ErrTaskNotFound,
		},
		"second level found": {
			path: []TaskID{t3.ID, t2.ID},
			storage: map[TaskID]*Task{
				t3.ID: t3,
			},
			expStorage: map[TaskID]*Task{
				t3.ID: t4,
			},
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		storage := NewTaskMemoryStorage()
		storage.storage = tc.storage

		err := storage.Delete(tc.path)
		if err != tc.err {
			t.Fatalf("expected err %s got %s", tc.err, err)
		}

		if err != nil {
			continue
		}

		if !reflect.DeepEqual(tc.expStorage, storage.storage) {
			t.Fatalf("expected storage %v got %v", tc.expStorage, storage.storage)
		}
	}
}

func TestTaskMemoryStorageSearch(t *testing.T) {
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
	t3 := &Task{
		ID:        TaskID(3),
		Label:     "baz",
		Completed: true,
		Children: map[TaskID]*Task{
			t2.ID: t2,
		},
	}

	tests := map[string]struct {
		path    []TaskID
		res     *Task
		err     error
		storage map[TaskID]*Task
	}{
		"empty path": {
			path: []TaskID{},
			err:  ErrTaskPathNotValid,
		},
		"first level": {
			path: []TaskID{
				TaskID(1),
			},
			res: t1,
			storage: map[TaskID]*Task{
				TaskID(1): t1,
			},
		},
		"first level not found": {
			path: []TaskID{
				TaskID(2),
			},
			err: ErrTaskNotFound,
			storage: map[TaskID]*Task{
				TaskID(1): t1,
			},
		},
		"second level": {
			path: []TaskID{
				TaskID(3),
				TaskID(2),
			},
			res: t2,
			storage: map[TaskID]*Task{
				t3.ID: t3,
			},
		},
		"second level not found": {
			path: []TaskID{
				TaskID(3),
				TaskID(1),
			},
			err: ErrTaskNotFound,
			storage: map[TaskID]*Task{
				t3.ID: t3,
			},
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		storage := NewTaskMemoryStorage()
		storage.storage = tc.storage

		res, err := storage.search(tc.path)
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

		if !reflect.DeepEqual(tc.res.Children, res.Children) {
			t.Fatalf("expected children %v got %v", tc.res.Children, res.Children)
		}
	}
}
