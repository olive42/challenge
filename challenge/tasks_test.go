package challenge

import "testing"

func TestDoGetTasksState(t *testing.T) {
	testCases := []struct {
		id   string
		want string
	}{
		{
			id:   "A0NVn8c9TVSqGStClwmsrw",
			want: "",
		},
	}
	for _, c := range testCases {
		// FIXME(olive): this is making an external call; mocking fetchRawJSON is an answer.
		got, err := DoGetTasksState(c.id)
		if err != nil {
			t.Errorf("DoGetTasksState(%s): %#v", c.id, err)
		}
		if string(got) != c.want {
			t.Errorf("DoGetTasksState(%s): got %s, expected %s", c.id, got, c.want)
		}
	}
}

func TestDoGetTasksStateError(t *testing.T) {
	testCases := []struct {
		id   string
		want string
	}{
		{
			id:   "pasglop",
			want: "",
		},
		// TODO(olive): ideally, taskcluster.net would support erroneous TaskIDs
	}
	for _, c := range testCases {
		_, err := DoGetTasksState(c.id)
		if err == nil {
			t.Errorf("DoGetTasksState(%s): did not error", c.id)
		} else if got := err.Error(); got != c.want {
			t.Errorf("DoGetTasksState(%s):\ngot error: %#v\nexpected error: %#v")
		}
	}
}
