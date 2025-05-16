package main

import (
	"fmt"
	"fyne.io/fyne/v2/driver/desktop"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 原始列表数据
var data = []string{
	"苹果", "香蕉", "橘子", "西瓜", "草莓",
	"蓝莓", "菠萝", "芒果", "柠檬", "樱桃",
}

// 创建一个切片用于存储过滤后的数据
var filteredData = make([]string, len(data))

var clipboardContent string

var list *widget.List
var myWindow fyne.Window

var copyChangeChan = make(chan string, 10)

func main2() {
	myApp := app.New()
	myWindow = myApp.NewWindow("搜索过滤示例")
	copy(filteredData, data)
	// 创建列表（使用 widget.List）
	list := widget.NewList(
		func() int {
			return len(filteredData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i int, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(filteredData[i])
		},
	)

	// 输入框
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("请输入关键词进行搜索")
	searchEntry.OnChanged = func(input string) {
		filteredData = nil
		lower := strings.ToLower(input)
		for _, item := range data {
			if strings.Contains(strings.ToLower(item), lower) {
				filteredData = append(filteredData, item)
			}
		}
		list.Refresh()
	}
	go func() {
		for {
			select {
			case <-copyChangeChan:
				filteredData = nil
				for _, item := range data {
					filteredData = append(filteredData, Format(item, 50))
				}
				fyne.Do(func() {
					list.Refresh()
				})
			}
		}
	}()

	// 垂直布局
	content := container.NewBorder(searchEntry, nil, nil, nil, container.NewVScroll(list))
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 400))

	clipboardContent = fyne.CurrentApp().Clipboard().Content()
	go ListenClipboardChange()
	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show", func() {
				myWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}

	myWindow.SetCloseIntercept(func() {
		myWindow.Hide()
	})
	myWindow.ShowAndRun()
}

func ListenClipboardChange() {
	for true {
		content := fyne.CurrentApp().Clipboard().Content()
		if content != "" && content != clipboardContent {
			fmt.Println(content)
			clipboardContent = content
			copyChangeChan <- content
			data = append([]string{content}, data...)
		}
		time.Sleep(time.Millisecond * 100)
	}

}
