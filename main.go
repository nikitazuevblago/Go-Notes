package main

import (
	"fmt"
	"image/color"
	"os"

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
	defaultNote = "Welcome!\nTap '+' in the toolbar to add a note.\nOr use the keyboard shortcut ctrl+N."
)

var (
	// Set colors
	fontColor     = color.NRGBA{R: 70, G: 58, B: 17, A: 255}
	leftSideColor = color.NRGBA{R: 242, G: 235, B: 155, A: 255}
	notesColor    = color.NRGBA{R: 216, G: 210, B: 140, A: 255}

	// Set temporary DB for notes
	notesDB     = make(map[string]string)
	currentNote string
)

// Custom theme
type customTheme struct{}

func (c customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameInputBackground {
		return color.Transparent
	}
	if name == theme.ColorNameForeground {
		return fontColor
	}
	if name == theme.ColorNameInputBorder {
		// Border color
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White
	}
	if name == theme.ColorNamePrimary {
		// Cursor (insert line) color
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White
	}
	if name == theme.ColorNameButton {
		return color.Transparent
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (c customTheme) Font(style fyne.TextStyle) fyne.Resource {
	fontData, err := os.ReadFile("GochiHand.ttf")
	if err != nil {
		fmt.Println(err)
	}
	return fyne.NewStaticResource("GochiHand.ttf", fontData)
}

func (c customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// Main

func main() {
	// Set window
	app := app.New()
	app.Settings().SetTheme(&customTheme{})
	window := app.NewWindow("Fyne Notes")
	window.Resize(fyne.NewSize(500, 320))

	// General background container
	generalBackgroundRect := canvas.NewRectangle(color.Black)
	generalBackgroundRect.Resize(window.Canvas().Size())

	// Right side of split container
	entry := widget.NewMultiLineEntry()
	rightSideWithBackground := container.NewStack(canvas.NewRectangle(notesColor), entry)

	// dynamic container for notes names
	zeroPaddingLayout := layout.NewCustomPaddedVBoxLayout(0) // Make zero padding between note names
	notesNameContainer := container.New(zeroPaddingLayout)
	entry.SetText(defaultNote)

	// wrap notesNameContainer in a scroll container
	scrollNotesNameContainer := container.NewVScroll(notesNameContainer)
	scrollNotesNameContainer.SetMinSize(fyne.NewSize(100, 300))

	// Left side of split container
	leftSide := container.NewVBox(
		container.NewHBox(widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			// Save the previous note text
			notesDB[currentNote] = entry.Text
			fmt.Println(entry.Text, "saved")
			// Clear the entry
			entry.SetText("")
			// Add new note
			noteName := fmt.Sprintf("Note %d", len(notesNameContainer.Objects)+1)
			noteNameButton := &widget.Button{
				Text:      noteName,
				Alignment: widget.ButtonAlignLeading,
				OnTapped: func() {
					// Save the previous note text
					notesDB[currentNote] = entry.Text
					fmt.Println(entry.Text, "saved")
					// Open the note
					entry.SetText(notesDB[noteName])
					// Change current note state
					currentNote = noteName
					// TODO: Highlight the note name background with white color
					for _, obj := range notesNameContainer.Objects {
						localContainer, ok := obj.(*fyne.Container)
						if ok && localContainer.Objects[1].(*widget.Button).Text == currentNote {
							// TODO: Highlight the note name background with white color
							// TODO: Remove the highlight from the previous note name
						}
					}
				},
			}
			// Highlight newly created note name
			bg := canvas.NewRectangle(color.White)
			highLightedNoteName := container.NewStack(bg, noteNameButton)
			// TODO: Remove the highlight from the previous note name
			// Add the note name to the container
			notesNameContainer.Add(highLightedNoteName)
			// Change current note state
			currentNote = noteName
			fmt.Println("Add Note")

		}), widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			// Remove the current note
			// TODO: Highlight the note name background with white color
			if len(notesNameContainer.Objects) > 0 {
				// Find the index of the current note
				indexToRemove := 0
				for i, obj := range notesNameContainer.Objects {
					button, ok := obj.(*widget.Button)
					if ok && button.Text == currentNote {
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
						currentNote = notesNameContainer.Objects[indexToRemove+1].(*widget.Button).Text
						entry.SetText(notesDB[currentNote])
					} else {
						fmt.Println("Change to previous note")
						currentNote = notesNameContainer.Objects[indexToRemove-1].(*widget.Button).Text
						entry.SetText(notesDB[currentNote])
					}
				}
				// Remove by index
				notesNameContainer.Objects = append(notesNameContainer.Objects[:indexToRemove],
					notesNameContainer.Objects[indexToRemove+1:]...)
				// Remove the note from the DB
				delete(notesDB, currentNote)
				fmt.Println("Remove Note")
			}
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

	// Padded container
	padded := container.NewPadded(split)

	// Stack container
	content := container.NewStack(generalBackgroundRect, padded)
	window.SetContent(content)
	window.ShowAndRun()
}
