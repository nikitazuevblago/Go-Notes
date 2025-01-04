package main

import (
	"encoding/gob"
	"fmt"
	"image/color"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// Main functions
func addNote(entry *widget.Entry, notesNameContainer *fyne.Container, newNoteName string, newNoteText string) {
	// Check if current note obj is in DB before adding a new note
	var isInDB bool
	var currentNoteName string
	if currentNoteNameContainer != nil {
		for i, note := range notesDB {
			if currentNoteNameContainer == note.noteNameContainer {
				isInDB = true
				// Update current note characteristics
				currentNoteName = note.noteNameContainer.Objects[1].(*widget.Button).Text
				notesDB[i] = Note{noteNameContainer: note.noteNameContainer, Name: currentNoteName, Text: entry.Text}
			}
		}
		if !isInDB {
			currentButtonText := currentNoteNameContainer.Objects[1].(*widget.Button).Text
			firstLine := strings.Split(currentButtonText, "\n")[0]
			if firstLine == "" {
				currentNoteName = "Untitled"
			} else {
				if len(firstLine) > 15 {
					currentNoteName = firstLine[:15] + "..."
				} else {
					currentNoteName = firstLine
				}
			}
			notesDB[len(notesDB)] = Note{noteNameContainer: currentNoteNameContainer, Name: currentNoteName, Text: entry.Text}
		}
	}
	// GUI manipulations with the current note
	if currentNoteNameContainer != nil {
		// Set current note name static (before it was connected to the entry)
		currentNoteNameContainer.Objects[1].(*widget.Button).SetText(currentNoteName)
		// Remove highlight from the current note
		currentNoteNameContainer.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
		// WARNING: removed currentNoteNameContainer.Refresh()
	}
	// Turn off the OnChanged logic so it won't be triggered when the entry is cleared
	entry.OnChanged = nil

	// Clear the entry
	entry.SetText(newNoteText)

	// Creating architecture for new note
	// Make a button for newNoteName
	noteNameButton := &widget.Button{
		Text:      newNoteName,
		Alignment: widget.ButtonAlignLeading,
	}
	// Implement OnTapped logic
	noteNameButton.OnTapped = func() {
		// Process if the noteNameButton is not the same as the current noteNameButton
		if currentNoteNameContainer.Objects[1].(*widget.Button) != noteNameButton {
			var isInDB bool
			var currentNoteName string
			if currentNoteNameContainer != nil {
				for i, note := range notesDB {
					if currentNoteNameContainer == note.noteNameContainer {
						isInDB = true
						// Update current note characteristics
						currentNoteName = note.noteNameContainer.Objects[1].(*widget.Button).Text
						notesDB[i] = Note{noteNameContainer: note.noteNameContainer, Name: currentNoteName, Text: entry.Text}
					}
				}
				if !isInDB {
					currentButtonText := currentNoteNameContainer.Objects[1].(*widget.Button).Text
					firstLine := strings.Split(currentButtonText, "\n")[0]
					if firstLine == "" {
						currentNoteName = "Untitled"
					} else {
						if len(firstLine) > 15 {
							currentNoteName = firstLine[:15] + "..."
						} else {
							currentNoteName = firstLine
						}
					}
					notesDB[len(notesDB)] = Note{noteNameContainer: currentNoteNameContainer, Name: currentNoteName, Text: entry.Text}
				}
			}
			// GUI manipulations with the current note
			if currentNoteNameContainer != nil {
				// Set current note name static (before it was connected to the entry)
				currentNoteNameContainer.Objects[1].(*widget.Button).SetText(currentNoteName)
				// Remove highlight from the current note
				currentNoteNameContainer.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
				// WARNING: removed currentNoteNameContainer.Refresh()
			}
			// Turn off the OnChanged logic so it won't be triggered when the entry is cleared
			entry.OnChanged = nil

			// Open the note
			for _, note := range notesDB {
				if note.noteNameContainer.Objects[1].(*widget.Button) == noteNameButton {
					// Set the text of the entry to the text of the tapped note
					entry.SetText(note.Text)
					// Highligh the tapped note
					note.noteNameContainer.Objects[0].(*canvas.Rectangle).FillColor = color.White
					// Turn on the OnChanged logic
					entry.OnChanged = func(s string) {
						var firstLine string
						strippedString := strings.TrimSpace(s)
						firstLine = strings.Split(strippedString, "\n")[0]
						if len(firstLine) > 15 {
							firstLine = firstLine[:15] + "..."
						}
						// Change the current note name
						currentNoteNameContainer.Objects[1].(*widget.Button).SetText(firstLine)
					}
					// Set the current note name container to the tapped note
					currentNoteNameContainer = note.noteNameContainer
				}
			}
		}
	}
	// Implement OnChanged logic to dynamically change the name of the note based on the text in the entry
	entry.OnChanged = func(s string) {
		var firstLine string
		strippedString := strings.TrimSpace(s)
		firstLine = strings.Split(strippedString, "\n")[0]
		if len(firstLine) > 15 {
			firstLine = firstLine[:15] + "..."
		}
		// Change the current note name
		currentNoteNameContainer.Objects[1].(*widget.Button).SetText(firstLine)
	}
	// Highlight newly created note name
	bg := canvas.NewRectangle(color.White)
	highLightedNoteName := container.NewStack(bg, noteNameButton)
	// Add the new NoteNameContainer to the container
	notesNameContainer.Add(highLightedNoteName)
	// Change current NoteNameContainer
	currentNoteNameContainer = highLightedNoteName
}

func removeNote(entry *widget.Entry, notesNameContainer *fyne.Container) {
	if len(notesNameContainer.Objects) > 0 {
		// Check if the current note is in DB
		var isInDB bool
		for _, note := range notesDB {
			if note.noteNameContainer == currentNoteNameContainer {
				isInDB = true
				break
			}
		}
		// If the current note is not in DB, then add it
		if !isInDB {
			currentNoteName := currentNoteNameContainer.Objects[1].(*widget.Button).Text
			notesDB[len(notesDB)] = Note{noteNameContainer: currentNoteNameContainer, Name: currentNoteName, Text: entry.Text}
		}

		// Switching, highlighting another note and removing the current note from the container and DB
		// Remove highlight from the current note
		currentNoteNameContainer.Objects[0].(*canvas.Rectangle).FillColor = color.Transparent
		if len(notesNameContainer.Objects) == 1 {
			// Turn off the OnChanged logic so it won't be triggered when the entry is set to defaultNote
			entry.OnChanged = nil
			// Set the text of the entry to the default note
			entry.SetText(defaultNote)
			notesNameContainer.Remove(currentNoteNameContainer)
			// Clear up the current note name container
			currentNoteNameContainer = nil
			delete(notesDB, 0)
		} else {
			for i, note_obj := range notesNameContainer.Objects {
				if note_obj == currentNoteNameContainer {
					// Remove the current note from notesDB
					delete(notesDB, i)
					// Switch to the next or previous note
					if len(notesNameContainer.Objects) > i+1 { // Check if the current note is not the last one
						currentNoteNameContainer = notesNameContainer.Objects[i+1].(*fyne.Container)
					} else {
						currentNoteNameContainer = notesNameContainer.Objects[i-1].(*fyne.Container)
					}
					// Highlight the new current note
					currentNoteNameContainer.Objects[0].(*canvas.Rectangle).FillColor = color.White
					// Remove the current note from the container
					notesNameContainer.Remove(note_obj)
					// Sort keys of notesDB map for easier reordering
					keys := make([]int, 0, len(notesDB))
					for key := range notesDB {
						keys = append(keys, key)
					}
					sort.Ints(keys)
					// Reorder the notesDB so indices don't have gaps
					intendedIndex := 0
					for _, key := range keys {
						if key != intendedIndex {
							notesDB[intendedIndex] = notesDB[key]
							delete(notesDB, key)
						}
						intendedIndex++
					}
					entry.SetText(currentNoteNameContainer.Objects[1].(*widget.Button).Text)
					break
				}
			}
		}
	}
}

// Function to save the DB
func saveNotesToDB(filename string, entry *widget.Entry, notesNameContainer *fyne.Container) error {
	// Check if current note is in DB, add if not
	if len(notesDB) < len(notesNameContainer.Objects) {
		currentNoteName := currentNoteNameContainer.Objects[1].(*widget.Button).Text
		notesDB[len(notesDB)] = Note{noteNameContainer: currentNoteNameContainer, Name: currentNoteName, Text: entry.Text}
	}
	// Remove Untitled notes from DB
	for i, note := range notesDB {
		if note.Name == "Untitled" {
			delete(notesDB, i)
		}
	}
	// Sort keys of notesDB map for easier reordering
	keys := make([]int, 0, len(notesDB))
	for key := range notesDB {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	// Reorder the notesDB so indices don't have gaps
	intendedIndex := 0
	for _, key := range keys {
		if key != intendedIndex {
			notesDB[intendedIndex] = notesDB[key]
			delete(notesDB, key)
		}
		intendedIndex++
	}

	// Open a writer using storage.Writer
	writer, err := storage.Writer(storage.NewFileURI(filename))
	if err != nil {
		return err
	}
	defer writer.Close()

	// Encode the notesDB map into the file
	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(notesDB)
	if err != nil {
		return err
	}
	return nil
}

// Function to load the DB
func loadNotesFromDB(filename string) (map[int]Note, error) {
	// Create a URI for the file
	fileURI := storage.NewFileURI(filename)

	// Check if the file exists
	exists, err := storage.Exists(fileURI)
	if err != nil {
		return nil, fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !exists {
		fmt.Println("Database file does not exist. Initializing new notes DB.")
		return make(map[int]Note), nil
	}

	// Open a reader using storage.Reader
	reader, err := storage.Reader(fileURI)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %w", err)
	}
	defer reader.Close()

	var data map[int]Note
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return data, nil
}
