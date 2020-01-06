package main

import (
	"github.com/harry1453/audioQ/api"
	_ "github.com/harry1453/audioQ/api"
	"github.com/harry1453/audioQ/console"
)

func main() {
	go console.Initialize()
	go api.Initialize()
	select {}
}
