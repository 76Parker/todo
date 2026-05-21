// entrypoint to application, handle OS signals and graceful shutdown
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/internal/app"
	"todo/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	defer stop()
	var cfgPath string
	if os.Getenv("CONFIG_PATH") == "" {
		cfgPath = "config.yaml"
	} else {
		cfgPath = os.Getenv("CONFIG_PATH")
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	starter, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := starter.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	disconnectedCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer starter.Close()
	defer cancel()
	if err := starter.Shutdown(disconnectedCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
