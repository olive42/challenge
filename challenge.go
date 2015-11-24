package main

import (
	"log"
	"net/http"

	"github.com/olive42/challenge/challenge"
)

func main() {
	http.HandleFunc("/task-state/", challenge.GetTasksState)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
