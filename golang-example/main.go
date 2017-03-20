package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
)

type Check struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Output string `json:"output"`
	Status int    `json:"status"`
	Ttl    int    `json:"ttl"`
}

func performBackup() (output string, err error) {
	if rand.Int()%10 == 0 {
		return "Unlucky day", errors.New("One in 10 check failed")
	} else {
		return "Everything is fine!", nil
	}
}

func main() {
	check := Check{"mysql", "mysql-backup", "", 0, 86400}
	output, err := performBackup()
	if err != nil {
		check.Status = 2
		check.Output = "Backup Failed:" + output
	} else {
		check.Output = "Backup Successful: " + output
	}
	jsonBody, _ := json.Marshal(check)
	http.Post("http://api:4567/results", "application/json", bytes.NewBuffer(jsonBody))
}
