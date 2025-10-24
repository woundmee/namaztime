package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type Database struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Database {
	return &Database{
		logger: logger,
	}
}

const DATABASE_PATH = "./users.db"
const DATABASE_TABLE = "users"

func (d *Database) connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", DATABASE_PATH)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) Create() error {
	db, err := d.connect()
	if err != nil {
		d.logger.Error("ошибка подключения к БД", "error", err)
		return err
	}

	defer db.Close()

	queryCreateTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        chatID INTEGER NOT NULL UNIQUE,
        username TEXT
    );`

	_, err = db.Exec(queryCreateTable)
	if err != nil {
		d.logger.Error("ошибка создания таблицы", "error", err)
		return err
	}

	return nil
}

func (d *Database) AddUser(chatID int64, username string) error {
	db, err := d.connect()
	if err != nil {
		d.logger.Error("ошибка подключения к БД", "error", err)
		return err
	}

	defer db.Close()

	// // fixme: добавить проверку на наличие пользователя в БД: возникает ошибка в логах
	users, err := d.GetUsers()
	if err != nil {
		d.logger.Error("не удалось получить список пользователей из БД", "error", err)
		return err
	}

	// проверка на наличие пользователя в БД
	for userChatID := range users {
		if chatID == userChatID {
			d.logger.Info("пользователь повторно вызвал команду /start", "chatID", chatID)
			return nil
		}
	}

	queryInsert := "insert into users (chatID, username) values (?, ?)"
	_, err = db.Exec(queryInsert, chatID, username)
	if err != nil {
		d.logger.Error("не удалось записать новые данные в БД", "chatID", chatID, "username", username, "erorr", err)
		return err
	}

	d.logger.Info("в БД добавлен новый пользователь", "chatID", chatID, "username", username)

	return nil
}

func (d *Database) DeleteUser(chatID int64) error {
	db, err := d.connect()
	if err != nil {
		d.logger.Error("ошибка подключения к БД", "error", err)
		return err
	}

	defer db.Close()

	queryDelete := fmt.Sprintf("delete from %s where chatID = ?", DATABASE_TABLE)
	res, err := db.Exec(queryDelete, chatID)
	if err != nil {
		d.logger.Error("не удалось удалить запись!", "chatID", chatID, "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		d.logger.Error("не удалось получить кол-во удаленных строк", "error", err)
		return err
	}

	if rowsAffected == 0 {
		d.logger.Warn("запись для удаления не найдена", "chatID", chatID)
	} else {
		d.logger.Info("Пользователь удален из БД", "chatID", chatID)
	}

	return nil
}

func (d *Database) GetUsers() (map[int64]string, error) {
	db, err := d.connect()
	if err != nil {
		d.logger.Error("ошибка подключения к БД", "error", err)
		return nil, err
	}

	defer db.Close()

	querySelect := fmt.Sprintf("select chatID, username from %s", DATABASE_TABLE)

	rows, err := db.Query(querySelect)
	if err != nil {
		d.logger.Error("ошибка получения данных из БД", "error", err)
		return nil, err
	}
	defer rows.Close()

	users := make(map[int64]string)

	for rows.Next() {
		var chatID int64
		var username string

		err = rows.Scan(&chatID, &username)
		if err != nil {
			d.logger.Error("ошибка сканирования данных из БД", "error", err)
			return nil, err
		}

		users[chatID] = username
	}

	// if len(users) == 0 {
	// 	return nil, fmt.Errorf("список пользователей пуст: %s", err)
	// }

	return users, nil
}
