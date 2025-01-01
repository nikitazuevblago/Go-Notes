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
	if currentHighlightedNoteName != nil {
		previousButtonText = currentHighlightedNoteName.Objects[1].(*widget.Button).Text
	}
	// Save the previous note text
	notesDB[currentNoteName] = entry.Text
	fmt.Println(currentNoteName, "saved")
	// Clear the entry
	entry.SetText("")
	if currentHighlightedNoteName != nil {
		currentHighlightedNoteName.Objects[1].(*widget.Button).SetText(previousButtonText)
		currentHighlightedNoteName.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
	}
	// Add new note
	noteName := fmt.Sprintf("Note %d", len(notesNameContainer.Objects)+1)
	noteNameButton := &widget.Button{
		Text:      noteName,
		Alignment: widget.ButtonAlignLeading,
	}
	noteNameButton.OnTapped = func() {
		if noteName != currentNoteName {
			var previousButtonText string
			if currentHighlightedNoteName != nil {
				previousButtonText = currentHighlightedNoteName.Objects[1].(*widget.Button).Text
			}
			// Save the previous note text
			notesDB[currentNoteName] = entry.Text

			fmt.Println(currentNoteName, "saved")
			// Open the note
			entry.SetText(notesDB[noteName])
			// Update the previous button's appearance
			if currentHighlightedNoteName != nil {
				currentHighlightedNoteName.Objects[1].(*widget.Button).SetText(previousButtonText)
				currentHighlightedNoteName.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
				currentHighlightedNoteName.Refresh()
			}
			// Change current note state
			currentNoteName = noteName
			// Highlight the note name background with white color
			for _, obj := range notesNameContainer.Objects {
				button := obj.(*fyne.Container).Objects[1].(*widget.Button)
				if button == noteNameButton {
					obj.(*fyne.Container).Objects[0].(*canvas.Rectangle).FillColor = color.White
					currentHighlightedNoteName = obj.(*fyne.Container)
				} else {
					obj.(*fyne.Container).Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
				}
			}
		}
	}

	// Dynamically change the note name based on first line of the note text
	entry.OnChanged = func(s string) {
		var firstLine string
		// Get the first 15 characters of the note text
		strippedString := strings.TrimSpace(s)
		firstLine = strings.Split(strippedString, "\n")[0]
		if len(firstLine) > 15 {
			firstLine = firstLine[:15] + "..."
		}
		// Change the note name
		currentHighlightedNoteName.Objects[1].(*widget.Button).SetText(firstLine)
	}
	// Highlight newly created note name
	bg := canvas.NewRectangle(color.White)
	highLightedNoteName := container.NewStack(bg, noteNameButton)
	// Add the note name to the container
	notesNameContainer.Add(highLightedNoteName)
	// Change current note state
	currentNoteName = noteNameButton.Text
	currentHighlightedNoteName = highLightedNoteName
	fmt.Println("Add Note")
}

func removeNote(entry *widget.Entry, notesNameContainer *fyne.Container) {
	currentNoteName = currentHighlightedNoteName.Objects[1].(*widget.Button).Text
	// Remove the current note
	if len(notesNameContainer.Objects) > 0 {
		// Find the index of the current note by currentHighlightedNoteName
		indexToRemove := 0
		for i, obj := range notesNameContainer.Objects {
			if obj == currentHighlightedNoteName {
				indexToRemove = i
				break
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
				currentNoteName = notesNameContainer.Objects[indexToRemove+1].(*fyne.Container).Objects[1].(*widget.Button).Text
				entry.SetText(notesDB[currentNoteName])
				// Highlight the note name background with white color
				bg := canvas.NewRectangle(color.White)
				highLightedNoteName := container.NewStack(bg, notesNameContainer.Objects[indexToRemove+1].(*fyne.Container).Objects[1].(*widget.Button))
				notesNameContainer.Objects[indexToRemove+1] = highLightedNoteName
				currentHighlightedNoteName = highLightedNoteName

			} else {
				fmt.Println("Change to previous note")
				currentNoteName = notesNameContainer.Objects[indexToRemove-1].(*fyne.Container).Objects[1].(*widget.Button).Text
				entry.SetText(notesDB[currentNoteName])
				// Highlight the note name background with white color
				bg := canvas.NewRectangle(color.White)
				highLightedNoteName := container.NewStack(bg, notesNameContainer.Objects[indexToRemove-1].(*fyne.Container).Objects[1].(*widget.Button))
				notesNameContainer.Objects[indexToRemove-1] = highLightedNoteName
				currentHighlightedNoteName = highLightedNoteName
			}
		}
		// Remove by index
		notesNameContainer.Objects = append(notesNameContainer.Objects[:indexToRemove],
			notesNameContainer.Objects[indexToRemove+1:]...)
		// Remove the note from the DB
		delete(notesDB, currentNoteName)
		fmt.Println("Remove Note")
	}
}
