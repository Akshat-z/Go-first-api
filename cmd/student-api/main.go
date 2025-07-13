package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Akshat-z/student-api/internal/config"
	"github.com/Akshat-z/student-api/internal/http/handlers/student"
	"github.com/Akshat-z/student-api/internal/storage/sqlite"
)

func main() {
	//+ setup config
	cfg := config.MustLoad()

	//+ setp database

	storage, err := sqlite.NewSqlite(cfg)

	if err != nil {
		log.Fatal("Can't connect to db")
	}

	//+ setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.Create(storage))

	//+ setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	fmt.Printf("Server started on Address %s", cfg.Address)

	signalChan := make(chan os.Signal, 1)
	serverErrChan := make(chan error, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //$ send message to channel and (<-sserverChan) which is blocking will reseave message .

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	select {
	case err := <-serverErrChan:
		slog.Error("error while starting the server", slog.String("error", err.Error()))

		os.Exit(1)

	case <-signalChan:
		//! as sign pass in channel to shutdown
		slog.Info("shutting down the server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdonw server", slog.String("error", err.Error()))
		}
		slog.Info("server shutdown successfully")
	}

}

//? to handle gracefull shutdown we listenandserve in different thread and handle the shutdown using channel.
