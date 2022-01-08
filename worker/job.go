package worker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// job status
type Status int

const (
	Running Status = iota
	Paused
	Stopped
)

type ConfigJob struct {
	Name string
	Tag  string
	URL  string
}

type Job struct {
	ConfigJob
	status       Status
	lastResponse Response

	togglePause chan struct{}
	stop        chan struct{}
	Event       chan Status
}

func (job *Job) checkStatus() {
	req, _ := http.NewRequest("GET", job.URL+"api/json?pretty=true", nil)
	req.Header.Set("Authorization", "basic "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	HandleResponse(responseData, job)
	log.Println("Job:", job.Name, "\tRan")
}

func (job *Job) StartCheckStatus() {
	log.Println("Started Job:", job.Name)
	go func() {
		var x time.Duration
		if interval == 0 {
			x = 60
		} else {
			x = time.Duration(interval)
		}
		ticker := time.NewTicker(x * time.Second)

		job.status = Running
		job.checkStatus()
		for {
			select {
			case <-job.togglePause:
				job.status = Paused
				log.Println("Job: ", job.Name, "\tPaused")
				job.Event <- job.status
				select {
				case <-job.togglePause:
					job.status = Running
					log.Println("Job: ", job.Name, "\tUnpaused")
					job.Event <- job.status
				case <-job.stop:
					job.status = Stopped
					log.Println("Job: ", job.Name, "\tStopped")
					job.Event <- job.status
					return
				}
			case <-job.stop:
				job.status = Stopped
				log.Println("Job: ", job.Name, "\tStopped")
				job.Event <- job.status
				return
			case <-ticker.C:
				job.checkStatus()
			}
		}
	}()
}

func (job *Job) setResponse(resp Response) {
	job.lastResponse = Response{
		Name:               resp.Name,
		URL:                resp.URL,
		Builds:             resp.Builds,
		Color:              resp.Color,
		LastBuild:          resp.LastBuild,
		LastCompletedBuild: resp.LastCompletedBuild,
	}
}

func (job *Job) Stop() {
	if job.status != Stopped {
		job.stop <- struct{}{}
	}
}

func (job *Job) TogglePause() {
	job.togglePause <- struct{}{}
}

func (job Job) IsRunning() bool {
	return job.status == Running
}

func (job Job) GetStatus() Status {
	return job.status
}

func (job Job) ToString() {
	fmt.Println("Job.name:", job.Name, ", Job.tag:", job.Tag, ", Job.url:", job.URL)
}
