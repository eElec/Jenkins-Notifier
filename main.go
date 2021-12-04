package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"os"
	"syscall"

	"golang.org/x/term"
	"jenkins-notifier/worker"

	"time"
)

var (
	client *http.Client = &http.Client{}
)

const (
	AUTH_TOKEN    = "token"
	AUTH_PASSWORD = "password"
)

type Configuration struct {
	Jobs []struct {
		Name string
		Tag  string
		URL  string
	}
	AuthType string
}

func main() {
	// notifier.PushNotification()
	// Load configurations
	var config Configuration
	err := loadConfig("config/config.json", &config)
	if err != nil {
		log.Fatalln("Error reading config.json")
	}

	// Get username and password
	basicAuthToken, err := getCredentials(config.AuthType)
	if err != nil {
		os.Exit(1)
	}

	tick := time.Tick(10 * time.Second)
	for range tick {
		for i := 0; i < len(config.Jobs); i++ {
			req, _ := http.NewRequest("GET", config.Jobs[i].URL+"api/json?pretty=true", nil)
			req.Header.Set("Authorization", "basic "+ basicAuthToken)
			response, err := client.Do(req)
	
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}
			worker.HandleResponse(responseData, i)
		}
	}
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
		fmt.Print("\nPassword: ")
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
