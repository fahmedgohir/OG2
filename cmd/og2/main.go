package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kelseyhightower/envconfig"
	"hunter.io/og2/internal/og2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	VolumePath string `required:"true" split_words:"true" desc:"Path to volume mount"`
}

func main() {
	var cfg config
	if err := envconfig.Process("og2", &cfg); err != nil {
		log.Panic(err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/sessions.db", cfg.VolumePath))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	sessions, err := og2.NewSessions(db)
	if err != nil {
		log.Panic(err)
	}
	sessions.Start()
	defer sessions.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	handler := og2.NewHandler(sessions)
	handler.Route(router)

	server := http.Server{Addr: fmt.Sprintf(":%d", 8081), Handler: router}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	fmt.Println("Server is running")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Panic(err)
	}
}
