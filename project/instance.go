package project

import "log"

var Instance *Project

func init() {
	Instance = new(Project)
	if err := Instance.Init(); err != nil {
		log.Fatal("Error initializing project:", err)
	}
}
