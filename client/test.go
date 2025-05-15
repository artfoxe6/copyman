package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main222() {

	data := []string{"Item 1", "Item 2", "Item 3"}
	a := app.New()

	w := a.NewWindow("Example")

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(data[id])
		},
	)
	w.SetContent(list)

	// 设置原始按键监听
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		fmt.Printf("Key pressed: %s\n", ke.Name)
	})

	w.Show()
	a.Run()

	//w := a.NewWindow("Raw Key Press Example")
	//
	//label := widget.NewLabel("Press any key...")
	//
	//// 设置原始按键监听
	//w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
	//	fmt.Printf("Key pressed: %s\n", ke.Name)
	//	label.SetText(fmt.Sprintf("You pressed: %s", ke.Name))
	//})
	//
	//w.SetContent(container.NewVBox(label))
	//w.ShowAndRun()
}
