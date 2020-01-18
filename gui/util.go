package gui

import (
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

var PromptCancelled = errors.New("prompt cancelled or nothing entered")

func setting(name string, onUpdate func(string), newValueChannel <-chan string) Widget {
	var textEdit *walk.TextEdit
	go func() {
		for {
			textEdit.SetText(<-newValueChannel)
		}
	}()
	return Composite{
		Layout: HBox{Spacing: 5},
		Children: []Widget{
			TextLabel{
				Text: name + ": ",
			},
			TextEdit{
				AssignTo:      &textEdit,
				CompactHeight: true,
				OnTextChanged: func() {
					onUpdate(textEdit.Text())
				},
			},
		},
	}
}

func messageBox(owner walk.Form, message string) error {
	var dialog *walk.Dialog
	_, err := Dialog{
		AssignTo: &dialog,
		Name:     "Message",
		Title:    "Message",
		Layout:   VBox{Spacing: 5},
		Children: []Widget{
			TextLabel{
				Text: message,
			},
			PushButton{
				Text: "OK",
				OnClicked: func() {
					dialog.Close(0)
				},
			},
		},
	}.Run(owner)
	if err != nil {
		return err
	}
	return nil
}

func getCurrentCueIndex(owner walk.Form, table *walk.TableView) (int, bool) {
	currentCue := table.CurrentIndex()
	if currentCue <= 0 {
		handleError(owner, fmt.Errorf("no cue selected"))
		return currentCue, false
	}
	return currentCue, true
}

func handleError(owner walk.Form, err error) {
	if newErr := messageBox(owner, "Error: "+err.Error()); newErr != nil {
		log.Println("Error displaying error, original error:", err, "Error displaying error:", newErr)
	}
}

// Blocking Prompt Dialog
func prompt(owner walk.Form, prompt string) (string, error) {
	var text *walk.TextEdit
	var dialog *walk.Dialog
	channel := make(chan string, 2)
	result, err := Dialog{
		AssignTo: &dialog,
		Name:     "Prompt",
		Title:    "Prompt",
		Layout:   VBox{Spacing: 5},
		Children: []Widget{
			TextLabel{
				Text: prompt,
			},
			TextEdit{
				AssignTo:      &text,
				CompactHeight: true,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "OK",
						OnClicked: func() {
							channel <- text.Text()
							dialog.Close(0)
						},
					},
					PushButton{
						Text: "Cancel",
						OnClicked: func() {
							dialog.Close(10000)
						},
					},
				},
			},
		},
	}.Run(owner)
	if err != nil {
		return "", err
	}
	if result == 10000 {
		return "", PromptCancelled
	}
	value := <-channel
	if value == "" {
		return "", PromptCancelled
	}
	return value, nil
}
