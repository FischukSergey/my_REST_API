package main

import (
	config "dev/myrestapi/internal"
	"dev/myrestapi/internal/logger/handlers/sl"
	"dev/myrestapi/internal/logger/handlers/slogpretty"
	"dev/myrestapi/internal/storage/sqlite"
	stdLog "log"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Обработчик главной страницы
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //если домашняя страница обозначена иначе чем "/", то прервать обработку
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Привет от моего сервера"))
}

// Обработчик для отображения заметки
func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Отображение заметки ..."))
}

// Обработчик для создания новой заметки
func createSnippet(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Форма для создания новой заметки ..."))
}

func main() {
	//TODO init config : cleanenv
	cfg := config.MustLoad()

	//TODO init logger: slog
	log := setuplogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	// log.Info("starting myrestapi", slog.String("env", cfg.Env))
	// log.Debug("debug massage are enabled")

	//TODO init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	/*
		// ручной тест метода сохранения записи
		testSave := save.Request{
			Genre: "портрет",
			Name:  "автопортрет",
			Size:  "30*40",
		}
		id, err := storage.SavePicture(&testSave)
		if err != nil {
			log.Error("failed to save picture", sl.Err(err))
			os.Exit(1)
		}
		log.Info("saved picture ", slog.Int64("id", id))
	*/

	/*
	// ручной тест метода выборки из базы данных

	pictureResp, err := storage.GetPicture("пейзаж")
	if err != nil {
		log.Error("failed to get picture", sl.Err(err))
		os.Exit(1)
	}
	for _, v := range pictureResp {
		fmt.Printf("Genre: %s, Name: %s, Size: %s.\n",v.Genre,v.Name,v.Size )
	}
	*/

	_ = storage //временное решение

	//TODO init router: chi, "chi render"
	//Регистрируем два обработчика и соответствующие URL-шаблоны
	//маршрутизатора servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//TODO run server:
	log.Info("Initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")
	stdLog.Printf("запущен сервер %s", cfg.Address) //просто проверка работы стандартного логера
	if err := http.ListenAndServe(cfg.Address, mux); err != nil {
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
