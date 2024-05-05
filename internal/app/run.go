package app

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/message-board/internal/config"
)

func Run(
	ctx context.Context,
	args []string,
	getEnv func(string) string,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// initialize and build services
	configuration, err := config.Read("appsettings.json")
	if err != nil {
		return err
	}

	database, err := sql.Open("pgx", configuration.DatabaseAddress)
	if err != nil {
		return err
	}

	// build app and routes
	app := NewApp(database)
	httpServer := &http.Server{
		Addr:    configuration.ServerAddress,
		Handler: app,
	}

	// start server
	go func() {
		fmt.Fprintf(stdout, "server listening to %s", configuration.ServerAddress)
		err = httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving %s\n", err)
		}
	}()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// stop server
	go func() {
		defer waitGroup.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}

	}()

	waitGroup.Wait()
	return nil
}
