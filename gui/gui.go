package gui

import (
	. "github.com/lxn/walk/declarative"
	"log"
)

func Initialize() {
	mainWindow := MainWindow{}
	exit, err := mainWindow.Run()
	if err != nil {
		log.Println("GUI Error:", exit, err)
	}
}
