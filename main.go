package main

import (
	"fmt"
	"image/color"

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

var (
	// Set colors
	fontColor     = color.NRGBA{R: 70, G: 58, B: 17, A: 255}
	leftSideColor = color.NRGBA{R: 242, G: 235, B: 155, A: 255}
	notesColor    = color.NRGBA{R: 216, G: 210, B: 140, A: 255}

	// Set temporary DB for notes
	notesDB                    = make(map[string]string)
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

	// Load the DB
	var err error
	notesDB, err = loadDB("notes.db")
	if err != nil {
		fmt.Println("Error loading DB:", err)
		notesDB = make(map[string]string)
	} else {
		fmt.Println(notesDB)
		for name, text := range notesDB {
			if name != "" {
				addNote(entry, notesNameContainer, name, text)
			}
		}
	}

	// wrap notesNameContainer in a scroll container
	scrollNotesNameContainer := container.NewVScroll(notesNameContainer)
	scrollNotesNameContainer.SetMinSize(fyne.NewSize(100, 300))

	// Add shortcut to add a note
	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		addNote(entry, notesNameContainer, "Untitled", "")
	})

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
	fmt.Println(notesDB)
	// Save last note if it is not in DB
	isFound := false
	currentNoteName := currentHighlightedNoteName.Objects[1].(*widget.Button).Text
	for name, text := range notesDB {
		if text == entry.Text && name == currentNoteName {
			isFound = true
		}
	}
	if !isFound {
		notesDB[currentNoteName] = entry.Text
	}
	saveDB("notes.db", notesDB)
}
