package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

var (
	device    string
	message   string
	priority  string
	sound     string
	timestamp string
	title     string
	tokenKey  string
	urlPath   string
	urlTitle  string
	userKey   string
)

var verbose bool
var version bool

// Semantic versioning - http://semver.org/
const (
	Major = 1
	Minor = 0
	Patch = 0
)

// init process command line flags and read in user created api keys file.
func init() {
	flag.BoolVar(&version, "version", false, "Display program version.")

	flag.BoolVar(&verbose, "v", false, "Display message response details.")

	flag.StringVar(&device, "device", "", "Send message directly to this "+
		"device, rather than all devices.")

	flag.StringVar(&message, "msg", "", "(Required) - Your message.")

	flag.StringVar(&priority, "priority", "", "Message priority. -2 = lowest "+
		", -1 = low, 0 = normal, 1 = high, 2 = emergency.")

	flag.StringVar(&sound, "sound", "", "Sound to play when user receives "+
		"notification, overrides the user's default sound choice.")

	flag.StringVar(&timestamp, "timestamp", "", "A Unix timestamp of your "+
		"message's date and time to display to the user.")

	flag.StringVar(&title, "title", "", "Your message's title, otherwise "+
		"your app's name is used.")

	flag.StringVar(&urlPath, "url", "", "A supplementary URL to show with "+
		"your message.")

	flag.StringVar(&urlTitle, "url-title", "", "A title for your "+
		"supplementary URL, otherwise just the URL is shown.")

	flag.Parse()

	if version {
		fmt.Println("Pushover", fmt.Sprintf("%d.%d.%d", Major, Minor, Patch))
		os.Exit(0)
	}

	if message == "" {
		fmt.Println("-msg is a required flag. You must specify a message, " +
			"use -h for help.")
		os.Exit(1)
	}

	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(u.HomeDir + "/.pushover")
	if err != nil {
		fmt.Println(err)
		fmt.Println("You need to create a dot file '.pushover' in your home " +
			"directory, containing your token and user api keys, in the " +
			"format:\ntoken=XXXXXXXXXX\nuser=XXXXXXXXXX")
		os.Exit(1)
	}

	// read in user created api keys file.
	r := bufio.NewReaderSize(file, 50)
	for l, err := r.ReadString('\n'); err == nil; l, err = r.ReadString('\n') {
		t := strings.Split(l, "=")
		if t[0] == "token" {
			tokenKey = strings.TrimSpace(t[1])
		} else if t[0] == "user" {
			userKey = strings.TrimSpace(t[1])
		}
	}

	if tokenKey == "" || userKey == "" {
		fmt.Println("Error processing 'token'/'user' from .pushover file.")
		os.Exit(1)
	}
}

// extractJSON extracts and returns a JSON structure.
func extractJSON(jsonBlob []byte) (map[string]interface{}, error) {
	var v interface{}

	err := json.Unmarshal(jsonBlob, &v)
	if err != nil {
		return nil, err
	}

	j := v.(map[string]interface{})

	return j, nil
}

func main() {
	// request payload.
	v := url.Values{}
	v.Set("token", tokenKey)
	v.Set("user", userKey)
	v.Set("message", message)

	if device != "" {
		v.Set("device", device)
	}

	if priority != "" {
		v.Set("priority", priority)
		if priority == "2" {
			// emergency priority notification.
			v.Set("retry", "60")
			v.Set("expire", "3600")
		}
	}

	if sound != "" {
		v.Set("sound", sound)
	}

	if timestamp != "" {
		v.Set("timestamp", timestamp)
	}

	if title != "" {
		v.Set("title", title)
	}

	if urlPath != "" {
		v.Set("url", urlPath)
	}

	if urlTitle != "" {
		v.Set("urlTitle", urlTitle)
	}

	p := v.Encode()

	req, err := http.NewRequest("POST",
		"https://api.pushover.net/1/messages.json",
		bytes.NewBufferString(p))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Host", "api.pushover.net")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(p)))

	c := &http.Client{}
	rsp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if rsp.StatusCode == http.StatusOK {
		j, err := extractJSON(body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Success")
		if verbose {
			fmt.Println("Request id:", j["request"])
			fmt.Println("Message limit:", rsp.Header["X-Limit-App-Limit"][0])
			fmt.Println("Remaining Messages:", rsp.Header["X-Limit-App-Remaining"][0])

			n, err := strconv.ParseInt(rsp.Header["X-Limit-App-Reset"][0], 10, 64)
			if err != nil {
				fmt.Println("Message limit resets on: ?")
			} else {
				t := time.Unix(n, 0)
				fmt.Println("Message limit resets on:", t.Format("Jan 2, 2006 at 3:04pm (UTC)"))
			}
		}
	} else if rsp.StatusCode >= 400 && rsp.StatusCode <= 499 {
		j, err := extractJSON(body)
		if err != nil {
			fmt.Println(err)
		} else {
			m := j["errors"].(interface{})
			fmt.Println("Fail", m)
		}
		os.Exit(1)
	} else {
		fmt.Println(http.StatusText(rsp.StatusCode))
		os.Exit(1)
	}
}
