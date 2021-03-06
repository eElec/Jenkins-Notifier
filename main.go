package main

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"net/http"
	"os"
	"syscall"

	"jenkins-notifier/utils"
	"jenkins-notifier/worker"

	"github.com/getlantern/systray"
	"golang.org/x/term"
)

var (
	client *http.Client = &http.Client{}
	jobs   map[string]*worker.Job
)

const (
	AUTH_TOKEN    = "token"
	AUTH_PASSWORD = "password"
)

type Configuration struct {
	Jobs     []worker.ConfigJob
	Interval uint
	AuthType string
}

func loadConfig(path string, ret interface{}) error {
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(ret)
	return err
}

/*
	Basic Authenticate using either password or token

	returns base64(user:password) or base64(user:token)
*/
func getCredentials(authType string) (string, error) {
	var (
		username string
		password string
	)

	var token struct {
		Username string
		ApiKey   string
	}
	err := loadConfig("config/keys.json", &token)
	if err != nil {
		log.Fatalln("Error reading key.json")
	}

	username = token.Username
	if authType == AUTH_TOKEN {
		password = token.ApiKey
		if password == "" {
			log.Fatalln("authType is set to \"token\", but no token is found in keys.json")
		}
	} else {
		log.Print("\nPassword: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalln(("Error while reading password."))
			os.Exit(1)
		}
		password = string(bytePassword)
	}

	tokenEncoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	return tokenEncoded, nil
}

func main() {
	// Load configurations
	var config Configuration
	err := loadConfig("config/config.json", &config)
	if err != nil {
		log.Fatalln("Error reading config.json")
	}

	// Get auth token
	basicAuthToken, err := getCredentials(config.AuthType)
	if err != nil {
		os.Exit(1)
	}

	jobs = worker.Initialize(client, basicAuthToken, config.Interval, config.Jobs)
	for k, _ := range jobs {
		jobs[k].StartCheckStatus()
	}
	systray.Run(onReady, onExit)
}

func onExit() {
	for k, _ := range jobs {
		jobs[k].Stop()
	}
}

func onReady() {
	systray.SetIcon(utils.GetIcon("icons/jenkins.ico"))
	systray.SetTitle("Jenkins Notifier")
	systray.SetTooltip("Jenkins Notifier")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		log.Println("Quiting...")
		systray.Quit()
	}()

	// job menu
	jobsMenu := systray.AddMenuItem("Jobs", "")

	for k, _ := range jobs {
		subMenu := jobsMenu.AddSubMenuItemCheckbox(jobs[k].Name, "", jobs[k].IsRunning())

		//
		go func(subMenu *systray.MenuItem, job *worker.Job) {
			for {
				select {
				case <-subMenu.ClickedCh:
					job.TogglePause()
					if subMenu.Checked() {
						subMenu.Uncheck()
					} else {
						subMenu.Check()
					}
				}
			}
		}(subMenu, jobs[k])
	}
}
