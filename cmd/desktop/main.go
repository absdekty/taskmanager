package main

import (
	"context"
	"github.com/absdekty/taskmanager/internal/repository/sqlite"
	"github.com/absdekty/taskmanager/internal/service"
	"github.com/absdekty/taskmanager/internal/ui"
	"github.com/absdekty/taskmanager/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.Init()

	/* Корневой контекст */
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer cancel()

	/* Инициализация репозитория */
	repo, err := repository.NewRepository("./tasks.db")
	if err != nil {
		logger.Error.Fatal(err)
	}
	defer repo.DB.Close()

	/* Инициализация сервиса */
	taskService := service.NewService(repo)

	/* Инициализация UI + запуск */
	uiApp := ui.NewUI(taskService)

	if err := uiApp.Run(ctx); err != nil {
		logger.Error.Fatalf("Ошибка запуска программы: %v", err)
	}

	/* Завершение программы, ожидание конца запросов БД */
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer shutdownCancel()

	<-shutdownCtx.Done()

	if shutdownCtx.Err() == context.DeadlineExceeded {
		logger.Info.Println("Таймаут завершения, принудительное закрытие")
	} else {
		logger.Info.Println("Все операции завершены")
	}
}
