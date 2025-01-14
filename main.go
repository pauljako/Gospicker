package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Json struct {
	Text        string `json:"text,omitempty"`
	Description string `json:"description,omitempty"`
	LinkText    string `json:"linktext,omitempty"`
	Buttons     []struct {
		Text     string `json:"text,omitempty"`
		NextPath string `json:"nextPath,omitempty"`
	} `json:"buttons,omitempty"`
}

func main() {
	var startUrl []string = []string{"https://kleefuchs.github.io/OS-Picker/assets"}
	err := initialize(startUrl)
	if err != nil {
		log.Fatalln(err)
	}
}

func initialize(url []string) error {
	json, err := getJson(strings.Join(append(url, "index.json"), "/"))
	if err != nil {
		return err
	}
	if json.Text != "" {
		fmt.Println("\n" + json.Text)
	}
	if json.Description != "" {
		fmt.Println("\n" + json.Description)
	}
	if len(url) > 1 {
		fmt.Printf("\n%v) %v", 0, "Back")
	}
	if len(json.Buttons) > 0 {
		for index, element := range json.Buttons {
			fmt.Printf("\n%v) %v", index + 1, element.Text)
		}
	}
	fmt.Print("\n\n")
	var input uint
	fmt.Print("> ")
	fmt.Scanln(&input)
	if int(input) == 0 {
		var newurl []string = url[:len(url)-1]
		return initialize(newurl)
	} else if int(input) > len(json.Buttons) {
		err = fmt.Errorf("invalid value: %v", input)
		return err
	} else {
		var newurl []string = append(url, json.Buttons[input - 1].NextPath)
		return initialize(newurl)
	}
	return nil
}

func getJson(url string) (Json, error) {
	response, err := http.Get(url)
	if err != nil {
		return Json{}, fmt.Errorf("cannot fetch URL %q: %v", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
        return Json{}, fmt.Errorf("unexpected http GET status: %s", response.Status)
    }
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Json{}, fmt.Errorf("failed to read body: %v", err)
	}
	var stringBody string = string(body)
	res := Json{}
	err = json.Unmarshal([]byte(stringBody), &res)
	if err != nil {
        return Json{}, fmt.Errorf("cannot decode JSON: %v", err)
    }
    return res, nil
}
