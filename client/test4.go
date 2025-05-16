package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// truncate 截断字符串

func main4() {
	myApp := app.New()
	myWindow := myApp.NewWindow("搜索 + 内容展示 + 键盘导航")

	data := []string{
		"苹果是一种非常常见且广泛种植的水果，富含维生素和纤维。",
		"香蕉含有丰富的钾元素，适合运动后食用。",
		"橘子味道酸甜，富含维生素C。",
		"西瓜是夏季消暑佳品，含水量高。",
		"草莓颜色鲜艳，口感酸甜，营养丰富。",
		"蓝莓具有抗氧化的功能，对眼睛有益。",
		"菠萝含有蛋白酶，能促进消化。",
		"芒果富含胡萝卜素，是热带水果之王。",
		"柠檬味酸，用于调味和美容。",
		"樱桃含铁量高，适合补血养颜。",
	}

	// 当前过滤后的数据与索引
	filteredData := append([]string{}, data...)
	selectedIndex := -1

	// 右侧展示内容
	detailLabel := widget.NewMultiLineEntry()
	detailLabel.Disable()

	// 列表控件
	list := widget.NewList(
		func() int {
			return len(filteredData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i int, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(truncate(filteredData[i], 30))
		},
	)
	list.OnSelected = func(id int) {
		selectedIndex = id
		detailLabel.SetText(filteredData[id])
	}

	// 搜索框
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("输入关键词过滤...")
	searchEntry.OnChanged = func(input string) {
		filteredData = nil
		lower := strings.ToLower(input)
		for _, item := range data {
			if strings.Contains(strings.ToLower(item), lower) {
				filteredData = append(filteredData, item)
			}
		}
		list.UnselectAll()
		selectedIndex = -1
		detailLabel.SetText("")
		list.Refresh()
	}

	// 创建键盘监听容器
	keyCapture := NewKeyCapture(func(direction int) {
		if len(filteredData) == 0 {
			return
		}
		if direction < 0 && selectedIndex > 0 {
			selectedIndex--
		} else if direction > 0 && selectedIndex < len(filteredData)-1 {
			selectedIndex++
		}
		list.Select(selectedIndex)
		detailLabel.SetText(filteredData[selectedIndex])
	})

	// 左侧区域
	left := container.NewBorder(searchEntry, nil, nil, nil, container.NewVScroll(list))
	split := container.NewHSplit(left, detailLabel)
	split.Offset = 0.4

	keyCapture.SetContent(split)

	// 设置焦点到 keyCapture 以便接收键盘事件
	myWindow.SetContent(keyCapture)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.Canvas().Focus(keyCapture)
	myWindow.ShowAndRun()
}

// KeyCapture 是一个可接收键盘上下键的容器
type KeyCapture struct {
	widget.BaseWidget
	content fyne.CanvasObject
	onArrow func(direction int) // -1 = up, 1 = down
}

func NewKeyCapture(onArrow func(direction int)) *KeyCapture {
	k := &KeyCapture{onArrow: onArrow}
	k.ExtendBaseWidget(k)
	return k
}

func (k *KeyCapture) SetContent(obj fyne.CanvasObject) {
	k.content = obj
	k.Refresh()
}

func (k *KeyCapture) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(k.content)
}

// Focusable 支持
func (k *KeyCapture) FocusGained()     {}
func (k *KeyCapture) FocusLost()       {}
func (k *KeyCapture) Focused() bool    { return true }
func (k *KeyCapture) TypedRune(r rune) {}
func (k *KeyCapture) TypedKey(e *fyne.KeyEvent) {
	switch e.Name {
	case fyne.KeyUp:
		k.onArrow(-1)
	case fyne.KeyDown:
		k.onArrow(1)
	}
}

// 接口断言
var _ fyne.Focusable = (*KeyCapture)(nil)
