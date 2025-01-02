package main

// TODO:
// - Make the notes preserve an order (map doesn't preserve order)
// - Remove notes.db from commit history
// - Make code more readable
// - Add an icon for app
// - Make app load without console
// - Compile for many platforms
// - Release executables for many platforms in release v1.0.0 - Github
// - Make README.md

import (
	"fmt"
	"image/color"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Set default note
const (
	defaultNote = "Welcome!\nTap '+' in the toolbar to add a note.\nOr use the keyboard shortcut ctrl+N."
)

type Note struct {
	Name string
	Text string
}

var (
	// Set colors
	fontColor     = color.NRGBA{R: 70, G: 58, B: 17, A: 255}
	leftSideColor = color.NRGBA{R: 242, G: 235, B: 155, A: 255}
	notesColor    = color.NRGBA{R: 216, G: 210, B: 140, A: 255}

	// Set temporary DB for notes
	notesDB                    = make(map[int]Note)
	currentNoteName            string
	currentHighlightedNoteName *fyne.Container
)

// Main

func main() {
	// Set window
	app := app.New()
	app.Settings().SetTheme(&customTheme{})
	window := app.NewWindow("Fyne Notes")
	window.Resize(fyne.NewSize(500, 320))
	window.SetPadded(false) // Removes padding
	window.CenterOnScreen() // Ensures it is properly centered

	// General background container
	generalBackgroundRect := canvas.NewRectangle(notesColor)
	generalBackgroundRect.Resize(window.Canvas().Size())

	// Right side of split container
	entry := widget.NewMultiLineEntry()
	rightSideWithBackground := container.NewStack(canvas.NewRectangle(notesColor), entry)

	// dynamic container for notes names
	zeroPaddingLayout := layout.NewCustomPaddedVBoxLayout(0) // Make zero padding between note names
	notesNameContainer := container.New(zeroPaddingLayout)
	entry.SetText(defaultNote)

	// // Load the DB
	var err error
	notesDB, err = loadDB("notes.db")
	if err != nil {
		fmt.Println("Error loading DB:", err)
		notesDB = make(map[int]Note)
	} else {
		if len(notesDB) > 0 {
			fmt.Println(notesDB, "LOADING DB")
			// Extract and sort the keys
			keys := make([]int, 0, len(notesDB))
			for key := range notesDB {
				keys = append(keys, key)
			}
			sort.Ints(keys)
			// Create a slice in sorted order
			orderedNotesDB := make([]Note, 0, len(notesDB))
			for _, key := range keys {
				orderedNotesDB = append(orderedNotesDB, notesDB[key])
			}
			for _, note := range orderedNotesDB {
				addNote(entry, notesNameContainer, note.Name, note.Text)
			}
		}
	}

	// wrap notesNameContainer in a scroll container
	scrollNotesNameContainer := container.NewVScroll(notesNameContainer)
	scrollNotesNameContainer.SetMinSize(fyne.NewSize(100, 300))

	// Add shortcut to add a note
	window.Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyN,
			Modifier: fyne.KeyModifierControl,
		},
		func(shortcut fyne.Shortcut) {
			addNote(entry, notesNameContainer, "Untitled", "")
		},
	)

	// Left side of split container
	leftSide := container.NewVBox(
		container.NewHBox(widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			addNote(entry, notesNameContainer, "Untitled", "")
		}), widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			removeNote(entry, notesNameContainer)
		})),
		scrollNotesNameContainer,
	)
	leftSideRect := canvas.NewRectangle(leftSideColor)
	leftSideWithBackground := container.NewStack(leftSideRect, leftSide)

	// Split container
	split := container.NewHSplit(
		leftSideWithBackground,
		rightSideWithBackground,
	)
	// Adjust the slider to be more to the left
	split.SetOffset(0.3)

	// Stack container
	content := container.NewStack(generalBackgroundRect, split)
	window.SetContent(content)
	window.ShowAndRun()
	// Save last note if it is not in DB
	isFound := false
	currentNoteName := currentHighlightedNoteName.Objects[1].(*widget.Button).Text
	for _, note := range notesDB {
		if note.Name == currentNoteName {
			isFound = true
		}
	}
	if !isFound {
		if entry.Text != defaultNote {
			notesDB[len(notesDB)] = Note{Name: currentNoteName, Text: entry.Text}
		}
	}
	fmt.Println(notesDB, "SAVING DB")
	saveDB("notes.db", notesDB)
}
