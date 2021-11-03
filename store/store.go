package store

import (
	"database/sql"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore(cfg *config.Config) (*Store, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `preset` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `uri` TEXT UNIQUE, `sid` INTEGER)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `stream` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `name` TEXT UNIQUE, `content` TEXT)")
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}
