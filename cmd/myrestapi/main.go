package main

import (
	config "dev/myrestapi/internal"
	"dev/myrestapi/internal/http-server/handlers/save"
	mwLogger "dev/myrestapi/internal/http-server/middleware/logger"
	"dev/myrestapi/internal/logger/handlers/sl"
	"dev/myrestapi/internal/logger/handlers/slogpretty"
	"dev/myrestapi/internal/storage/sqlite"
	stdLog "log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//TODO init config : cleanenv
	cfg := config.MustLoad()

	//TODO init logger: slog
	log := setuplogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	//TODO init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	//_ = storage //временное решение

	//init router: chi, "chi render"
	router := chi.NewRouter()
	router.Use(middleware.RequestID) //добавляет request_id в каждый запрос
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))    // переопределение штатного логера
	router.Use(middleware.URLFormat) //Парсер URLов поступающих запросов

	router.Post("/name", save.New(log, storage)) //Парсер URLов поступающих запросов

	//TODO run server:
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")
	stdLog.Printf("запущен сервер %s", cfg.Address) //просто проверка работы стандартного логера

	if err := srv.ListenAndServe(); err != nil {
		log.Info("невозможно запустить сервер:", cfg.Address, err)
	}

}

func setuplogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
		// log = slog.New(
		// 	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// )
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
