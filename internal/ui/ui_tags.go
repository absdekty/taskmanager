package ui

import (
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"github.com/absdekty/taskmanager/pkg/logger"
)

func (ui *UI) tagsAddButton(tagName string) {
	var btn *widget.Button

	btn = widget.NewButton(tagName, func() {
		ui.tags.scroll.Remove(btn)
		ui.tags.scroll.Refresh()
		
		err := ui.service.RemoveTag(ui.ctx, ui.currentTask.ID, tagName)
		if err != nil {
			logger.Error.Printf("tagsAddButton: %v", err)
		}
		
		ui.currentTask.RemoveTag(tagName)
		ui.tags.updateContent()
	})

	ui.tags.scroll.Add(btn)
	ui.tags.scroll.Refresh()
}

func (ui *UI) tagsAddDefaultButton() *widget.Button {
	btn := widget.NewButton("тут будут теги..", nil)
	btn.Disable()
	return btn
}

func (ui *UI) InitTags() {
	ui.tags.scroll = container.NewHBox()
	
	/* Добавить тег */
	entryAdd := widget.NewEntry()
	entryAdd.SetPlaceHolder("тег..")
	
	btnAdd := widget.NewButton("Добавить", func() {
		if ui.currentTask == nil {
			entryAdd.SetText("")
			ui.tags.updateContent()
			return
		}
		
		if entryAdd.Text == "" {
			return
		}
		
		text := entryAdd.Text
		entryAdd.SetText("")
		
		err := ui.service.AddTag(ui.ctx, ui.currentTask.ID, text)
		if err != nil {
			logger.Error.Printf("tags btnAdd: %v", err)
			return
		}
		
		ui.currentTask.AddTag(text)
		ui.tagsAddButton(text)
		ui.tags.updateContent()
	})
	
	/* Функция обновления контента */
	ui.tags.updateContent = func() {
		entryAdd.SetText("")
		ui.tags.scroll.RemoveAll()
		
		if ui.currentTask == nil {
			btnAdd.Disable()
			ui.tags.scroll.Add(ui.tagsAddDefaultButton())
		} else {
			btnAdd.Enable()
			
			if len(ui.currentTask.Tags) == 0 {
				ui.tags.scroll.Add(ui.tagsAddDefaultButton())
			} else {
				for _, val := range ui.currentTask.Tags {
					ui.tagsAddButton(val)
				}
			}
		}
		
		ui.tags.scroll.Refresh()
	}
	
	/* Контент */
	ui.sections.tags = container.NewMax(
		ui.tags.bg,
		container.NewVBox(
			widget.NewLabel("Теги:"),
			container.NewBorder(
				nil, nil, btnAdd, nil, entryAdd),
			container.NewHScroll(ui.tags.scroll)))
}