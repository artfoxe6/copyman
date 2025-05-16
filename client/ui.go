package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("搜索过滤 + 内容展示")

	// 原始数据
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

	filteredData := make([]string, len(data))
	copy(filteredData, data)

	// 用于右侧内容展示的标签（多行）
	detailLabel := widget.NewLabel("")
	//detailLabel.SetPlaceHolder("请选择左侧列表中的项以查看完整内容")
	//detailLabel.Disable() // 设置为只读

	// 左侧列表
	list := widget.NewList(
		func() int {
			return len(filteredData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // 每项显示前30字符
		},
		func(i int, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(Format(filteredData[i], 30))
		},
	)
	list.OnSelected = func(id int) {
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
		list.Refresh()
		detailLabel.SetText("") // 清空右侧内容
	}

	// 左侧：搜索框 + 滚动列表
	left := container.NewBorder(searchEntry, nil, nil, nil, container.NewVScroll(list))

	// 整体布局：左右并排
	content := container.NewHSplit(left, detailLabel)
	content.Offset = 0.4 // 左右比例

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
