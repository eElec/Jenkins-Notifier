package worker

import (
	"gopkg.in/toast.v1"
)

func pushNotification(title string, message string, color string, url string) {

	/* if color == "blue_anime" {
		notify = "Build In Progress"
	} else {
		notify = "Build Completed with color = " + color
		var icon = ""
		if color == "blue" {
			icon = `C:\Users\aayushi.a\Desktop\greenTick.png`
		} else if color == "yellow" {
			icon = `C:\Users\aayushi.a\Desktop\yellowExclamation.png`
		} else if color == "red" {
			icon = `C:\Users\aayushi.a\Desktop\redCross.jpg`
		} else if color == "aborted" {
			icon = `C:\Users\aayushi.a\Desktop\blackAborted.png`
		}
	*/

	notification := toast.Notification{
		AppID:  "Jenkins Notifier",
		Title:   title,
		Message: message,
		//Icon:    "go.png", // This file must exist (remove this line if it doesn't)
		Actions: []toast.Action{
			{Type: "protocol", Label: "View Status", Arguments: url},
			{Type: "protocol", Label: "Okay", Arguments: ""},
		},
	}
	err := notification.Push()
	if err != nil {
		// handle
	}
}
