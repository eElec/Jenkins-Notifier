package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	client   *http.Client
	token    string
	interval uint

	iconSuccess  string
	iconUnstable string
	iconFailure  string
)

var JobStates []Response

func Initialize(cl *http.Client, authToken string, intervalArg uint, configJob []ConfigJob) (jobs map[string]*Job) {
	client = cl
	token = authToken
	interval = intervalArg

	pwd, _ := os.Getwd()
	iconSuccess = filepath.Join(pwd, "icons", "success.png")
	iconUnstable = filepath.Join(pwd, "icons", "unstable.png")
	iconFailure = filepath.Join(pwd, "icons", "fail.png")

	jobs = make(map[string]*Job)

	for _, job := range configJob {
		jobs[job.Name] = &Job{
			ConfigJob:   job,
			status:      Stopped,
			togglePause: make(chan struct{}),
			stop:        make(chan struct{}),
			Event:       make(chan Status),
		}
	}
	return jobs
}

func HandleResponse(resp []byte, job *Job) {
	var responseObject Response
	json.Unmarshal(resp, &responseObject)
	if job.lastResponse.Color != "" {
		var color = responseObject.Color
		var number = responseObject.LastBuild.Number
		log.Println(strconv.Itoa(number) + ":\t" + color)
		if job.lastResponse.Color != color || job.lastResponse.LastBuild.Number != number {
			prepareNotification(*job, color, number, responseObject)
			job.setResponse(responseObject)
		}
	} else {
		job.setResponse(responseObject)
	}
}

func prepareNotification(job Job, color string, buildNumber int, resp Response) {
	var message = "Build #" + strconv.Itoa(buildNumber)

	// Build started
	if strings.Contains(color, "_anime") {
		message += "\nBuild Started"
	} else {
		message += "\nBuild Completed"
	}

	// Icons
	var icon string
	switch {
	case strings.Contains(color, "blue"):
		icon = iconSuccess
	case strings.Contains(color, "red"):
		icon = iconFailure
	case strings.Contains(color, "yellow"):
		icon = iconUnstable
	// case strings.Contains(color, "aborted"):
	//
	default:
		icon = ""
	}

	var url = resp.URL + "/" + strconv.Itoa(buildNumber)
	pushNotification(job.Name, message, icon, url)
}
