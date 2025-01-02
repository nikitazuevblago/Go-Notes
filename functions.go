package main

import (
	"encoding/gob"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Sub functions
func addNote(entry *widget.Entry, notesNameContainer *fyne.Container, noteName string, noteText string) {
	var previousButtonText string
	if currentHighlightedNoteName != nil {
		previousButtonText = currentHighlightedNoteName.Objects[1].(*widget.Button).Text
		if previousButtonText == "" {
			previousButtonText = "Untitled"
		}
	}
	// Save the previous note text
	notesDB[previousButtonText] = entry.Text
	// Clear the entry
	//entry.SetText("")
	entry.SetText(noteText)
	if currentHighlightedNoteName != nil {
		currentHighlightedNoteName.Objects[1].(*widget.Button).SetText(previousButtonText)
		currentHighlightedNoteName.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
	}
	// Add new note
	// noteName := "Untitled"
	noteNameButton := &widget.Button{
		Text:      noteName,
		Alignment: widget.ButtonAlignLeading,
	}
	noteNameButton.OnTapped = func() {
		var previousButtonText string
		if currentHighlightedNoteName != nil {
			previousButtonText = currentHighlightedNoteName.Objects[1].(*widget.Button).Text
			if previousButtonText == "" {
				previousButtonText = "Untitled"
			}
		}
		// Save the previous note text
		notesDB[previousButtonText] = entry.Text

		// Open the note
		entry.SetText(notesDB[noteNameButton.Text])
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
		// Clear the entry if there are no notes
		if len(notesNameContainer.Objects) == 1 {
			entry.SetText(defaultNote)
		} else if len(notesNameContainer.Objects) > 1 {
			// Change current note
			if len(notesNameContainer.Objects) > indexToRemove+1 { // Check if there is a next note
				currentNoteName = notesNameContainer.Objects[indexToRemove+1].(*fyne.Container).Objects[1].(*widget.Button).Text
				entry.SetText(notesDB[currentNoteName])
				// Highlight the note name background with white color
				bg := canvas.NewRectangle(color.White)
				highLightedNoteName := container.NewStack(bg, notesNameContainer.Objects[indexToRemove+1].(*fyne.Container).Objects[1].(*widget.Button))
				notesNameContainer.Objects[indexToRemove+1] = highLightedNoteName
				currentHighlightedNoteName = highLightedNoteName

			} else {
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
	}
}

// Function to save the DB
func saveDB(filename string, data map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

// Function to load the DB
func loadDB(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]string
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
