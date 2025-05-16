package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main5() {
	a := app.New()
	w := a.NewWindow("Custom List with Key Handling")

	items := []string{"苹果", "香蕉", "葡萄", "西瓜", "橘子"}

	myList := NewMyList(items, func(index int, value string) {
		fmt.Printf("选中项：%d - %s\n", index, value)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("使用 ↑ / ↓ 选择项："),
		myList,
	))

	// 设置焦点到组件以便接收键盘事件
	w.Canvas().Focus(myList)

	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}
