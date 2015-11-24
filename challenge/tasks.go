package challenge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	TaskClusterURIPrefix = "https://scheduler.taskcluster.net/v1"
)

type TaskGraphStatus struct {
	TaskGraphID string `json:"taskGraphId"`
	SchedulerID string `json:"schedulerId"`
	State       string `json:"state"`
}

type TaskInformation struct {
	TaskID       string   `json:"taskId"`
	Name         string   `json:"name"`
	Requires     []string `json:"requires"`
	RequiresLeft []string `json:"requiresLeft"`
	Reruns       int      `json:"reruns"`
	RerunsLeft   int      `json:"rerunsLeft"`
	State        string   `json:"state"`
	Satisfied    bool     `json:"satisfied"`
	Dependents   []string `json:"dependents"`
}

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Source      string `json:"source"`
}

type TaskGraphResponse struct {
	Status   TaskGraphStatus   `json:"status"`
	Tasks    []TaskInformation `json:"tasks"`
	Metadata Metadata          `json:"metadata"`
	Tags     map[string]string `json:"tags"`
	Scopes   []string          `json:"scopes"`
}

// ChallengeResponse maps task IDs to their status.
type ChallengeResponse struct {
	TaskID string `json:"taskId"`
	State  string `json:"state"`
}

// fetchRawJSON returns the JSON returned by GETting the proper URI.
func fetchRawJSON(uri string) ([]byte, error) {
	result, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// jsonDecode returns a decoded JSON response
func jsonDecode(body []byte) (TaskGraphResponse, error) {
	var tg TaskGraphResponse
	err := json.Unmarshal(body, &tg)
	return tg, err
}

// DoGetTasksState returns a JSON representation of the state of all
// tasks in the TaskGraph.
func DoGetTasksState(id string) ([]ChallengeResponse, error) {
	uri := strings.Join([]string{TaskClusterURIPrefix, "/task-graph/", id, "/inspect"}, "")
	body, err := fetchRawJSON(uri)
	if err != nil {
		return nil, err
	}

	tgResp, err := jsonDecode(body)
	if err != nil {
		return nil, err
	}

	chResp := make([]ChallengeResponse, len(tgResp.Tasks))
	for i, task := range tgResp.Tasks {
		chResp[i].TaskID = task.TaskID
		chResp[i].State = task.State
	}
	return chResp, nil
}

// htmlize returns HTML instead of bare JSON.
func htmlize(in []ChallengeResponse) []byte {
	var buf bytes.Buffer
	for _, ch := range in {
		buf.WriteString(fmt.Sprintf("task %s: state %s<br>", ch.TaskID, ch.State))
	}
	return buf.Bytes()
}

// GetTasksState handles HTTP responses for GET /blah
func GetTasksState(w http.ResponseWriter, r *http.Request) {
	// get prout from DoGetTasksState
	taskID := r.URL.Path[len("/task-state/"):]
	taskStates, err := DoGetTasksState(taskID)
	if err != nil {
		http.Error(w, "Error writing task status", 500)
		return
	}
	if r.FormValue("output") == "html" {
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlize(taskStates))
	} else {
		out, err := json.Marshal(taskStates)
		if err != nil {
			http.Error(w, "Error writing tasks status", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}
