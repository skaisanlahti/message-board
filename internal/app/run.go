package app

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/message-board/internal/assert"
	"github.com/skaisanlahti/message-board/internal/config"
)

func Run(
	ctx context.Context,
	args []string,
	getEnv func(string) string,
	stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	configuration, err := config.Read("appsettings.json")
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(stderr, nil))
	assert.SetLogger(logger)

	database, err := sql.Open("pgx", configuration.DatabaseAddress)
	if err != nil {
		return err
	}
	defer database.Close()

	appHandler := NewApp(database, logger)
	httpServer := &http.Server{
		Addr:    configuration.ServerAddress,
		Handler: appHandler,
	}

	go func() {
		logger.Info(
			"server started",
			slog.String("address", configuration.ServerAddress),
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
