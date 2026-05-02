package ui

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"github.com/absdekty/taskmanager/internal/model"
	"github.com/absdekty/taskmanager/pkg/logger"
)

func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return "0 м."
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d д.", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d ч.", hours))
	}
	if minutes > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d м.", minutes))
	}

	return strings.Join(parts, " ")
}

func (ui *UI) InitTask() {
	/* Вкладка "Информация" */
	// Название
	infoNameEntry := widget.NewEntry()
	infoNameEntry.OnSubmitted = func(text string) {
		task := ui.currentTask
		
		if task == nil {
			infoNameEntry.SetText("")
			infoNameEntry.SetPlaceHolder("задача..")
			infoNameEntry.Disable()
			return
		}
	
		if text == task.Title {
			infoNameEntry.SetText("")
			return
		}
	
		if text != "" {
			task.Title = text
		
			if err := ui.service.UpdateTask(ui.ctx, task); err != nil {
				logger.Error.Printf("task infoNameEntry: %v", err)
				return
			}
			
			infoNameEntry.SetText("")
			infoNameEntry.SetPlaceHolder(task.Title)
			
			ui.search.updateContent()
		}
	}
	
	// Создана
	infoCreatedEntry := widget.NewEntry()
	infoCreatedEntry.Disable()
	infoCreatedEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
	
	// Дедлайн
	infoDueTimeEntry := widget.NewEntry()
	infoDueTimeEntry.Disable()
	infoDueTimeEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
	infoDueTimeEntry.OnSubmitted = func(text string) {
		task := ui.currentTask
		
		if task == nil {
			infoDueTimeEntry.SetText("")
			infoDueTimeEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
			infoDueTimeEntry.Disable()
			return
		}
		
		if text == "" {
			infoDueTimeEntry.SetText("")
			return
		}
		
		t, err := time.ParseInLocation("15:04 02-01-2006", text, time.Now().Location())
		if err != nil {
			infoDueTimeEntry.SetText("")
			logger.Debug.Printf("task infoDueTimeEntry timeParse: %v", err)
			return
		}
		
		if t.UTC().Equal(task.DueTime) || t.Before(time.Now()) {
			infoDueTimeEntry.SetText("")
			return
		}
	
		task.DueTime = t.UTC()
	
		if err := ui.service.UpdateTask(ui.ctx, task); err != nil {
			logger.Error.Printf("task infoDueTimeEntry: %v", err)
			return
		}
			
		infoDueTimeEntry.SetText("")
		infoDueTimeEntry.SetPlaceHolder(t.Format("15:04 02-01-2006"))
		
		ui.task.updateContent()
	}
	
	// Напоминание
	infoNotifyEntry := widget.NewEntry()
	infoNotifyEntry.Disable()
	infoNotifyEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
	infoNotifyEntry.OnSubmitted = func(text string) {
		task := ui.currentTask
		
		if task == nil {
			infoNotifyEntry.SetText("")
			infoNotifyEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
			infoNotifyEntry.Disable()
			return
		}
		
		if text == "" {
			infoNotifyEntry.SetText("")
			return
		}
		
		t, err := time.ParseInLocation("15:04 02-01-2006", text, time.Now().Location())
		if err != nil {
			infoNotifyEntry.SetText("")
			logger.Debug.Printf("task infoNotifyEntry timeParse: %v", err)
			return
		}
		
		if t.UTC().Equal(task.DueTime) || t.Before(time.Now()) {
			infoNotifyEntry.SetText("")
			return
		}
	
		task.NotifyAt = t.UTC()
	
		if err := ui.service.UpdateTask(ui.ctx, task); err != nil {
			logger.Error.Printf("task infoNotifyEntry: %v", err)
			return
		}
			
		infoNotifyEntry.SetText("")
		infoNotifyEntry.SetPlaceHolder(t.Format("15:04 02-01-2006"))
		
		ui.task.updateContent()
	}
	
	// Субзадачи
	infoSubtasksEntry := widget.NewEntry()
	infoSubtasksEntry.SetText("0/0")
	infoSubtasksEntry.Disable()
	
	// Прочая информация
	infoDeadLineAt := widget.NewLabel("Будет просрочено:")
	infoNotifyAt := widget.NewLabel("Напомнить о задаче:")
	infoTotalProgress := widget.NewLabel("Общий прогресс: 0.0%")
	
	// Контент
	informationContent := container.NewVBox(
		container.NewPadded(
			container.NewBorder(
				nil, nil,
				widget.NewLabel("Название:"),
				nil,
				infoNameEntry)),
		container.NewPadded(container.NewBorder(nil, nil, widget.NewLabel("Субзадачи:"), nil, infoSubtasksEntry)),
		container.NewPadded(infoTotalProgress),
		widget.NewSeparator(),
		container.NewPadded(
			container.NewBorder(
				nil, nil,
				widget.NewLabel("Создана:"),
				nil,
				infoCreatedEntry)),
		container.NewPadded(
			container.NewBorder(
				nil, nil,
				widget.NewLabel("Просрочена будет через:"),
				nil,
				infoDueTimeEntry)),
		container.NewPadded(
			container.NewBorder(
				nil, nil,
				widget.NewLabel("Напоминание будет через:"),
				nil,
				infoNotifyEntry)),
		widget.NewSeparator(),
		widget.NewButton("Удалить", func() {
			task := ui.currentTask
			if task != nil {
				if err := ui.service.DeleteTask(ui.ctx, task.ID); err != nil {
					logger.Error.Printf("task informationContent button(delete): %v", err)
					return
				}
				
				ui.currentTask = nil
				ui.search.updateContent()
				ui.tags.updateContent()
				ui.task.updateContent()
			}
		}),
		widget.NewSeparator(),
		infoDeadLineAt,
		infoNotifyAt,
	)

	/* Вкладка "Субзадачи" */
	filteredSubtasks := []*model.Subtask{}
	
	subtasksList := widget.NewList(
		func() int { return len(filteredSubtasks) },
		func() fyne.CanvasObject {
			
			return container.NewVBox(
				container.NewBorder(nil, nil, widget.NewLabel("Название:"), nil, widget.NewEntry()),
				container.NewHBox(
					widget.NewLabel("Прогресс:"),
					widget.NewEntry(),
					widget.NewLabel("/"),
					widget.NewEntry(),
				))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			vbox := obj.(*fyne.Container)
			nameEntry := vbox.Objects[0].(*fyne.Container).Objects[0].(*widget.Entry)
			
			hbox := vbox.Objects[1].(*fyne.Container)
			progressEntry := hbox.Objects[1].(*widget.Entry)
			needEntry := hbox.Objects[3].(*widget.Entry)
			
			subtask := filteredSubtasks[id]
			
			nameEntry.SetPlaceHolder(subtask.Name)
			progressEntry.SetPlaceHolder(strconv.Itoa(subtask.Progress))
			needEntry.SetPlaceHolder(strconv.Itoa(subtask.NeedProgress))
			
			nameEntry.OnSubmitted = func(text string) {
				if text == subtask.Name {
					nameEntry.SetText("")
					return
				}
			
				if text != "" {
					subtask.Name = text
				
					if err := ui.service.UpdateSubtask(ui.ctx, subtask); err != nil {
						logger.Error.Printf("task nameEntry: %v", err)
						return
					}
					
					nameEntry.SetText("")
					nameEntry.SetPlaceHolder(subtask.Name)
				}
			}
			
			progressEntry.OnSubmitted = func(text string) {
				if text == strconv.Itoa(subtask.Progress) {
					progressEntry.SetText("")
					return
				}
			
				if text != "" {
					num, err := strconv.Atoi(text)
					if err != nil {
						progressEntry.SetText("")
						return
					}
					
					if num < 0 {
						num = 0
					} else if num > subtask.NeedProgress {
						num = subtask.NeedProgress
					}
					
					subtask.Progress = num
				
					if err := ui.service.UpdateSubtask(ui.ctx, subtask); err != nil {
						logger.Error.Printf("task progressEntry: %v", err)
						return
					}
					
					progressEntry.SetText("")
					progressEntry.SetPlaceHolder(fmt.Sprintln(num))
					
					ui.task.updateContent()
				}
			}
			
			needEntry.OnSubmitted = func(text string) {
				if text == strconv.Itoa(subtask.NeedProgress) {
					needEntry.SetText("")
					return
				}
			
				if text != "" {
					num, err := strconv.Atoi(text)
					if err != nil {
						needEntry.SetText("")
						return
					}
					
					subtask.NeedProgress = num
				
					if err := ui.service.UpdateSubtask(ui.ctx, subtask); err != nil {
						logger.Error.Printf("task needEntry: %v", err)
						return
					}
					
					needEntry.SetText("")
					needEntry.SetPlaceHolder(text)
					
					ui.task.updateContent()
				}
			}
		})
	
	subtasksBtnAdd := widget.NewButton("Добавить", func() {
		task := ui.currentTask
		if task != nil {
			_, err := ui.service.AddSubtask(ui.ctx, task.ID, "Новая субзадача", 1)
			if err != nil {
				logger.Error.Printf("task subtasksBtnAdd: %v", err)
				return
			}
			
			ui.search.updateContent()
			
			subtask, err := model.NewSubtask("Новая субзадача 1", 1)
			if err != nil {
				logger.Debug.Printf("task subtasksBtnAdd: %v", err)
				return
			}
			
			filteredSubtasks = append(filteredSubtasks, subtask)
			
			subtasksList.Refresh()
		}
	})
	
	subtasksContent := container.NewBorder(
		subtasksBtnAdd,
		nil, nil, nil,
		subtasksList,
	)
	
	/* Функция обновления контента */
	ui.task.updateContent = func() {
		task := ui.currentTask
		
		infoNameEntry.OnSubmitted("")
		
		if task != nil {
			/* Вкладка "Информация" */
			infoNameEntry.SetText("")
			infoNameEntry.SetPlaceHolder(task.Title)
			infoNameEntry.Enable()
			
			infoCreatedEntry.SetPlaceHolder(task.Created.Local().Format("15:04 02-01-2006"))
			
			infoDueTimeEntry.SetText("")
			if task.DueTime.IsZero() {
				infoDueTimeEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
			} else {
				infoDueTimeEntry.SetPlaceHolder(task.DueTime.Local().Format("15:04 02-01-2006"))
			}
			infoDueTimeEntry.Enable()
			
			infoNotifyEntry.SetText("")
			if task.NotifyAt.IsZero() {
				infoNotifyEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
			} else {
				infoNotifyEntry.SetPlaceHolder(task.NotifyAt.Local().Format("15:04 02-01-2006"))
			}
			infoNotifyEntry.Enable()
			
			infoSubtasksEntry.SetText(fmt.Sprintf("%d/%d", task.GetRemainingSubtasksCount(), len(task.Subtasks)))
			
			if !task.DueTime.IsZero() {
				infoDeadLineAt.SetText(fmt.Sprintf("Будет просрочено через: %s", FormatDuration(task.DueTime.Sub(time.Now().UTC()))))
			} else { infoDeadLineAt.SetText("") }
			
			if !task.NotifyAt.IsZero() {
				infoNotifyAt.SetText(fmt.Sprintf("Будет напоминание через: %s", FormatDuration(task.NotifyAt.Sub(time.Now().UTC()))))
			} else { infoNotifyAt.SetText("") }
			
			totalProgress, err := ui.service.GetTotalProgress(ui.ctx, task.ID)
			if err != nil {
				infoTotalProgress.SetText(fmt.Sprintln("Общий прогресс: 0.0% (ошибка загрузки)"))
			} else {
				infoTotalProgress.SetText(fmt.Sprintf("Общий прогресс: %.1f%%", totalProgress))
			}
			
			/* Вкладка "Субзадачи" */
			subtasksBtnAdd.Enable()
			
			filteredSubtasks = make([]*model.Subtask, 0, len(task.Subtasks))
			
			for _, subtask := range task.Subtasks {
				filteredSubtasks = append(filteredSubtasks, subtask)
			}
		} else {
			/* Вкладка "Информация" */
			infoCreatedEntry.SetPlaceHolder("HH:MM DD-MM-YYYY")
			
			infoSubtasksEntry.SetText("0/0")
			
			infoDeadLineAt.SetText("")
			infoNotifyAt.SetText("")
			infoTotalProgress.SetText(fmt.Sprintln("Общий прогресс: 0.0%"))
			
			/* Вкладка "Субзадачи" */
			filteredSubtasks = []*model.Subtask{}
			
			subtasksBtnAdd.Disable()
		}
		
		subtasksList.Refresh()
	}
	
	/* Контент */
	ui.sections.task = container.NewMax(
		ui.task.bg,
		container.NewBorder(
			nil, nil, nil, nil,
			container.NewAppTabs(
				container.NewTabItem("Информация", informationContent),
				container.NewTabItem("Субзадачи", subtasksContent))))
	
	/* Фоновое обновление "будет просрочено/напоминание через:" */
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			task := ui.currentTask
			if task != nil {
				now := time.Now().UTC()
				dueRemaining := task.DueTime.Sub(now)
				notifyRemaining := task.NotifyAt.Sub(now)
			
				fyne.Do(func() {
					if !task.DueTime.IsZero() {
						infoDeadLineAt.SetText(fmt.Sprintf("Будет просрочено через: %s", FormatDuration(dueRemaining)))
					} else { infoDeadLineAt.SetText("") }
					
					if !task.NotifyAt.IsZero() {
						infoNotifyAt.SetText(fmt.Sprintf("Будет напоминание через: %s", FormatDuration(notifyRemaining)))
					} else { infoNotifyAt.SetText("") }
				})
			}
		}
	}()
}