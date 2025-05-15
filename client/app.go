package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"image/color"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Grid Layout")

	input := widget.NewEntry()
	input.Resize(fyne.NewSize(300, 50))
	text2 := canvas.NewText("2", color.White)
	text3 := canvas.NewText("3", color.White)
	grid := container.NewVBox(input, text2, text3)
	w.SetContent(grid)
	w.Show()
	a.Run()
}
