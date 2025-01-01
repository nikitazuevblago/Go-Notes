package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Sub functions
func addNote(entry *widget.Entry, notesNameContainer *fyne.Container) {
	var previousButtonText string
	if previousButton != nil {
		previousButtonText = previousButton.Text
	}
	// Save the previous note text
	notesDB[currentNote] = entry.Text
	fmt.Println(currentNote, "saved")
	// Clear the entry
	entry.SetText("")
	if previousButton != nil {
		fmt.Println(previousButtonText, "PREVIOUS BUTTON TEXT")
		previousButton.SetText(previousButtonText)
	}
	// Add new note
	noteName := fmt.Sprintf("Note %d", len(notesNameContainer.Objects)+1)
	noteNameButton := &widget.Button{
		Text:      noteName,
		Alignment: widget.ButtonAlignLeading,
		OnTapped: func() {
			if noteName != currentNote {
				// Save the previous note text
				notesDB[currentNote] = entry.Text
				fmt.Println(currentNote, "saved")
				// Open the note
				entry.SetText(notesDB[noteName])
				// Change current note state
				currentNote = noteName
				// Iterate over the notes name container and highlight the current note
				for _, obj := range notesNameContainer.Objects {
					localContainer, ok := obj.(*fyne.Container)
					if ok && localContainer.Objects[1].(*widget.Button).Text == currentNote {
						// Highlight the note name background with white color
						localContainer.Objects[0].(*canvas.Rectangle).FillColor = color.White
					} else {
						// Remove the highlight from the previous note name
						localContainer.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
					}
				}
			}
		},
	}
	// Dynamically change the note name based on first line of the note text
	entry.OnChanged = func(s string) {
		var firstLine string
		// Get the first line of the note text
		if len(s) > 15 {
			firstLine = strings.Split(s, "\n")[0][:15] + "..."
		} else {
			firstLine = strings.Split(s, "\n")[0]
		}
		// Change the note name
		noteNameButton.SetText(firstLine)
	}
	// Highlight newly created note name
	bg := canvas.NewRectangle(color.White)
	highLightedNoteName := container.NewStack(bg, noteNameButton)
	// Iterate over the notes name container and remove the highlight from the previous note name
	for _, obj := range notesNameContainer.Objects {
		localContainer, ok := obj.(*fyne.Container)
		if ok && localContainer.Objects[1].(*widget.Button).Text == currentNote {
			// Remove the highlight from the previous note name
			localContainer.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
		}
	}
	// Add the note name to the container
	notesNameContainer.Add(highLightedNoteName)
	// Change current note state
	currentNote = noteName
	previousButton = noteNameButton
	fmt.Println("Add Note")
}

func removeNote(entry *widget.Entry, notesNameContainer *fyne.Container) {
	// Remove the current note
	if len(notesNameContainer.Objects) > 0 {
		// Find the index of the current note
		indexToRemove := 0
		for i, obj := range notesNameContainer.Objects {
			localContainer, ok := obj.(*fyne.Container)
			if ok {
				button, ok := localContainer.Objects[1].(*widget.Button)
				if ok && button.Text == currentNote {
					indexToRemove = i
					break
				}
			}
		}
		fmt.Println(indexToRemove, "removed index")
		// Clear the entry if there are no notes
		if len(notesNameContainer.Objects) == 1 {
			entry.SetText(defaultNote)
		} else if len(notesNameContainer.Objects) > 1 {
			// Change current note
			if len(notesNameContainer.Objects) > indexToRemove+1 { // Check if there is a next note
				fmt.Println("Change to next note")
				currentNote = notesNameContainer.Objects[indexToRemove+1].(*fyne.Container).Objects[1].(*widget.Button).Text
				entry.SetText(notesDB[currentNote])
				// Highlight the note name background with white color
				bg := canvas.NewRectangle(color.White)
				highLightedNoteName := container.NewStack(bg, notesNameContainer.Objects[indexToRemove+1].(*fyne.Container).Objects[1].(*widget.Button))
				notesNameContainer.Objects[indexToRemove+1] = highLightedNoteName

			} else {
				fmt.Println("Change to previous note")
				currentNote = notesNameContainer.Objects[indexToRemove-1].(*fyne.Container).Objects[1].(*widget.Button).Text
				entry.SetText(notesDB[currentNote])
				// Highlight the note name background with white color
				bg := canvas.NewRectangle(color.White)
				highLightedNoteName := container.NewStack(bg, notesNameContainer.Objects[indexToRemove-1].(*fyne.Container).Objects[1].(*widget.Button))
				notesNameContainer.Objects[indexToRemove-1] = highLightedNoteName
			}
		}
		// Remove by index
		notesNameContainer.Objects = append(notesNameContainer.Objects[:indexToRemove],
			notesNameContainer.Objects[indexToRemove+1:]...)
		// Remove the note from the DB
		delete(notesDB, currentNote)
		fmt.Println("Remove Note")
	}
}
