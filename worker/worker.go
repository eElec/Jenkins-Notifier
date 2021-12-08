package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Build struct {
	Class  string
	Number int
	URL    string
}
type Response struct {
	Name               string  `json:"displayName"`
	URL                string  `json:"url"`
	Builds             []Build `json:"builds"`
	Color              string  `json:"color"`
	LastBuild          Build   `json:"lastBuild"`
	LastCompletedBuild Build   `json:"lastCompletedBuild"`
}

var (
	client *http.Client
	token  string
)

var JobStates []Response

func Initialize(cl *http.Client, authToken string, configJob []ConfigJob) (jobs map[string]Job) {
	client = cl
	token = authToken

	jobs = make(map[string]Job)

	for _, job := range configJob {
		jobs[job.Name] = Job{
			ConfigJob:   job,
			status:      stopped,
			togglePause: make(chan struct{}),
			stop:        make(chan struct{}),
		}
	}
	return jobs
}

func HandleResponse(resp []byte, index int) {
	var responseObject Response
	json.Unmarshal(resp, &responseObject)

	if index < len(JobStates) && JobStates[index].Name == responseObject.Name {
		var color = responseObject.Color
		//if (JobStates[index].Color != color) {
		var notify string

		if color == "blue_anime" {
			notify = "Build In Progress"
		} else {
			notify = "Build Completed with color = " + color
		}
		fmt.Print("-----> " + color + "\n")
		pushNotification(responseObject.Name, notify, color, responseObject.URL)
		//}
	} else {
		fmt.Print("not Exist\n")
		JobStates = append(JobStates, responseObject)
	}
}
