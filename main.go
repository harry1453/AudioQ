package main

import (
	"github.com/harry1453/audioQ/api"
	"github.com/harry1453/audioQ/console"
	"github.com/harry1453/audioQ/gui"
)

func main() {
	go console.Initialize()
	go api.Initialize()
	go gui.Initialize()
	select {}
}
