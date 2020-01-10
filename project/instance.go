package project

import "log"

func init() {
	if err := Init(); err != nil {
		log.Fatal("Error initializing project:", err)
	}
}
