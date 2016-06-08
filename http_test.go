package tasks

import (
	"net/http"
	"testing"
)

func TestParseTaskIDPath(t *testing.T) {
	tests := map[string]struct {
		path string
		res  []TaskID
		err  error
	}{
		"/tasks": {
			path: "http://foo.com/tasks",
			res:  []TaskID{},
		},
		"/tasks/1": {
			path: "http://foo.com/tasks/1",
			res: []TaskID{
				TaskID(1),
			},
		},
		"/tasks/1/2": {
			path: "http://foo.com/tasks/1/2",
			res: []TaskID{
				TaskID(1),
				TaskID(2),
			},
		},
		"/tasks/kekeke": {
			path: "http://foo.com/tasks/kekeke",
			err:  ErrHandlerURLNotValid,
		},
		"/": {
			path: "http://foo.com/",
			err:  ErrHandlerURLNotValid,
		},
		"no slash": {
			path: "http://foo.com",
			err:  ErrHandlerURLNotValid,
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		r, err := http.NewRequest("GET", tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := parseTaskIDPath(r)
		if err != tc.err {
			t.Fatalf("expected err %s got %s", tc.err, err)
		}

		if err != nil {
			continue
		}

		if len(tc.res) != len(res) {
			t.Fatalf("expected taskid path len %d got %d", len(tc.res), len(res))
		}

		for i, taskID := range res {
			if taskID != tc.res[i] {
				t.Fatalf("expected taskid path %d got %d", tc.res[i], taskID)
			}
		}
	}
}
