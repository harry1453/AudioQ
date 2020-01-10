package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

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

func getFirstSelected(table *walk.TableView) (int, error) {
	selected := table.SelectedIndexes()
	if selected == nil || len(selected) == 0 {
		return 0, fmt.Errorf("nothing selected")
	} else {
		return selected[0], nil
	}
}

// Blocking Prompt Dialog
func prompt(owner walk.Form, prompt string) (string, error) {
	var text *walk.TextEdit
	var dialog *walk.Dialog
	channel := make(chan string, 2)
	_, err := Dialog{
		AssignTo: &dialog,
		Name:     "Prompt",
		Title:    "Prompt",
		Layout:   VBox{Spacing: 5},
		Children: []Widget{
			TextLabel{
				Text: prompt,
			},
			TextEdit{
				AssignTo: &text,
			},
			PushButton{
				Text: "OK",
				OnClicked: func() {
					channel <- text.Text()
					dialog.Close(0)
				},
			},
		},
	}.Run(owner)
	if err != nil {
		return "", err
	}
	return <-channel, nil
}
