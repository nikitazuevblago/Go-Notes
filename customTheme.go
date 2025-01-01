package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
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
