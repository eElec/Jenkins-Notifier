package worker

import (
	"fmt"
	// "io/ioutil"
	"log"
	// "net/http"
	"time"
)

// job status
const (
	running = iota
	paused  = iota
	stopped = iota
)

type ConfigJob struct {
	Name string
	Tag  string
	URL  string
}

type Job struct {
	ConfigJob
	status       int
	lastResponse Response

	togglePause chan struct{}
	stop        chan struct{}
}

func (job Job) checkStatus() {
	// req, _ := http.NewRequest("GET", job.URL+"api/json?pretty=true", nil)
	// req.Header.Set("Authorization", "basic "+token)
	// resp, err := client.Do(req)
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer resp.Body.Close()

	// responseData, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// }
	log.Println("Job:", job.Name, "\tRan")
}

func (job Job) StartCheckStatus() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-job.togglePause:
			job.status = paused
			fmt.Println("Job: ", job.Name, "\tPaused")
			select {
			case <-job.togglePause:
			}
		case <-job.stop:
			job.status = stopped
			fmt.Println("Job: ", job.Name, "\tStopped")
			// wg.Done
			return
		case <-ticker.C:
			job.status = running
			job.checkStatus()
		}
	}
}

func (job Job) Stop() {
	if job.status != stopped {
		job.stop <- struct{}{}
	}
}

func (job Job) TogglePause() {
	job.togglePause <- struct{}{}
}

func (job Job) IsRunning() bool {
	return job.status == running
}

func (job Job) ToString() {
	fmt.Println("Job.name:", job.Name, ", Job.tag:", job.Tag, ", Job.url:", job.URL)
}
