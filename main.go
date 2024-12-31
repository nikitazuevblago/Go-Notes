package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	fontColor     = color.NRGBA{R: 70, G: 58, B: 17, A: 255}
	leftSideColor = color.NRGBA{R: 242, G: 235, B: 155, A: 255}
	notesColor    = color.NRGBA{R: 216, G: 210, B: 140, A: 255}
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

	// dynamic container for notes names
	notesNameContainer := container.NewVBox()
	noteNameButton := &widget.Button{
		Text:      "Note 1",
		Alignment: widget.ButtonAlignLeading,
		OnTapped: func() {
			// TODO: open the note
		},
	}
	notesNameContainer.Add(noteNameButton)

	// wrap notesNameContainer in a scroll container
	scrollNotesNameContainer := container.NewVScroll(notesNameContainer)
	scrollNotesNameContainer.SetMinSize(fyne.NewSize(100, 300))

	// Left side of split container
	leftSide := container.NewVBox(
		container.NewHBox(widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			// Add new note
			noteNameButton := &widget.Button{
				Text:      fmt.Sprintf("Note %d", len(notesNameContainer.Objects)+1),
				Alignment: widget.ButtonAlignLeading,
				OnTapped: func() {
					// TODO: open the note
				},
			}
			notesNameContainer.Add(noteNameButton)
			fmt.Println("Add Note")
		}), widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			// Remove the last note
			if len(notesNameContainer.Objects) > 0 {
				notesNameContainer.Remove(notesNameContainer.Objects[len(notesNameContainer.Objects)-1])
				fmt.Println("Remove Note")
			}
		})),
		scrollNotesNameContainer,
	)
	leftSideRect := canvas.NewRectangle(leftSideColor)
	leftSideWithBackground := container.NewStack(leftSideRect, leftSide)

	// Right side of split container
	entry := widget.NewMultiLineEntry()
	rightSideWithBackground := container.NewStack(canvas.NewRectangle(notesColor), entry)

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
