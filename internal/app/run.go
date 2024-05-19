package app

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
	"github.com/skaisanlahti/message-board/internal/app/web"
	"github.com/skaisanlahti/message-board/internal/pkg/assert"
	"github.com/skaisanlahti/message-board/internal/pkg/file"
	"github.com/skaisanlahti/message-board/internal/pkg/password"
	"github.com/skaisanlahti/message-board/internal/pkg/session"
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

	settings, err := file.ReadJSON[appSettings](*settingsFilePath)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(stderr, nil))
	assert.SetLogger(logger)
	logger.InfoContext(
		ctx,
		"settings loaded",
		slog.String("file", *settingsFilePath),
	)

	templates := web.ParseTemplates()
	webService := web.NewService(logger, templates)

	database, err := sql.Open("pgx", settings.DatabaseAddress)
	if err != nil {
		return err
	}
	defer database.Close()

	sessionOptions := session.Options{
		CookieName:      "sid",
		SessionDuration: 1 * time.Hour,
	}

	sessionService := session.NewService(sessionOptions)
	passwordOptions := password.Options{
		Time:    5,
		Memory:  1024 * 7,
		Threads: 1,
		SaltLen: 32,
		KeyLen:  64,
	}

	passwordService := password.NewService(passwordOptions)

	server := newServer(
		logger,
		database,
		webService,
		passwordService,
		sessionService,
	)

	httpServer := http.Server{
		Addr:    settings.ServerAddress,
		Handler: server,
	}

	go func() {
		logger.InfoContext(
			ctx,
			"server started",
			slog.String("address", httpServer.Addr),
		)

		err = httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			logger.InfoContext(ctx, "server closed")
			return
		}

		if err != nil {
			logger.ErrorContext(
				ctx,
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
			logger.ErrorContext(
				shutdownCtx,
				"error shutting down http server",
				slog.Any("error", err),
			)
		}
	}()

	waitGroup.Wait()
	return nil
}
