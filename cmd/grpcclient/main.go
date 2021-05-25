package main

import (
	"os"
	"strings"
	"task/services/endpoints/grpc/client"
)

func main() {
	client.RunClient(nil, getNumbersFromEnv()...)
}

func getNumbersFromEnv() []string {
	numbString := os.Getenv("DATES")
	split := strings.Split(numbString, " ")
	return split
}
