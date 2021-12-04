package worker

import (
	"encoding/json"
)

type Build struct {
	Class string 
	Number int 
	URL string
}

type Response struct {
	Name string `json:"displayName"`
	URL string `json:"url"`
	Builds []Build `json:"builds"`
	Color string `json:"color"`
	LastBuild Build `json:"lastBuild"`
	LastCompletedBuild Build `json:"lastCompletedBuild"`
}


func HandleResponse(resp []byte) {

	var responseObject Response
	json.Unmarshal(resp, &responseObject)

	var color = responseObject.Color
	var notify string
		
	if color == "blue_anime" {
		notify = "Build In Progress"
	} else {
		notify = "Build Completed with color = " + color
		pushNotification(responseObject.Name, notify, color, responseObject.URL)
	}
}