package main

import (
	"image/color"
	"log"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Set default note
const (
	defaultNote = "Welcome!\nTap '+' in the toolbar to add a note."
)

// Type of database element
type Note struct {
	noteNameContainer *fyne.Container
	Name              string
	Text              string
}

var (
	// Set colors
	fontColor  = color.NRGBA{R: 70, G: 58, B: 17, A: 255}
	mainColor  = color.NRGBA{R: 242, G: 235, B: 155, A: 255}
	entryColor = color.NRGBA{R: 216, G: 210, B: 140, A: 255} // Color of the note entry (text area)

	// Set local DB for notes
	notesDB = make(map[int]Note)

	// Set current note container
	currentNoteNameContainer *fyne.Container
)

func main() {
	// Create new app
	app := app.NewWithID("com.example.gonotes")
	log.Println("App started")

	// Set icon
	appIcon := resourceIconPng // This is the variable name generated by `fyne bundle` in bundled.go
	app.SetIcon(appIcon)
	log.Println("CUSTOM_LOG: Icon set")

	// Set theme
	app.Settings().SetTheme(&customTheme{})

	// Set window
	window := app.NewWindow("Go Notes")
	window.Resize(fyne.NewSize(500, 320))
	window.SetPadded(false) // Removes padding
	//window.CenterOnScreen() // Ensures it is properly centered
	log.Println("CUSTOM_LOG: window set")

	// General background container
	generalBackgroundRect := canvas.NewRectangle(mainColor) // WARNING: MAYBE NO NEED FOR THIS
	generalBackgroundRect.Resize(window.Canvas().Size())

	// Splitting general background container into two parts, note names on the left and note text on the right
	// Right side of split container
	entry := widget.NewMultiLineEntry() // Text area for note text
	entry.SetText(defaultNote)

	// Left side of split container
	rightSideWithBackground := container.NewStack(canvas.NewRectangle(entryColor), entry) // Background and text area
	zeroPaddingLayout := layout.NewCustomPaddedVBoxLayout(0)                              // Make zero padding between note names
	notesNameContainer := container.New(zeroPaddingLayout)                                // Container for note names
	scrollNotesNameContainer := container.NewVScroll(notesNameContainer)                  // wrap notesNameContainer in a scroll container
	scrollNotesNameContainer.SetMinSize(fyne.NewSize(100, 300))
	leftSide := container.NewVBox(
		container.NewHBox(widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			addNote(entry, notesNameContainer, "Untitled", "")
		}), widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			removeNote(entry, notesNameContainer)
		})),
		scrollNotesNameContainer,
	)
	leftSideRect := canvas.NewRectangle(mainColor)
	leftSideWithBackground := container.NewStack(leftSideRect, leftSide)

	// Combine left and right sides in a horizontal split container
	split := container.NewHSplit(
		leftSideWithBackground,
		rightSideWithBackground,
	)
	// Adjust the slider of the split container to be more to the left
	split.SetOffset(0.3)

	// Making background for the whole GUI
	content := container.NewStack(generalBackgroundRect, split) // WARNING: MAYBE NO NEED FOR THIS

	// Load notes from DB, add to notesNameContainer in order
	var err error
	var loadedNotesDB map[int]Note
	loadedNotesDB, err = loadNotesFromDB("notes.db")
	if err != nil {
		loadedNotesDB = make(map[int]Note)
		log.Printf("CUSTOM_LOG: failed to load notes from DB: %v", err)
	} else {
	    log.Println("CUSTOM_LOG: DB loaded")
	}
	// Sort notesDB keys
	keys := make([]int, 0, len(loadedNotesDB))
	for key := range loadedNotesDB {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	// Add notes to notesNameContainer in order
	for _, key := range keys {
		addNote(entry, notesNameContainer, loadedNotesDB[key].Name, loadedNotesDB[key].Text)
	}

	// Set window content
	window.SetContent(content)

	// Handle lifecycle events for cross-platform Android
	lifecycle := app.Lifecycle()
    lifecycle.SetOnStopped(func() {
        log.Println("App is stopping")
        // Perform cleanup tasks here
        err := saveNotesToDB("notes.db", entry, notesNameContainer)
        if err != nil {
            log.Printf("CUSTOM_LOG: failed to save notes to DB: %v", err)
        } else {
            log.Println("CUSTOM_LOG: DB saved")
        }
    })

	window.ShowAndRun()
}
