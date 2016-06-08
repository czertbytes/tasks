package tasks

import "testing"

func TestTaskValidator(t *testing.T) {
	empty := ""
	tooLong := "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123"
	label := "foobar"
	completed := true

	tests := map[string]struct {
		validator TaskActionValidator
		jsonTask  *JSONTask
		err       error
	}{
		"create label nil": {
			validator: NewCreateValidator(),
			jsonTask:  &JSONTask{},
			err:       ErrTaskLabelIsRequired,
		},
		"create label empty": {
			validator: NewCreateValidator(),
			jsonTask: &JSONTask{
				Label: &empty,
			},
			err: ErrTaskLabelIsNotValid,
		},
		"create label too long": {
			validator: NewCreateValidator(),
			jsonTask: &JSONTask{
				Label: &tooLong,
			},
			err: ErrTaskLabelIsNotValid,
		},
		"create pass": {
			validator: NewCreateValidator(),
			jsonTask: &JSONTask{
				Label: &label,
			},
		},
		"update both nil": {
			validator: NewUpdateValidator(),
			jsonTask:  &JSONTask{},
			err:       ErrTaskLabelOrCompletedRequired,
		},
		"update label nil": {
			validator: NewUpdateValidator(),
			jsonTask: &JSONTask{
				Completed: &completed,
			},
		},
		"update completed nil": {
			validator: NewUpdateValidator(),
			jsonTask: &JSONTask{
				Label: &label,
			},
		},
		"update label empty": {
			validator: NewUpdateValidator(),
			jsonTask: &JSONTask{
				Label: &empty,
			},
			err: ErrTaskLabelIsNotValid,
		},
		"update label too long": {
			validator: NewUpdateValidator(),
			jsonTask: &JSONTask{
				Label: &tooLong,
			},
			err: ErrTaskLabelIsNotValid,
		},
	}

	for desc, tc := range tests {
		t.Log(desc)

		err := tc.validator.Validate(tc.jsonTask)
		if err != tc.err {
			t.Fatalf("expected err %s got %s", tc.err, err)
		}

		if err != nil {
			continue
		}
	}
}
