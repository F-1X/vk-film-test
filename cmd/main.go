package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vk/config"
	"vk/database"
	"vk/server"
	
)

//@title vk film library
//@version 1.0
//@host localhost:8001
//@BasePath /

func main() {

	config, err := config.Read("config.yml")
	if err != nil {
		log.Fatalf("config read failed: %v", err)
	}

	db, err := database.New(config.Database)
	if err != nil {
		log.Fatalf("failed create database: %v", err)
	}

	db.InitDb()

	server := server.NewAPIServer(config.Server, *db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создание канала для получения сигнала завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Ожидание сигнала завершения работы
	go func() {
		sig := <-quit
		log.Printf("Received signal: %v", sig)
		cancel()
	}()

	// Запуск сервера
	go func() {
		// log.Printf("Server started on port %s", cfg.Server.Port)
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Ожидание сигнала отмены
	<-ctx.Done()

	// Завершение работы сервера
	log.Println("Shutting down server...")

	// Создание контекста с таймаутом для graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Завершение работы сервера с таймаутом
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("Server shutdown complete")
}
