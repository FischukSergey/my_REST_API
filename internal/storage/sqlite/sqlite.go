package sqlite

import (
	"database/sql"
	"dev/myrestapi/internal/http-server/handlers/save"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //инициализация драйвера sqlite
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New" //имя текущей функции для логов и ошибок

	db, err := sql.Open("sqlite3", storagePath) //подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем таблицу picture если ее еще нет
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS picture(
			id INTEGER PRIMARY KEY,
			genre TEXT NOT NULL,
			name TEXT NOT NULL,
			size TEXT);
		CREATE INDEX IF NOT EXISTS idx_alias ON picture(alias);
		`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SavePicture метод добавления записи в базу данных.
// возвращает ID записи или ошибку
func (s *Storage) SavePicture(r *save.Request) (int64, error) {
	const op = "storage.sqlite.SavePicture"

	// TODO: проверка на наличие записи с одинаковым именем картины
	/*
	   stmtExist, err := s.db.Prepare("SELECT EXISTS (SELECT id FROM picture WHERE name=?)")
	   	if err != nil {
	   		return 0, fmt.Errorf("%s: prepare search: %w", op, err)
	   	}
	   exist,err:=stmtExist.Exec(r.Name)
	   	if err != nil {
	   		return 0, fmt.Errorf("%s: executive search: %w", op, err)
	   	}
	   	if exist>0 || true {
	   		return 0,
	   	}
	*/

	//Подготавливаем запрос
	stmt, err := s.db.Prepare("INSERT INTO picture(genre,name,size) VALUES (?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	//Выполняем запрос
	res, err := stmt.Exec(r.Genre, r.Name, r.Size)
	if err != nil {
		return 0, fmt.Errorf("%s: executive statement: %w", op, err)
	}

	// Получаем ID созданной записи
	// Используем метод LastInsertID из пакета SQL
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	//Возвращаем ID новой записи
	return id, nil
}

// GetPicture() метод получения выборки картин по жанру.
// возвращает срез строк (фильтр) или ошибку
func (s *Storage) GetPicture(alias string) ([]save.Request, error) {
	const op = "storage.sqlite.GetPicture"

	// готовим запрос по жанру картину
	stmt, err := s.db.Prepare("SELECT genre,name,size FROM picture WHERE genre=?")
	if err != nil {
		return []save.Request{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	// переменная для сохранения результат
	resGenre := []save.Request{}

	// выполняем запрос, выбираем строки соответствующие запросу
	resp, err := stmt.Query(alias)
	if err != nil {
		return []save.Request{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	
	// заносим строки в результирующий срез
	for resp.Next() {
	
		picture := save.Request{}

		err := resp.Scan(&picture.Genre, &picture.Name, &picture.Size)
		if err != nil {
			fmt.Printf("%s: not row execute statement: %s \n", op, err)
			continue
		}
		resGenre = append(resGenre, picture)
	}

	return resGenre, nil
}
