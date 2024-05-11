package program

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/message-board/library/assert"
	"github.com/skaisanlahti/message-board/library/file"
)

type appSettings struct {
	ServerAddress   string `json:"serverAddress"`
	DatabaseAddress string `json:"databaseAddress"`
}

func Run(
	ctx context.Context,
	stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	settingsFilePath := flag.String("settings", "", "application settings")
	flag.Parse()

	if *settingsFilePath == "" {
		return errors.New("no settings were provided")
	}

	config, err := file.ReadJSON[appSettings](*settingsFilePath)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(stderr, nil))
	assert.SetLogger(logger)
	logger.Info(
		"settings loaded",
		slog.String("file", *settingsFilePath),
	)

	database, err := sql.Open("pgx", config.DatabaseAddress)
	if err != nil {
		return err
	}
	defer database.Close()

	handler := newHandler(database, logger)
	httpServer := http.Server{
		Addr:    config.ServerAddress,
		Handler: handler,
	}

	go func() {
		logger.Info(
			"server started",
			slog.String("address", httpServer.Addr),
		)

		err = httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			logger.Info("server closed")
			return
		}

		if err != nil {
			logger.Error(
				"error listening and serving",
				slog.Any("error", err),
			)
		}
	}()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancelShutdown := context.WithTimeout(ctx, 10*time.Second)
		defer cancelShutdown()

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			logger.Error(
				"error shutting down http server",
				slog.Any("error", err),
			)
		}
	}()

	waitGroup.Wait()
	return nil
}
