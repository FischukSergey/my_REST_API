package main

import (
	config "dev/myrestapi/internal"
	"log"
	"net/http"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Обработчик главной страницы
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {  //если домашняя страница обозначена иначе чем "/", то прервать обработку
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

	//TODO init storage: sqlite

	//TODO init router: chi, "chi render"
	//Регистрируем два обработчика и соответствующие URL-шаблоны
	//маршрутизатора servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//TODO run server:
	log.Printf("Запуск сервера на %s\n", cfg.Address)
	err := http.ListenAndServe(cfg.Address, mux)
	log.Fatalf("невозможно запустить сервер: %s", err)
}
