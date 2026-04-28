package ui

import (
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"github.com/absdekty/taskmanager/internal/model"
	"github.com/absdekty/taskmanager/pkg/logger"
)

func (ui *UI) InitSearch() {
	ui.search.tasks = []*model.Task{}
	
	/* Добавить задачу */
	AddBtn := widget.NewButton("Добавить", func() {
		task, err := ui.service.CreateTask(ui.ctx, "Новая задача")
		if err != nil {
			logger.Error.Printf("search addBtn: %v", err)
			return
		}
		
		ui.search.tasks = append(ui.search.tasks, task)
		ui.search.list.Refresh()
	})

	/* Поиск */
	Entry := widget.NewEntry()
	Entry.SetPlaceHolder("поиск..")
	Entry.OnChanged = func(text string) {
		tasks, err := ui.service.ListTasks(ui.ctx)
		if err != nil {
			logger.Error.Printf("search updateContent: %v", err)
			tasks = []*model.Task{}
		}
	
		if text == "" || tasks == nil || len(tasks) == 0 {
			ui.search.tasks = []*model.Task{}
			ui.search.list.Refresh()
			
			return
		}
	
		text = strings.ToLower(text)
		
		ui.search.tasks = []*model.Task{}
		
		for _, task := range tasks {
			if strings.Contains(strings.ToLower(task.Title), text) {
				ui.search.tasks = append(ui.search.tasks, task)
			}
		}
		
		ui.search.list.Refresh()
	}
	
	ui.search.list = widget.NewList(
		func() int {
			return len(ui.search.tasks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(ui.search.tasks[id].Title)
		})
	
	ui.search.list.OnSelected = func(id widget.ListItemID) {
		ui.currentTask = ui.search.tasks[id]
		
		logger.Info.Printf("выбрана задача %+v", ui.currentTask)
		
		/* Вызов обновления тегов, задачи*/
	}
	
	/* Функция обновления контента */
	ui.search.updateContent = func() {
		tasks, err := ui.service.ListTasks(ui.ctx)
		if err != nil {
			logger.Error.Printf("search updateContent: %v", err)
			tasks = []*model.Task{}
		}
	
		ui.search.tasks = tasks
		if tasks == nil {
			ui.search.tasks = []*model.Task{}
		}
		Entry.SetText("")
		ui.search.list.Refresh()
	}
	
	/* Контент */
	ui.sections.search = container.NewMax(
		ui.search.bg,
		container.NewBorder(
			container.NewVBox(
				widget.NewLabel("Задачи:"),
				AddBtn,
				Entry),
			nil, nil, nil,
			ui.search.list))
}