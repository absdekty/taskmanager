package model

import "errors"

/* Domain-ошибки */
var (
	/* Общие */
	ErrEmptyName = errors.New("Название пустое")

	/* Задача */
	ErrTaskIsOverdue = errors.New("Задача в дедлайне")
	ErrPastDeadline  = errors.New("Дедлайн в прошлом")
	ErrPastNotify    = errors.New("Напоминание в прошлом")
	ErrNotExisting   = errors.New("Тег не существует")

	/* Субзадача */
	ErrInvalidProgress     = errors.New("Прогресс  меньше единицы")
	ErrMaxProgressExceeded = errors.New("Прогресс выше максимально")
	ErrMinProgressExceeded = errors.New("Прогресс отрицателен")
)

var (
	ErrTaskNotFound    = errors.New("Задача не существует")
	ErrSubtaskNotFound = errors.New("Субзадача не существует")
)
