package main

import (
	"flag"
)

func main() {
	var action string
	var data string
	flag.StringVar(&data, "data", "bar", "data to publish")
	flag.StringVar(&action, "action", "", "sub or pub")
	flag.Parse()
	if action == "sub" {
		jsSub()
	}
	if action == "pub" {
		jsPub(data)
	}
}
