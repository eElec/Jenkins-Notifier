package worker

import (
	"gopkg.in/toast.v1"
)

func pushNotification(title string, message string, icon string, url string) {
	println(icon)
	notification := toast.Notification{
		AppID:  "Jenkins Notifier",
		Title:   title,
		Message: message,
		Icon:    icon,
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
